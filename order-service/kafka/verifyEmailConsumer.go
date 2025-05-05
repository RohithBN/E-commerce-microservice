package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/RohithBN/shared/utils"
	"github.com/segmentio/kafka-go"
)

type OTPEmail struct {
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

func VerifyOTPEmailConsumer(ctx context.Context) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "send-verify-otp-email",
		GroupID: "otp-email-group",
	})
	defer reader.Close()

	for {
		select {
		case <-ctx.Done():
			log.Println("OTP Email consumer shutting down...")
			return ctx.Err()
		default:
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return ctx.Err() // Context was canceled
				}
				log.Printf("Error reading message: %v", err)
				continue
			}

			var sendOTP OTPEmail
			err = json.Unmarshal(m.Value, &sendOTP)
			if err != nil {
				log.Printf("Error decoding OTP email payload: %v", err)
				continue
			}

			fmt.Println("receiverd at sonsumer , sending to email function")

			err = utils.SendOTPMail(sendOTP.Email, sendOTP.CreatedAt)
			if err != nil {
				log.Printf("Error sending OTP mail: %v", err)
				continue
			}

			log.Printf("Email: %s, CreatedAt: %s", sendOTP.Email, sendOTP.CreatedAt)
			log.Printf("Successfully sent OTP mail to %s", sendOTP.Email)
		}
	}
}
