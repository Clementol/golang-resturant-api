package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	db "github.com/Clementol/restur-manag/database"
	"github.com/Clementol/restur-manag/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cartCollection (*mongo.Collection) = db.OpenCollection(db.Client, "cart")

func AddItemToCart() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var cart models.Cart
		userId := c.MustGet("user_id").(string)

		cart.User_id = userId

		if err := c.Bind(&cart); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(&cart)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		cartObj := bson.M{}
		cart.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		cart.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		cartObj["created_at"] = cart.Created_at
		cartObj["updated_at"] = cart.Updated_at
		cartObj["cart_items"] = cart.Cart_items

		userCart := bson.M{}

		for _, cartItem := range cart.Cart_items {
			var update primitive.M

			filter := bson.M{
				"user_id": cart.User_id,
			}
			checkForItem := bson.M{
				"user_id":            cart.User_id,
				"cart_items.food_id": cartItem.Food_id,
			}
			cartCounts, _ := cartCollection.CountDocuments(ctx, checkForItem)
			log.Fatal(cartCounts)
			if cartCounts == 1 {
				update = bson.M{
					"$set": bson.M{
						"cart_items.$": cartItem,
						"updated_at":   cart.Updated_at,
					},
				}
				opt := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
				err := cartCollection.FindOneAndUpdate(ctx, filter, update, opt).Decode(&userCart)
				if err != nil {
					msg := "error adding cart " + err.Error()
					c.JSON(http.StatusBadRequest, gin.H{"error": msg})
					return
				}
			} else {
				update = bson.M{
					"$push": bson.M{
						"cart_items": cartItem,
					},
				}
				opt := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
				err := cartCollection.FindOneAndUpdate(ctx, filter, update, opt).Decode(&userCart)
				if err != nil {
					msg := "error adding cart " + err.Error()
					c.JSON(http.StatusBadRequest, gin.H{"error": msg})
					return
				}
			}

		}

		c.JSON(http.StatusAccepted, userCart)

	}
}

func GetCartItems() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func RemoveCartItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
