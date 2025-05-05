package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

var writer *kafka.Writer

func VerifyOTPEmailProducer(email string) error {
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	// inititalise the kafka writer
	writer = &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "send-verify-otp-email",
		Balancer: &kafka.LeastBytes{},
	}

	event := map[string]string{
		"email":     email,
		"createdAt": createdAt,
	}

	payload, _ := json.Marshal(event)
	fmt.Println("Payload: ", string(payload))
	return writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(email),
			Value: payload,
		})

}
