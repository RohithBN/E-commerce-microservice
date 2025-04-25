package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type Product struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name"`
	Price       float64            `json:"price"`
	Description string             `json:"description"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
	AddedToCart bool               `json:"added_to_cart"`
	Category    string             `json:"category"`
	Stock       int                `json:"stock"`
}

type Cart struct {
	UserId   int      `json:"user_id"`
	Products   []Product `json:"products"`
	TotalPrice float64   `json:"total_price"`
}

type Order struct {
	UserId       User      `json:"user_id"`
	OrderId      primitive.ObjectID `json:"order_id" bson:"_id,omitempty"`
	Products   []Product `json:"products"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  string    `json:"created_at"`
	Status     string    `json:"status"`
}
