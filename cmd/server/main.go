package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sparxfort1ano/order-service/internal/cache"
	"github.com/sparxfort1ano/order-service/internal/config"
	"github.com/sparxfort1ano/order-service/internal/handler"
	"github.com/sparxfort1ano/order-service/internal/kafka"
	"github.com/sparxfort1ano/order-service/internal/repository"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Config Loading
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cfg.AppPort)

	// DB init
	repo, err := repository.NewPostgresRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	repo.Migrate("./migrations/000001_init.up.sql")
	log.Println("Database initialized and migrated successfully!")

	// Cache pre-warming
	orderCache := cache.NewOrderCache()
	orders, err := repo.GetAllOrders(ctx)
	if err != nil {
		log.Printf("warning: failed to warm up cache: %v", err)
	} else {
		for _, o := range orders {
			orderCache.Set(o)
		}
		log.Printf("Successfully warmed up cache with %d orders", len(orders))
	}

	// Init producer
	producer := kafka.NewOrderProducer(cfg.KafkaBroker, cfg.KafkaTopic)
	go producer.Run(ctx)

	// Init consumer
	consumer := kafka.NewOrderConsumer(cfg.KafkaBroker, cfg.KafkaTopic, repo, orderCache)
	go consumer.Start(ctx)

	// HTTP server config
	h := handler.NewOrderHandler(orderCache, repo)
	mux := http.NewServeMux()
	mux.HandleFunc("/order", h.GetOrder)
	mux.Handle("/", http.FileServer(http.Dir("./web/static")))
	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: mux,
	}

	// Server run
	go func() {
		log.Printf("Server started on http://localhost:%s", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	<-ctx.Done()
	log.Println("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
