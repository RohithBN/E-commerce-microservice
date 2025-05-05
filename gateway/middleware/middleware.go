package middleware

import (
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

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

type ClientData struct {
	tokensRemaining float64
	lastRefillTime  int64
	mu              sync.Mutex // Per-client mutex for better concurrency
}

var (
    clientDataMap      = make(map[string]*ClientData)
    clientDataMapMutex sync.RWMutex
    maxBucketSize      = 5.0
    refillRate         = 5.0 / 60.0 // 5 tokens per minute = 0.0833 tokens per second
    cleanupInterval    = 10 * time.Minute
)

func init() {
	// Cleanup routine for inactive users
	go func() {
		for {
			time.Sleep(cleanupInterval)
			cleanupInactiveClients()
		}
	}()
}

func cleanupInactiveClients() {
    //represents the last 30 min from current time
	threshold := time.Now().Add(-30*time.Minute).UnixNano() / 1e6

	clientDataMapMutex.Lock()
	defer clientDataMapMutex.Unlock()

	for userId, client := range clientDataMap {
		client.mu.Lock()
        //if users last req is older than threshold , delete the user from map
		if client.lastRefillTime < threshold {
			delete(clientDataMap, userId)
		}
		client.mu.Unlock()
	}
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetHeader("X-User-ID")
		if userId == "" {
			c.Next() 
			return
		}
		clientDataMapMutex.RLock()
        //check if client exists
		client, exists := clientDataMap[userId]
		clientDataMapMutex.RUnlock()

		if !exists {
			clientDataMapMutex.Lock()
			client, exists = clientDataMap[userId]
			if !exists {
                // if client doesnt exist then create a new client , with current time as last req
				client = &ClientData{
					tokensRemaining: maxBucketSize,
					lastRefillTime:  time.Now().UnixNano() / 1e6,
				}
				clientDataMap[userId] = client
			}
			clientDataMapMutex.Unlock()
		}

		// Lock this client's data for the duration of our check/update
		client.mu.Lock()
		defer client.mu.Unlock()

		// Calculate token refill based on time elapsed
		currentTime := time.Now().UnixNano() / 1e6
        //how many sec diff btwn now nad last req 
		elapsedSeconds := float64(currentTime-client.lastRefillTime) / 1000.0
        //twlls how many tokens to add based on refillrate
		tokensToAdd := elapsedSeconds * refillRate

		if tokensToAdd > 0 {
            //if there are tokens to add , add it to remaining tokens , but ensure using min() that it doesnt exceed the max bucket size
			client.tokensRemaining = math.Min(maxBucketSize, client.tokensRemaining+tokensToAdd)
			client.lastRefillTime = currentTime // last req=now
		}

		if client.tokensRemaining >= 1.0 {
			// Consume one token
			client.tokensRemaining -= 1.0
			c.Next()
		} else {
			c.JSON(429, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": int(math.Ceil((1.0 - client.tokensRemaining) / refillRate)),
			})
			c.Abort()
		}
	}
}
