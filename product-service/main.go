package main

import (
	"log"

	"github.com/RohithBN/product-service/handlers"
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

	router.GET("/products", handlers.GetProducts)
	router.POST("/add-product", handlers.AddProduct)
	router.GET("/products/:id", handlers.GetProductByID)
	router.PUT("/update-product/:id", handlers.UpdateProduct)
	router.DELETE("/delete-product/:id", handlers.DeleteProduct)

	log.Printf("Product service starting on port 8082")
	router.Run(":8082")
}
