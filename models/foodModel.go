package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Food struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       *string            `form:"name" validate:"required,min=2,max=100"`
	Price      *float64           `form:"price" validate:"required"`
	Food_image string             `bson:"food_image"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
	Food_id    string             `json:"food_id"`
	Menu_id    *string            `form:"menu_id" validate:"required"`
}
