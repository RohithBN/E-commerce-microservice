package utils

import (
	"context"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/RohithBN/shared/redis"
	"github.com/RohithBN/shared/types"
)

type EmailConfig struct {
	From     string
	Password string
	Host     string
	Port     string
	Address  string
}

func NewEmailConfig() *EmailConfig {
	host := "smtp.gmail.com"
	port := "587"
	return &EmailConfig{
		From:     os.Getenv("EMAIL_FROM"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		Host:     host,
		Port:     port,
		Address:  fmt.Sprintf("%s:%s", host, port),
	}
}

func SendEmail(to []string, subject string, body string) error {
	config := NewEmailConfig()
	if config.From == "" || config.Password == "" {
		return fmt.Errorf("missing email configuration")
	}

	auth := smtp.PlainAuth("", config.From, config.Password, config.Host)

	message := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n",
		config.From,
		strings.Join(to, ","),
		subject,
		body))

	err := smtp.SendMail(config.Address, auth, config.From, to, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func SendEmailAfterRegistration(email string, name string, createdAt string) error {
	subject := "Welcome to E-Commerce Store"

	body := fmt.Sprintf(`
        <html>
        <head>
            <style>
                body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
                .container { max-width: 600px; margin: 0 auto; padding: 20px; }
                .header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
                .footer { text-align: center; margin-top: 20px; color: #666; }
            </style>
        </head>
        <body>
            <div class="container">
                <div class="header">
                    <h1>Welcome to E-Commerce Store</h1>
                </div>
                
                <p>Dear %s,</p>
                <p>Thank you for registering with us! We're excited to have you on board.</p>
                
                <p>Your account was created on %s.</p>
                
                <div class="footer">
                    <p>Thank you for joining us!</p>
                    <small>This is an automated email, please do not reply.</small>
                </div>
            </div>
        </body>
        </html>
    `, name, createdAt)

	return SendEmail([]string{email}, subject, body)
}

func SendOrderConfirmationEmail(toEmail string, order *types.Order) error {
	subject := "Order Confirmation - E-Commerce Store"

	// Format order ID to be more user-friendly
	orderIDString := order.OrderId.Hex()

	body := fmt.Sprintf(`
        <html>
        <head>
            <style>
                body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
                .container { max-width: 600px; margin: 0 auto; padding: 20px; }
                .header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
                .order-details { background-color: #f9f9f9; padding: 20px; margin: 20px 0; border-radius: 5px; }
                .footer { text-align: center; margin-top: 20px; color: #666; }
                .button { background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; }
            </style>
        </head>
        <body>
            <div class="container">
                <div class="header">
                    <h1>Order Confirmation</h1>
                </div>
                
                <p>Dear Customer,</p>
                <p>Thank you for your order! We're pleased to confirm that we've received your order.</p>
                
                <div class="order-details">
                    <h3>Order Details:</h3>
                    <p><strong>Order ID:</strong> #%s</p>
                    <p><strong>Total Amount:</strong> $%.2f</p>
                    <p><strong>Status:</strong> %s</p>
                    <p><strong>Order Date:</strong> %s</p>
                </div>
                
                <p>We'll notify you when your order ships. If you have any questions, please don't hesitate to contact us.</p>
                
                <div class="footer">
                    <p>Thank you for shopping with us!</p>
                    <small>This is an automated email, please do not reply.</small>
                </div>
            </div>
        </body>
        </html>
    `, orderIDString[:8], order.TotalPrice, order.Status, order.CreatedAt)

	return SendEmail([]string{toEmail}, subject, body)
}

func SendOTPMail(email string, createdAt string) error {
	uniqueCode := generateUniqueCode()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := redis.RedisClient.Set(ctx, email, uniqueCode,  10*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to set OTP in Redis: %v", err)
	}
    fmt.Println("OTP set in Redis:", uniqueCode)
	subject := "Your OTP Code - E-Commerce Store"
	body := fmt.Sprintf(`

        <html>
        <head>
            <style>
                body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
                .container { max-width: 600px; margin: 0 auto; padding: 20px; }
                .header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
                .footer { text-align: center; margin-top: 20px; color: #666; }
            </style>
        </head>
        <body>

            <div class="container">
                <div class="header">
                    <h1>Your OTP Code</h1>
                </div>
                <p>Dear User,</p>
                <p>Your OTP code is: <strong>%s</strong></p>
                <p>This code is valid for a short period of time. Please use it to complete your registration.</p>
                <p>If you did not request this code, please ignore this email.</p>
                <div class="footer">

                    <p>Thank you for using our service!</p>
                    <small>This is an automated email, please do not reply.</small>
                </div>
            </div>
        </body>
        </html>
    `, uniqueCode)
    fmt.Println("Sending OTP email to:", email)
	return SendEmail([]string{email}, subject, body)
}

func generateUniqueCode() string {
	// Generate a unique code (e.g., a random 6-digit number)
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
