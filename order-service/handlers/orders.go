package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/RohithBN/order-service/kafka"
	"github.com/RohithBN/shared/redis"
	"github.com/RohithBN/shared/types"
	"github.com/RohithBN/shared/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TwoFacAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDstr := c.GetHeader("X-User-ID")
		if userIDstr == "" {
			c.JSON(401, gin.H{"error": "User ID not found"})
			c.Abort()
			return
		}
		userID, err := strconv.Atoi(userIDstr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid user ID format"})
			c.Abort()
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		userVerified, err := redis.RedisClient.Get(ctx, fmt.Sprintf("user:%d:verified", userID)).Result()
		if err != nil {
			log.Printf("Redis error: %v", err)
			c.JSON(401, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}
		if userVerified != "true" {
			c.JSON(401, gin.H{"error": "User not verified"})
			c.Abort()
			return
		}
		c.Next()
	}

}

func VerifyOTP(c *gin.Context) {

	userIDstr := c.GetHeader("X-User-ID")
	if userIDstr == "" {
		c.JSON(401, gin.H{"error": "User ID not found"})
		return
	}
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID format"})
		return
	}
	// get the otp from req
	var otpData struct {
		OTP string `json:"otp"`
	}
	if err := c.BindJSON(&otpData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}
	otp := otpData.OTP

	userEmail := c.GetHeader("X-User-Email")
	if userEmail == "" {
		c.JSON(401, gin.H{"error": "User email not found"})
		return
	}

	//get the otp for the user email from redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	otpFromRedis, err := redis.RedisClient.Get(ctx, userEmail).Result()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get OTP from redis"})
		return
	}
	if otpFromRedis == "" {
		c.JSON(401, gin.H{"error": "OTP not found"})
		return
	}
	//compare the otp from req and otp from redis
	if otpFromRedis != otp {
		c.JSON(401, gin.H{"error": "Invalid OTP"})
		return
	}
	//if otp valid , set flag of user in redis as true (verified)
	err = redis.RedisClient.Set(ctx, fmt.Sprintf("user:%d:verified", userID), "true", 30*time.Minute).Err()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to set user verified flag in redis"})
		return
	}
	//delete the otp from redis
	err = redis.RedisClient.Del(ctx, userEmail).Err()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete OTP from redis"})
		return
	}

	c.JSON(200, gin.H{"message": "OTP verified successfully"})
}

func SendOTP(c *gin.Context) {
	//get the user email
	userEmail := c.GetHeader("X-User-Email")
	if userEmail == "" {
		c.JSON(401, gin.H{"error": "User email not found"})
		return
	}
	fmt.Println("Sending OTP to email:", userEmail)
	err := kafka.VerifyOTPEmailProducer(userEmail)
	if err != nil {
		fmt.Println("Error Sending verify mail")
	}
}

func CreateOrder(c *gin.Context) {
	userIdStr := c.GetHeader("X-User-ID")
	if userIdStr == "" {
		c.JSON(401, gin.H{"error": "User ID not found"})
		return
	}

	user_id, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID format"})
		return
	}
	userEmail := c.GetHeader("X-User-Email")
	if userEmail == "" {
		c.JSON(401, gin.H{"error": "User email not found"})
		return
	}

	// Get user's cart
	cartCollection := utils.MongoDB.Collection("carts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cart types.Cart
	err = cartCollection.FindOne(ctx, bson.M{"userid": user_id}).Decode(&cart)
	if err != nil {
		c.JSON(404, gin.H{"error": "Cart not found"})
		return
	}

	// Create new order
	order := types.Order{
		UserId:     user_id,
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
	_, err = cartCollection.DeleteOne(ctx, bson.M{"userid": user_id})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to clear cart"})
		return
	}

	order.OrderId = result.InsertedID.(primitive.ObjectID)
	err = utils.SendOrderConfirmationEmail(userEmail, &order)
	if err != nil {
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
		OrderId string  `json:"order_id"`
		Amount  float64 `json:"amount"`
		Method  string  `json:"method"` // "stripe" or "paypal"
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
		"status":  "paid",
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
		"status":  updateInfo.Status,
	})
}

func GetOrders(c *gin.Context) {
	userIdStr := c.GetHeader("X-User-ID")
	if userIdStr == "" {
		c.JSON(401, gin.H{"error": "User ID not found"})
		return
	}

	user_id, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID format"})
		return
	}

	orderCollection := utils.MongoDB.Collection("orders")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := orderCollection.Find(ctx, bson.M{"userid": user_id})
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
