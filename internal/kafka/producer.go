package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/segmentio/kafka-go"
	"github.com/sparxfort1ano/order-service/internal/repository"
)

type OrderProducer struct {
	writer *kafka.Writer
}

func (p *OrderProducer) Run(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	log.Println("Kafka Producer started...")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			order := generateFakeOrder()

			jsonData, err := json.Marshal(order)
			if err != nil {
				log.Printf("failed to unmarshal data: %v", err)
				continue
			}

			err = p.writer.WriteMessages(ctx, kafka.Message{
				Key:   []byte(order.OrderUid),
				Value: jsonData,
			})

			if err != nil {
				log.Printf("failed [producer] to write message: %v", err)
			} else {
				log.Printf("Producer sent order: %s", order.OrderUid)
			}
		}
	}
}

func NewOrderProducer(broker, topic string) *OrderProducer {
	return &OrderProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func generateFakeOrder() repository.Order {
	// These fields occur multiple times in struct
	orderID := gofakeit.Regex("[a-zA-Z0-9]{15}")
	trackerNumber := "WBTESTTRACK" + gofakeit.Regex("[a-zA-Z0-9]{4}")

	return repository.Order{
		OrderUid:          orderID,
		TrackNumber:       trackerNumber,
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        gofakeit.Regex("[a-zA-Z0-9]{7}"),
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              int64(gofakeit.Number(0, 100)),
		DateCreated:       time.Now(),
		OofShard:          "1",
		Delivery: repository.Delivery{
			Name:    gofakeit.FirstName() + " " + gofakeit.LastName(),
			Phone:   "+7" + gofakeit.Phone(),
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: "Lenina 10",
			Region:  gofakeit.State(),
			Email:   gofakeit.Email(),
		},
		Payment: repository.Payment{
			Transaction:  orderID,
			RequestID:    "",
			Currency:     gofakeit.CurrencyShort(),
			Provider:     "wbpay",
			Amount:       int64(gofakeit.Number(1000, 5000)),
			PaymentDt:    time.Now().Unix(),
			Bank:         "alpha",
			DeliveryCost: int64(gofakeit.Number(100, 500)),
			GoodsTotal:   int64(gofakeit.Number(500, 4500)),
			CustomFee:    int64(gofakeit.Number(0, 100)),
		},
		Items: []repository.Item{
			{
				ChrtID:      int64(gofakeit.Number(1000000, 9999999)),
				TrackNumber: trackerNumber,
				Price:       int64(gofakeit.Number(100, 1000)),
				Rid:         gofakeit.UUID(),
				Name:        "Product Name",
				Sale:        int64(gofakeit.Number(0, 50)),
				Size:        "XL",
				TotalPrice:  int64(gofakeit.Number(100, 1000)),
				NmID:        int64(gofakeit.Number(1000000, 9999999)),
				Brand:       gofakeit.Company(),
				Status:      202,
			},
		},
	}
}
