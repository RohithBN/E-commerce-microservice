package main

import (
	"log"

	"github.com/RohithBN/gateway/handlers"
	"github.com/RohithBN/gateway/middleware"
	"github.com/RohithBN/shared/metrics"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()
	router.Use(metrics.PrometheusMiddleware())

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	// Public routes
	router.POST("/api/register", handlers.ProxyHandler("auth", "/register"))
	router.POST("/api/login", handlers.ProxyHandler("auth", "/login"))

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	api.Use(middleware.RateLimitMiddleware())
	{
		// Products
		api.GET("/products", handlers.ProxyHandler("products", "/products"))
		api.POST("/add-product", handlers.ProxyHandler("products", "/add-product"))
		api.GET("/products/:id", handlers.ProxyHandler("products", "/products/:id"))
		api.PUT("/update-product/:id", handlers.ProxyHandler("products", "/update-product/:id"))
		api.DELETE("/delete-product/:id", handlers.ProxyHandler("products", "/delete-product/:id"))

		// Cart
		api.POST("/cart/:productId", handlers.ProxyHandler("cart", "/cart/:productId"))
		api.GET("/cart", handlers.ProxyHandler("cart", "/cart"))
		api.DELETE("/cart/:productId", handlers.ProxyHandler("cart", "/cart/:productId"))
		// Orders
		api.POST("/create-order", handlers.ProxyHandler("orders", "/create-order"))
		api.POST("/orders/send-otp", handlers.ProxyHandler("orders", "/orders/send-otp"))
		api.POST("/orders/verify-otp", handlers.ProxyHandler("orders", "/orders/verify-otp"))
		api.POST("/orders/payment", handlers.ProxyHandler("orders", "/orders/payment"))
		api.PUT("/orders/:orderId/status", handlers.ProxyHandler("orders", "/orders/:orderId/status"))
		api.GET("/orders", handlers.ProxyHandler("orders", "/orders"))
	}

	log.Printf("API Gateway starting on port 8080")
	router.Run(":8080")
}
