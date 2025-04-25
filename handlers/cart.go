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

func AddToCart(c *gin.Context) {
    // get userId from token
    user := c.MustGet("user").(*types.User)
    user_id := user.Id
    productId := c.Param("productId")
    objectId, err := primitive.ObjectIDFromHex(productId)
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid product ID"})
        return
    }
    
    // get product from db
    var product types.Product
    productCollection := utils.MongoDB.Collection("products")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    err = productCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&product)
    if err != nil {
        c.JSON(404, gin.H{"error": "Product not found"})
        return
    }
    
    cartCollection := utils.MongoDB.Collection("carts")
    var cart types.Cart
    
    err = cartCollection.FindOne(ctx, bson.M{"userid": user_id}).Decode(&cart)
    if err != nil {
        // cart doesn't exist, create new
        newCart := types.Cart{
            UserId:     user_id,
            Products:   []types.Product{product},
            TotalPrice: product.Price,
        }
        _, err := cartCollection.InsertOne(ctx, newCart)
        if err != nil {
            c.JSON(500, gin.H{"error": "Failed to create cart"})
            return
        }
        c.JSON(200, gin.H{"message": "Product added to new cart", "cart": newCart})
        return
    }
    
    // update existing cart
    cart.Products = append(cart.Products, product)
    cart.TotalPrice += product.Price
    _, err = cartCollection.UpdateOne(
        ctx,
        bson.M{"userid": user_id},  
        bson.M{"$set": bson.M{
            "products":    cart.Products,
            "totalprice": cart.TotalPrice, 
        }},
    )
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to update cart"})
        return
    }
    c.JSON(200, gin.H{"message": "Product added to cart", "cart": cart})
}
