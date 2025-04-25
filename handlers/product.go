package handlers

import (
	"context"
	"time"

	"github.com/RohithBN/types"
	"github.com/RohithBN/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProducts(c *gin.Context) {
	var products []types.Product

	collection := utils.MongoDB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch products",
			"details": err.Error(),
		})

		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var product types.Product
		if err := cursor.Decode(&product); err != nil {
			c.JSON(500, gin.H{"error": "Failed to decode product"})
			return
		}
		products = append(products, product)
	}
	if err := cursor.Err(); err != nil {
		c.JSON(500, gin.H{"error": "Cursor error"})
		return
	}
	c.JSON(200, gin.H{"products": products})
}

func AddProduct(c *gin.Context) {
	var product types.Product
	if err := c.BindJSON(&product); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid Product Details",
		})
		return
	}
	product.CreatedAt = time.Now().Format(time.RFC3339)

	collection := utils.MongoDB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, product)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to add product"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Product added successfully",
		"product": product,
	})
}
func GetProductByID(c *gin.Context) {
	var product types.Product
	productID := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid product ID format"})
		return 
	}

	collection := utils.MongoDB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return 
	}
	c.JSON(200, gin.H{"product": product})
}

func UpdateProduct(c *gin.Context) {
	var product types.Product
	if err := c.BindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": "Invalid Product Details"})
		return
	}
	product.UpdatedAt = time.Now().Format(time.RFC3339)

	// Convert ID from string to ObjectID
	productID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid product ID"})
		return
	}

	collection := utils.MongoDB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"name":         product.Name,
		"price":        product.Price,
		"description":  product.Description,
		"updated_at":   product.UpdatedAt,
		"added_to_cart": product.AddedToCart,
		"category":     product.Category,
		"stock":        product.Stock,
	}
	
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
		if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Product updated successfully",
	})
}

func DeleteProduct(c *gin.Context) {
	productID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid product ID"})
		return
	}

	collection := utils.MongoDB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete product"})
		return
	}
	if result.DeletedCount == 0 {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Product deleted successfully"})
}
