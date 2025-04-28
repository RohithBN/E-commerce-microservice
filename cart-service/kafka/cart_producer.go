package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

var writer *kafka.Writer

func InitKafkaWriter() {
	writer = &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "cart-add-item-topic",
		Balancer: &kafka.LeastBytes{},
	}
}

func ProduceCartAddItem(quantity int, productId string, userId int) error {
	event := map[string]interface{}{
		"quantity":  quantity,
		"productId": productId,
		"userId":    userId,
	}

	payload, _ := json.Marshal(event)


	return writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(productId),
		Value: payload,
	})
}
