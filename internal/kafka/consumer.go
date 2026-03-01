package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/sparxfort1ano/order-service/internal/cache"
	"github.com/sparxfort1ano/order-service/internal/repository"
)

type OrderConsumer struct {
	reader *kafka.Reader
	repo   *repository.PostgresRepository
	cache  *cache.Cache
}

func (c *OrderConsumer) Start(ctx context.Context) {
	log.Println("Kafka Consumer started...")
	defer c.reader.Close()

	for {
		// Kafka reading
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				log.Println("Kafka reader stopped by context")
				return
			}
			log.Printf("failed to read from Kafka: %v", err)
			return
		}

		// Unmarshaling
		var order repository.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("failed to unmarshal order: %v, Data: %s", err, string(msg.Value))
			continue
		}

		// Validation
		if err := order.Validate(); err != nil {
			log.Printf("invalid order data for UID %s: %v", order.OrderUid, err)
			continue
		}

		// Inserting into DB
		if err := c.repo.SaveOrder(ctx, &order); err != nil {
			log.Printf("failed to save order %s to DB: %v", order.OrderUid, err)
			continue
		}

		// Do an offset
		c.reader.CommitMessages(ctx, msg)

		// Cache update
		c.cache.Set(&order)

		log.Printf("Successfully processed order: %s", order.OrderUid)
	}
}

// Init consumer
func NewOrderConsumer(broker, topic string,
	repo *repository.PostgresRepository, cache *cache.Cache) *OrderConsumer {
	return &OrderConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{broker},
			Topic:    topic,
			GroupID:  "order-group",
			MinBytes: 1000,
			MaxBytes: 1000000,
		}),
		repo:  repo,
		cache: cache,
	}
}
