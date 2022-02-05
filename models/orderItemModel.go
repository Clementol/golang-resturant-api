package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderItems struct {
	ID            primitive.ObjectID `bson:"_id"`
	Quantity      *int32             `json:"quantity" validate:"required"`
	Unit_price    *float64           `json:"unit_price" validate:"required"`
	Food_id       *string            `json:"food_id" validate:"required"`
	Order_item_id string             `bson:"order_item_id" validate:"required"`
}
