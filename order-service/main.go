package main

import (
	"log"

	"github.com/RohithBN/order-service/handlers"
	"github.com/RohithBN/shared/metrics"
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

	router := gin.Default()
	router.Use(metrics.PrometheusMiddleware())

	metrics.RegisterMetricsEndpoint(router)

	router.POST("/create-order", handlers.CreateOrder)
	router.POST("/orders/payment", handlers.ProcessPayment)
	router.PUT("/orders/:orderId/status", handlers.UpdateOrderStatus)
	router.GET("/orders", handlers.GetOrders)

	log.Printf("Order service starting on port 8084")
	router.Run(":8084")
}
