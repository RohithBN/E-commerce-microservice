package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/RohithBN/shared/utils"
	"github.com/segmentio/kafka-go"
)

func ConsumeEmail() {
	// Call the context version with a background context
	ConsumeEmailWithContext(context.Background())
}

func ConsumeEmailWithContext(ctx context.Context) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "email-topic",
		GroupID: "email-group",
	})
	defer reader.Close()

	for {
		select {
		case <-ctx.Done():
			// Context was canceled, exit gracefully
			log.Println("Email consumer shutting down...")
			return nil
		default:
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			var event struct {
				Email     string `json:"email"`
				Name      string `json:"name"`
				CreatedAt string `json:"createdAt"`
			}

			if err := json.Unmarshal(m.Value, &event); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}

			log.Printf("Received message: %s", string(m.Value))
			log.Printf("Email: %s, Name: %s, CreatedAt: %s", event.Email, event.Name, event.CreatedAt)

			if err := utils.SendEmailAfterRegistration(event.Email, event.Name, event.CreatedAt); err != nil {
				log.Printf("Error sending email: %v", err)
			}
		}
	}
}
