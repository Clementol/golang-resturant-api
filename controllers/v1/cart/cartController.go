package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	db "github.com/Clementol/restur-manag/database"
	"github.com/Clementol/restur-manag/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cartCollection (*mongo.Collection) = db.OpenCollection(db.Client, "cart")
var validate = validator.New()

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
		cart.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		cart.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		var filter primitive.M
		userCart := bson.M{}

		filter = bson.M{
			"user_id": cart.User_id,
		}
		for _, cartItem := range cart.Cart_items {
			var update primitive.M

			checkForItem := bson.M{
				"user_id":            cart.User_id,
				"cart_items.food_id": cartItem.Food_id,
			}
			userCounts, err := cartCollection.CountDocuments(ctx, filter)
			if err != nil {
				log.Fatal(err.Error())
			}

			if userCounts == 1 {

				cartCounts, _ := cartCollection.CountDocuments(ctx, checkForItem)

				if cartCounts == 1 { // if cart item exist

					checkForItem = bson.M{
						"cart_items.food_id": bson.M{"$eq": cartItem.Food_id},
					}
					update := bson.M{"$set": bson.M{
						"cart_items.$.quantity": cartItem.Quantity,
						"updated_at":            cart.Updated_at,
					},
					}

					opt := options.Update()
					updatedCart, err := cartCollection.UpdateOne(ctx, checkForItem, update, opt)
					if err != nil {
						msg := "error adding cart " + err.Error()
						c.JSON(http.StatusBadRequest, gin.H{"error": msg})
						return
					}
					if updatedCart.MatchedCount != 1 {
						msg := "no cart item found"
						c.JSON(http.StatusBadRequest, gin.H{"error": msg})
						return
					}
					c.JSON(http.StatusAccepted, updatedCart.ModifiedCount)
					return

				} else {
					update = bson.M{
						"$push": bson.M{
							"cart_items": bson.M{
								"food_id":  cartItem.Food_id,
								"quantity": cartItem.Quantity,
							},
						},
						"$set": bson.M{
							"updated_at": cart.Updated_at,
						},
					}

					opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
					err := cartCollection.FindOneAndUpdate(ctx, filter, update, opt).Decode(&userCart)
					if err != nil {
						msg := "error adding cart"
						c.JSON(http.StatusBadRequest, gin.H{"error": msg})
						return
					}

					c.JSON(http.StatusAccepted, userCart)
					return
				}
			}

		}
		newCart := bson.M{}
		newCart["user_id"] = cart.User_id
		newCart["cart_items"] = cart.Cart_items
		newCart["created_at"] = cart.Created_at
		newCart["updated_at"] = cart.Updated_at

		result, insertErr := cartCollection.InsertOne(ctx, newCart)
		if insertErr != nil {
			msg := "error adding cart "
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusCreated, result)

	}
}

func GetCartItems() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.MustGet("user_id").(string)
		cartItems := []bson.M{}

		matchCartStage := bson.D{{Key: "$match", Value: bson.D{
			{Key: "user_id", Value: userId},
		},
		}}

		projectCartStage := bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "user_id", Value: 0},
			{Key: "created_at", Value: 0},
			{Key: "updated_at", Value: 0},
		},
		}}

		result, err := cartCollection.Aggregate(ctx, mongo.Pipeline{
			matchCartStage,
			projectCartStage,
		})

		if err != nil {
			msg := "unable to get cart items "
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if err := result.All(ctx, &cartItems); err != nil {
			log.Fatal(err.Error())
		}

		c.JSON(http.StatusOK, cartItems)

	}
}

func RemoveCartItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.MustGet("user_id").(string)

		var cart models.RemoveCartItem

		if err := c.Bind(&cart); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(&cart)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		cart.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		checkForItem := bson.M{
			"user_id":            userId,
			"cart_items.food_id": cart.Food_id,
		}

		update := bson.M{
			"$pull": bson.M{
				"cart_items": bson.M{
					"food_id": cart.Food_id,
					// "quantity": cartItem.Quantity,
				},
			},
			"$set": bson.M{
				"updated_at": cart.Updated_at,
			},
		}

		result, err := cartCollection.UpdateOne(ctx, checkForItem, update)

		if err != nil {
			msg := "unable to remove cart item"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusAccepted, result.ModifiedCount)

	}
}
