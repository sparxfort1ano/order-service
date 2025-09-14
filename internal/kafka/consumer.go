package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/sparxfort1ano/order-service/internal/repository"
	"github.com/sparxfort1ano/order-service/internal/service"
)

// Потребитель слушает Kafka и отправляет заказы в сервис
type Consumer struct {
	reader *kafka.Reader
	svc    *service.Service
}

// Создание нового потребителя
func New(brokers []string, topic, group string, svc *service.Service) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,           // адреса брокеров
		Topic:       topic,             // топик для чтения
		GroupID:     group,             // группа потребителя
		MinBytes:    1e3,               // минимум байт в сообщении
		MaxBytes:    10e6,              // максимум байт
		MaxWait:     time.Second,       // ждать не дольше этой паузы
		StartOffset: kafka.FirstOffset, // читаем все сообщения с начала
	})
	return &Consumer{reader: r, svc: svc}
}

// Запуск фоновой горутины для чтения сообщений из Kafka
func (c *Consumer) Start(ctx context.Context) {
	go func() {
		defer c.reader.Close()
		for {
			// читаем сообщение из топика
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return // завершение по контексту
				}
				log.Printf("kafka read: %v", err)
				continue
			}

			// Парсим JSON в структуру заказа
			var ord repository.Order
			if err := json.Unmarshal(m.Value, &ord); err != nil {
				log.Printf("kafka bad json (skip): %v", err)
				continue
			}

			// игнорируем сообщения без order_uid
			if ord.OrderUID == "" {
				log.Printf("kafka message without order_uid (skip)")
				continue
			}

			// сохраняем заказ через сервис
			if err := c.svc.Save(ctx, ord); err != nil {
				log.Printf("save failed: %v", err)
				continue
			}
			log.Printf("kafka saved order %s", ord.OrderUID)
		}
	}()
}
