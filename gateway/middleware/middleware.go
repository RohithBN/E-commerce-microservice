package middleware

import (
	"fmt"
	"os"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Add debug logging
        fmt.Printf("JWT_SECRET_KEY: %s\n", os.Getenv("JWT_SECRET_KEY"))
        
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Authorization header is missing"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            secretKey := os.Getenv("JWT_SECRET_KEY")
            if secretKey == "" {
                return nil, fmt.Errorf("JWT_SECRET_KEY not found in environment")
            }
            return []byte(secretKey), nil
        })

        if err != nil {
            c.JSON(401, gin.H{"error": fmt.Sprintf("Token validation error: %v", err)})
            c.Abort()
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            // Forward user info in headers
            c.Request.Header.Set("X-User-ID", fmt.Sprintf("%.0f", claims["id"].(float64)))
            c.Request.Header.Set("X-User-Email", claims["email"].(string))
            
            c.Next()
        } else {
            c.JSON(401, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }
    }
}

