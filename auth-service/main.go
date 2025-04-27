package main

import (
	"log"

	"github.com/RohithBN/auth-service/handlers"
	"github.com/RohithBN/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := utils.ConnectDB()
	if err := handlers.InitDB(); err != nil {
		log.Fatalf("Error initializing database: %v\n", err)
	}
	if err != nil {
		log.Fatalf("Error connecting to the database: %v\n", err)
	}
	defer db.Close()

	router := gin.Default()

	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)

	log.Printf("Auth service starting on port 8081")
	router.Run(":8081")
}
