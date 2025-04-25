package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func ConnectMongoDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoUri := os.Getenv(("MONGOURI"))
	clientOptions := options.Client().ApplyURI(mongoUri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("MongoDB connect error: %v", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("MongoDB ping error: %v", err)
	}

	MongoClient = client
	MongoDB = client.Database("e-commerce") // use your DB name
	fmt.Println(" Connected to MongoDB")
	return nil
}
