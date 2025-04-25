package middlewares

import (
	"github.com/RohithBN/handlers"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context){
		authHeader:=c.GetHeader(("Authorization"))
		if authHeader==""{
			c.JSON(401,gin.H{"error":"Authorization header is missing"})
			c.Abort()
			return
		}
		token:=authHeader[len("Bearer "):]
		if token==""{
			c.JSON(401,gin.H{"error":"Token is missing"})
			c.Abort()
			return
		}
		claims, err := handlers.ValidateJWT(token)

		if err!=nil{
			c.JSON(401,gin.H{"error":"Invalid token"})
			c.Abort()
			return
		}
		c.Set("user", claims)
		c.Next()
	}

}
