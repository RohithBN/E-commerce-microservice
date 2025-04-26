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

func CreateOrder(c *gin.Context) {
    user := c.MustGet("user").(*types.User)
    
    // Get user's cart
    cartCollection := utils.MongoDB.Collection("carts")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    var cart types.Cart
    err := cartCollection.FindOne(ctx, bson.M{"userid": user.Id}).Decode(&cart)
    if err != nil {
        c.JSON(404, gin.H{"error": "Cart not found"})
        return
    }

    // Create new order
    order := types.Order{
        UserId:     user.Id,
        Products:   cart.Products,
        TotalPrice: cart.TotalPrice,
        Status:     "pending",
        CreatedAt:  time.Now().Format(time.RFC3339),
    }

    orderCollection := utils.MongoDB.Collection("orders")
    result, err := orderCollection.InsertOne(ctx, order)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to create order"})
        return
    }

    // Clear the cart after order creation
    _, err = cartCollection.DeleteOne(ctx, bson.M{"userid": user.Id})
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to clear cart"})
        return
    }

	order.OrderId = result.InsertedID.(primitive.ObjectID)
	err=utils.SendOrderConfirmationEmail(user.Email,&order)
	if err!=nil{
		c.JSON(500, gin.H{"error": "Failed to send order confirmation email"})
		return
	}

    
    c.JSON(200, gin.H{
        "message": "Order created successfully",
        "order":   order,
    })
}

func ProcessPayment(c *gin.Context) {
    var paymentInfo struct {
        OrderId string `json:"order_id"`
        Amount  float64 `json:"amount"`
        Method  string `json:"method"` // "stripe" or "paypal"
    }

    if err := c.BindJSON(&paymentInfo); err != nil {
        c.JSON(400, gin.H{"error": "Invalid payment info"})
        return
    }

    // Simulate payment processing
    time.Sleep(1 * time.Second)

    // Update order status to paid
    orderId, err := primitive.ObjectIDFromHex(paymentInfo.OrderId)
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid order ID"})
        return
    }

    orderCollection := utils.MongoDB.Collection("orders")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err = orderCollection.UpdateOne(
        ctx,
        bson.M{"_id": orderId},
        bson.M{"$set": bson.M{"status": "paid"}},
    )
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to update order status"})
        return
    }

    c.JSON(200, gin.H{
        "message": "Payment processed successfully",
        "status": "paid",
    })
}

func UpdateOrderStatus(c *gin.Context) {
    orderId := c.Param("orderId")
    var updateInfo struct {
        Status string `json:"status"`
    }

    if err := c.BindJSON(&updateInfo); err != nil {
        c.JSON(400, gin.H{"error": "Invalid status update"})
        return
    }

    // Validate status
    validStatuses := map[string]bool{
        "pending": true,
        "paid":    true,
        "shipped": true,
    }

    if !validStatuses[updateInfo.Status] {
        c.JSON(400, gin.H{"error": "Invalid status"})
        return
    }

    objectId, err := primitive.ObjectIDFromHex(orderId)
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid order ID"})
        return
    }

    orderCollection := utils.MongoDB.Collection("orders")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err = orderCollection.UpdateOne(
        ctx,
        bson.M{"_id": objectId},
        bson.M{"$set": bson.M{"status": updateInfo.Status}},
    )
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to update order status"})
        return
    }

    c.JSON(200, gin.H{
        "message": "Order status updated successfully",
        "status": updateInfo.Status,
    })
}

func GetOrders(c *gin.Context) {
    user := c.MustGet("user").(*types.User)
    
    orderCollection := utils.MongoDB.Collection("orders")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cursor, err := orderCollection.Find(ctx, bson.M{"userid": user.Id})
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to fetch orders"})
        return
    }
    defer cursor.Close(ctx)

    var orders []types.Order
    if err = cursor.All(ctx, &orders); err != nil {
        c.JSON(500, gin.H{"error": "Failed to decode orders"})
        return
    }

    c.JSON(200, gin.H{"orders": orders})
}