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
		Topic:    "email-topic",
		Balancer: &kafka.LeastBytes{},
	}

}

func ProduceEmail(email string, name string, createdAt string) error {
	event := map[string]string{
		"email":     email,
		"name":      name,
		"createdAt": createdAt,
	}

	payload, _ := json.Marshal(event)

	return writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(email),
			Value: payload,
		},
)
}
