package middlewares

import (
	"github.com/RohithBN/handlers"
	"github.com/RohithBN/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader(("Authorization"))
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Authorization header is missing"})
            c.Abort()
            return
        }
        token := authHeader[len("Bearer "):]
        if token == "" {
            c.JSON(401, gin.H{"error": "Token is missing"})
            c.Abort()
            return
        }
        claims, err := handlers.ValidateJWT(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        mapClaims, ok := claims.(jwt.MapClaims)
        if !ok {
            c.JSON(401, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }
        // Safe extraction with type assertions
        idFloat, ok := mapClaims["id"].(float64)
        if !ok {
            c.JSON(400, gin.H{"error": "Invalid user ID in token"})
            c.Abort()
            return
        }
        
        // Extract name and email as strings directly
        name, _ := mapClaims["name"].(string)
        email, _ := mapClaims["email"].(string)
        
        // Create a pointer to a User instead of a User value
        user := &types.User{
            Id:    int(idFloat), // Convert float64 to int safely
            Name:  name,
            Email: email,
        }
        c.Set("user", user)
        c.Next()
    }
}
