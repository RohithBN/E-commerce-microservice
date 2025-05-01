package main

import (
	"log"

	"github.com/RohithBN/cart-service/handlers"
	"github.com/RohithBN/cart-service/kafka"
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

	//inititalise kafka writer
	kafka.InitKafkaWriter()

	router := gin.Default()

	router.Use(metrics.PrometheusMiddleware())

	metrics.RegisterMetricsEndpoint(router)

	// Routes aligned with gateway
	router.POST("/cart/:productId", handlers.AddToCart)
	router.GET("/cart", handlers.GetCart)
	router.DELETE("/cart/:productId", handlers.DeleteFromCart)

	log.Printf("Cart service starting on port 8083")
	router.Run(":8083")
}
