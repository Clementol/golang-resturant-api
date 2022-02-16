package models

import (
	"time"
)

type Cart struct {
	// ID         primitive.ObjectID `bson:"_id"`
	User_id    string      `json:"user_id" validate:"required"`
	Cart_items []CartItems `bson:"cart_items" json:"cart_items" validate:"required"`
	// Cart_id    string             `bson:"cart_id" validate:"required"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

type CartItems struct {
	Food_id  string `json:"food_id" validate:"required"`
	Quantity *int32 `json:"quantity" validate:"required,min=1"`
}

type RemoveCartItem struct {
	Food_id  string `json:"food_id" validate:"required"`
	Updated_at time.Time `json:"updated_at"`
}