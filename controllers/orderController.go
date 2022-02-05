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

var orderCollection = (*mongo.Collection)(db.OpenCollection(db.Client, "order"))
var tableCollection = (*mongo.Collection)(db.OpenCollection(db.Client, "table"))

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		result, err := orderCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			msg := "order occured while listing order items"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		var allOrders []bson.M

		if err := result.All(ctx, &allOrders); err != nil {
			log.Fatal(err.Error())
		}
		c.JSON(http.StatusOK, allOrders)
	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		orderId := c.Param("order_id")
		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": `error getting order item`})
			return
		}
		c.JSON(http.StatusOK, order)
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		defer cancel()

		validationErr := validate.Struct(&order)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		

		order.ID = primitive.NewObjectID()
		order.Order_id = order.ID.Hex()
		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		var orderObj  primitive.D
		var orderItemsObj  []interface{}

		orderObj = append(orderObj, primitive.E{Key: "_id", Value: order.ID})
		orderObj = append(orderObj, primitive.E{Key: "order_date", Value: order.Order_Date})
		orderObj = append(orderObj, primitive.E{Key: "created_at", Value: order.Created_at})
		orderObj = append(orderObj, primitive.E{Key: "updated_at", Value: order.Updated_at})
		orderObj = append(orderObj, primitive.E{Key: "order_id", Value: order.Order_id})
		orderObj = append(orderObj, primitive.E{Key: "user_id", Value: order.User_id})
		orderObj = append(orderObj, primitive.E{Key: "total_amount", Value: order.Total_amount})
		orderObj = append(orderObj, primitive.E{Key: "payment_status", Value: order.Payment_status})
		orderObj = append(orderObj, primitive.E{Key: "payment_method", Value: order.Payment_method})
		orderObj = append(orderObj, primitive.E{Key: "order_status", Value: order.Order_status})


		for _, orderItem := range order.OrderItems {


			// validationErr := validate.Struct(orderItem)
			// if validationErr != nil {
			// 	c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			// 	return
			// }
			orderItem.ID = primitive.NewObjectID()
			orderItem.Order_item_id = orderItem.ID.Hex()
			orderItemsObj = append(orderItemsObj, orderItem)
		}
		orderObj = append(orderObj, primitive.E{Key: "order_items", Value: orderItemsObj})

		// for _, orderItem := range order.OrderItems {
			
		// 	log.Println(orderItem.Order_item_id, "here")
		// }

		// log.Fatal()
		
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

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var table models.Table
		var order models.Order

		var updateObj primitive.D
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		orderId := c.Param("order_id")

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if order.User_id != "" {
			err := orderCollection.FindOne(ctx, bson.M{"table_id": order.User_id}).Decode(&table)
			defer cancel()
			if err != nil {
				msg := "message: Menu was not found"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "table_id", Value: order.User_id})

		}
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: order.Updated_at})
		upsert := true
		filter := bson.M{"order_id": orderId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, err := orderCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)
		if err != nil {
			msg := "order item update failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
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
