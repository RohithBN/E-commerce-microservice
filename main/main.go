package main

import (
	"context"
	"log"

	"github.com/RohithBN/handlers"
	"github.com/RohithBN/middlewares"
	"github.com/RohithBN/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	db, err := utils.ConnectDB()
	if err := handlers.InitDB(); err != nil {
		log.Fatalf("Error initializing database: %v\n", err)
	}
	if err != nil {
		log.Fatalf("Error connecting to the database: %v\n", err)
	}
	defer db.Close()

	if err := utils.ConnectMongoDB(); err != nil {
		log.Fatalf("Error connecting to MongoDB: %v\n", err)
	}
	defer utils.MongoClient.Disconnect(context.Background())

	router := gin.Default()

	router.POST("/register", handlers.Register)

	router.POST("/login", handlers.Login)

	auth := router.Group("/auth")
	auth.Use(middlewares.AuthMiddleware())
	{
		//product routes
		auth.GET("/products", handlers.GetProducts)
		auth.POST("/add-product", handlers.AddProduct)
		auth.GET("/products/:id", handlers.GetProductByID)
		auth.PUT("/update-product/:id", handlers.UpdateProduct)
		auth.DELETE("/delete-product/:id", handlers.DeleteProduct)

		//cart routes
		auth.POST("/add-to-cart/:productId", handlers.AddToCart)
		auth.GET("/cart", handlers.GetCart)

		//order routes
		auth.POST("/create-order", handlers.CreateOrder)
		auth.POST("/orders/payment", handlers.ProcessPayment)
		auth.PUT("/orders/:orderId/status", handlers.UpdateOrderStatus)
		auth.GET("/orders", handlers.GetOrders)
	}

	router.Run(":8080")

}
