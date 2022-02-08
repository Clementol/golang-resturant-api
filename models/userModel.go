package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	First_Name    *string            `form:"first_name" validate:"required,min=2,max=100"`
	Last_Name     *string            `form:"last_name" validate:"required,min=2,max=100"`
	Password      *string            `form:"password" validate:"required,min=6"`
	Email         *string            `form:"email" validate:"email,required"`
	Avatar        *string            `json:"avatar"`
	Phone         *string            `form:"phone" validate:"required"`
	Token         *string            `json:"token"`
	Refresh_Token *string            `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"user_id"`
}
