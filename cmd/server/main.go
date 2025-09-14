package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sparxfort1ano/order-service/internal/cache"
	"github.com/sparxfort1ano/order-service/internal/config"
	"github.com/sparxfort1ano/order-service/internal/handler"
	"github.com/sparxfort1ano/order-service/internal/kafka"
	"github.com/sparxfort1ano/order-service/internal/repository"
	"github.com/sparxfort1ano/order-service/internal/service"
)

func main() {
	// Загружаем конфигурацию 
	cfg := config.Load()
	log.Printf("cfg: http=%s db=%s kafka=%s topic=%s group=%s",
		cfg.HTTPAddr, cfg.DBDSN, cfg.KafkaBroker, cfg.KafkaTopic, cfg.KafkaGroupID)

	// Контекст с отменой при завершении программы (Ctrl+C, SIGTERM)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Подключение к Postgres
	pool, err := pgxpool.New(ctx, cfg.DBDSN)
	if err != nil {
		log.Fatal("pg:", err)
	}
	defer pool.Close()

	// Репозиторий для работы с б/д
	repo := repository.NewPostgresRepo(pool)
	if err := repo.Init(ctx); err != nil {
		log.Fatal("pg init:", err)
	}

	// Кэшерирование заказов в памяти
	c := cache.New()

	// Сервис, объединяющий б/д и кэширование
	svc := service.New(repo, c)

	// Прогрев кэша: загружаем последние заказы из б/д
	if err := svc.Warm(ctx, 10000); err != nil {
		log.Println("warm:", err)
	}

	// Регистрируем эндпоинты HTTP-сервера
	mux := http.NewServeMux()
	handler.New(svc).Routes(mux)

	// Запуск HTTP-сервера
	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: mux}
	go func() {
		log.Println("HTTP on", cfg.HTTPAddr)
		log.Fatal(srv.ListenAndServe())
	}()

	// Запуск Kafka consumer для чтения заказов из топика
	k := kafka.New([]string{cfg.KafkaBroker}, cfg.KafkaTopic, cfg.KafkaGroupID, svc)
	k.Start(ctx)
	log.Println("Kafka consumer started on", cfg.KafkaBroker, "topic", cfg.KafkaTopic)

	// Ждем завершения (Ctrl+C или сигнал)
	<-ctx.Done()
	log.Println("shutting down...")

	// Плавное завершение работы HTTP-сервера
	sh, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(sh)
}
