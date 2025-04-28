package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"github.com/RohithBN/shared/utils"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConsumeCartAddItem()  {
	ConsumeCartAddItemWithContext(context.Background())
}

func ConsumeCartAddItemWithContext(ctx context.Context) error {

		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   "cart-add-item-topic",
			GroupID: "cart-group",
		})

		defer reader.Close()

		for {
			select {
			case <-ctx.Done():
				log.Println("Cart consumer shutting down...")
				return nil
			default:
				m, err := reader.ReadMessage(ctx)
				if err != nil {
					log.Printf("Error reading message: %v", err)
					continue
				}
				var event struct {
					Quantity  int    `json:"quantity"`
					ProductId string `json:"productId"`
					UserId    int    `json:"userId"`
				}
				if err := json.Unmarshal(m.Value, &event); err != nil {
					log.Printf("Error unmarshalling message: %v", err)
					continue
				}
				log.Printf("Received message: %s", string(m.Value))
				log.Printf("Quantity: %d, ProductId: %s, UserId: %d", event.Quantity, event.ProductId, event.UserId)

				//logic to update the product stock

				productCollection := utils.MongoDB.Collection("products")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				// convert productId(stirng) to ObjectID
				objectId, err := primitive.ObjectIDFromHex(event.ProductId)
				if err != nil {
					log.Printf("Error converting productId to ObjectID: %v", err)
					continue
				}

				//update the product stock

				_, err = productCollection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$inc": bson.M{"stock": -event.Quantity}})
				if err != nil {
					log.Printf("Error updating product stock: %v", err)
					continue
				}
				log.Printf("Product stock updated successfully")
			}
		}


}
