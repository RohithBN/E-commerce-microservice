package handlers

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/RohithBN/types"
	"github.com/RohithBN/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var Users []types.User

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
var dbPool *pgxpool.Pool

func InitDB() error {
    var err error
    dbPool, err = utils.GetDB()
    if err != nil {
        return fmt.Errorf("failed to initialize database: %v", err)
    }
    return nil
}

func Register(c *gin.Context) {
	var user types.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now().Format(time.RFC3339)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to connect to database"})
		return
	}
	_, err = dbPool.Exec(
		c,
		`INSERT INTO USERS (name, email, password, created_at) VALUES ($1, $2, $3, $4)`,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to register user",
			"details": err.Error()})

		return
	}
	Users = append(Users, user)

	c.JSON(200, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

func Login(c *gin.Context) {
	var user types.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var err error
	
	// Check if the user exists in the database
	var password string
	err=dbPool.QueryRow(
		c,
		`SELECT id,password FROM USERS WHERE email = $1`,
		user.Email,
	).Scan(&user.Id,&password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(401, gin.H{"error": "User not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to query user",
			"details": err.Error()})
		return
	}
	// Simulate a successful login

	isValidPassword := bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
	if isValidPassword != nil {
		c.JSON(401, gin.H{"error": "Invalid password"})
		return
	}
	tokenString, err := GenerateJWT(user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token",
			"details": err.Error()})
		return
	}
	fmt.Println("Token:", tokenString)
	c.JSON(200, gin.H{
		"message": "User logged in successfully",
		"user":    user,
		"token":   tokenString,
	})

}

func GenerateJWT(user types.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    user.Id,
			"name":  user.Name,
			"email": user.Email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token.Claims, nil
}

