package main

import (
	"context"
	"log"

	"github.com/RohithBN/order-service/handlers"
	"github.com/RohithBN/order-service/kafka"
	"github.com/RohithBN/shared/metrics"
	"github.com/RohithBN/shared/redis"
	"github.com/RohithBN/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := utils.ConnectMongoDB(); err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	if err := redis.ConnectRedis(); err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	router := gin.Default()
	router.Use(metrics.PrometheusMiddleware())

	// Create a background context that won't time out
ctx, cancel := context.WithCancel(context.Background())
defer cancel() // Will only be called when main() exits

// Start Kafka consumer
go func() {
    if err := kafka.VerifyOTPEmailConsumer(ctx); err != nil {
        log.Printf("OTP Email consumer stopped: %v", err)
    }
}()

	metrics.RegisterMetricsEndpoint(router)
	//public routes
	router.GET("/orders", handlers.GetOrders)

	//two factor authentication routes
	router.POST("/orders/send-otp", handlers.SendOTP)
	router.POST("/orders/verify-otp", handlers.VerifyOTP)

	//protected routes
	router.Use(handlers.TwoFacAuthMiddleware())
	router.POST("/create-order", handlers.CreateOrder)
	router.POST("/orders/payment", handlers.ProcessPayment)
	router.PUT("/orders/:orderId/status", handlers.UpdateOrderStatus)

	log.Printf("Order service starting on port 8084")
	router.Run(":8084")
}
