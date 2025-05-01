package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/RohithBN/auth-service/handlers"
	"github.com/RohithBN/auth-service/kafka"
	"github.com/RohithBN/shared/metrics"
	"github.com/RohithBN/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// First load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Then initialize database
	db, err := utils.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v\n", err)
	}
	defer db.Close()

	if err := handlers.InitDB(); err != nil {
		log.Fatalf("Error initializing database: %v\n", err)
	}

	// Initialize Kafka components
	kafka.InitKafkaWriter()

	// Create a cancelable context for the consumer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consumer with context
	go func() {
		if err := kafka.ConsumeEmailWithContext(ctx); err != nil {
			log.Printf("Error starting Kafka consumer: %v", err)
		}
	}()

	// Setup router
	router := gin.Default()
	router.Use(metrics.PrometheusMiddleware())

	metrics.RegisterMetricsEndpoint(router)
	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)

	// Handle shutdown gracefully
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")
		cancel() // Cancel the Kafka consumer context
	}()

	log.Printf("Auth service starting on port 8081")
	router.Run(":8081")
}
