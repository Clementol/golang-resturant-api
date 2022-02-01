package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vendor struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string             `form:"name" validate:"required,min=2,max=30"`
	Image        string             `bson:"vendor_image"`
	Location     string             `form:"location" validate:"required"`
	Latitude     *float64           `form:"latitude" validate:"required"`
	Longitude    *float64           `form:"longitude" validate:"required"`
	Open_time    *float64           `form:"open_time" validate:"required"`
	Close_time   *float64           `form:"close_time" validate:"required"`
	Delivery_fee *float64           `form:"delivery_fee" validate:"required"`
	Created_at   time.Time          `json:"created_at"`
	Updated_at   time.Time          `json:"updated_at"`
	Vendor_id    string             `json:"vendor_id"`
}
