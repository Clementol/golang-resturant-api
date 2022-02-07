package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID             primitive.ObjectID `bson:"_id"`
	Order_Date     time.Time          `json:"order_date"`
	Created_at     time.Time          `json:"created_at"`
	Updated_at     time.Time          `json:"updated_at"`
	Order_id       string             `json:"order_id"`
	User_id        string             `json:"user_id" validate:"required"`
	Vendor_id      string             `json:"vendor_id" validate:"required"`
	Total_amount   float64            `json:"total_amount" validate:"required" `
	Payment_status *string            `json:"payment_status" validate:"eq=PENDING|eq=PAID"`
	Payment_method *string            `json:"payment_method" validate:"eq=CARD|eq=CASH|eq="`
	Order_status   string             `json:"order_status" validate:"eq=ORDERED|eq=COMING|eq=DELIVERED"`
	OrderItems     []OrderItems       `bson:"order_items" json:"order_items" validate:"required"`
}
