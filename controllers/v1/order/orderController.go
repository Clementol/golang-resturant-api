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
)

var orderCollection = (*mongo.Collection)(db.OpenCollection(db.Client, "order"))
var vendorCollection = (*mongo.Collection)(db.OpenCollection(db.Client, "vendor"))

var validate = validator.New()

func GetUserOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.MustGet("user_id").(string)

		matchOrderStage := bson.D{{Key: "$match", Value: bson.D{
			{Key: "user_id", Value: userId},
		},
		}}

		lookupFoodStage := bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "food"},
			{Key: "foreignField", Value: "food_id"},
			{Key: "localField", Value: "order_items.food_id"},
			{Key: "as", Value: "order_foods"},
		},
		}}

		lookupVendorStage := bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "vendor"},
			{Key: "foreignField", Value: "vendor_id"},
			{Key: "localField", Value: "vendor_id"},
			{Key: "as", Value: "vendor"},
		},
		}}

		groupOrderStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "_id", Value: nil}}},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		},
		}}

		projectOrderStage := bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: "$total_count"},
			{Key: "orders", Value: "$data"},
		},
		}}

		unsetOrderStage := bson.D{{Key: "$unset", Value: []interface{}{
			"orders.order_items",
			"orders.vendor_id",
		},
		}}

		var allOrders []bson.M
		result, err := orderCollection.Aggregate(ctx, mongo.Pipeline{
			matchOrderStage,
			lookupFoodStage,
			lookupVendorStage,
			groupOrderStage,
			projectOrderStage,
			unsetOrderStage,
		})

		if err != nil {
			msg := "could not get orders " + err.Error()
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		if err := result.All(ctx, &allOrders); err != nil {
			log.Fatal(err.Error())
		}
		// result, err := orderCollection.Find(context.TODO(), bson.M{})
		defer cancel()

		c.JSON(http.StatusOK, allOrders)
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var order models.Order

		userId := c.MustGet("user_id").(string)
		if userId == "" {
			msg := "unable to get user"
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}
		order.User_id = userId

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		defer cancel()

		validationErr := validate.Struct(&order)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		vendorCount, err := vendorCollection.CountDocuments(ctx, bson.M{"vendor_id": order.Vendor_id})

		if err != nil {
			msg := "error occurred while checking for vendor"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		if vendorCount == 0 {
			msg := "vendor doesn't exist"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		order.ID = primitive.NewObjectID()
		order.Order_id = order.ID.Hex()
		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		var orderObj primitive.D
		var orderItemsObj []interface{}

		orderObj = append(orderObj, primitive.E{Key: "_id", Value: order.ID})
		orderObj = append(orderObj, primitive.E{Key: "order_date", Value: order.Order_Date})
		orderObj = append(orderObj, primitive.E{Key: "created_at", Value: order.Created_at})
		orderObj = append(orderObj, primitive.E{Key: "updated_at", Value: order.Updated_at})
		orderObj = append(orderObj, primitive.E{Key: "order_id", Value: order.Order_id})
		orderObj = append(orderObj, primitive.E{Key: "user_id", Value: order.User_id})
		orderObj = append(orderObj, primitive.E{Key: "vendor_id", Value: order.Vendor_id})
		orderObj = append(orderObj, primitive.E{Key: "total_amount", Value: order.Total_amount})
		orderObj = append(orderObj, primitive.E{Key: "payment_status", Value: order.Payment_status})
		orderObj = append(orderObj, primitive.E{Key: "payment_method", Value: order.Payment_method})
		orderObj = append(orderObj, primitive.E{Key: "order_status", Value: order.Order_status})

		for _, orderItem := range order.OrderItems {

			orderItem.ID = primitive.NewObjectID()
			orderItem.Order_item_id = orderItem.ID.Hex()
			orderItemsObj = append(orderItemsObj, orderItem)
		}
		orderObj = append(orderObj, primitive.E{Key: "order_items", Value: orderItemsObj})

		result, insertErr := orderCollection.InsertOne(ctx, orderObj)

		if insertErr != nil {
			msg := "Order was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusAccepted, result)
	}

}

func OrderItemOrderCreator(order models.Order) string {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	orderCollection.InsertOne(ctx, order)
	defer cancel()
	return order.Order_id
}
