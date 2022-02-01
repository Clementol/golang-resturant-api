package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	// db "github.com/Clementol/restur-manag/database"
	"github.com/Clementol/restur-manag/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var tableCollection *mongo.Collection = db.OpenCollection(db.Client, "table")

func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		result, err := orderCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			msg := "order occured while listing table items"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		var allTables []bson.M

		if err := result.All(ctx, &allTables); err != nil {
			log.Fatal(err.Error())
		}
		c.JSON(http.StatusOK, allTables)
	}
}

func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		tableId := c.Param("table_id")
		var table models.Table

		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": `error getting table item`})
			return
		}
		c.JSON(http.StatusOK, table)
	}
}

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var table models.Table

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(&table)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		table.ID = primitive.NewObjectID()
		table.Table_id = table.ID.Hex()

		result, insertErr := tableCollection.InsertOne(ctx, table)

		if insertErr != nil {
			msg := "table item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated, result)
	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var table models.Table

		tableId := c.Param("table_id")

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(&table)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		var updateObj primitive.D

		if table.Number_of_guest != nil {
			updateObj = append(updateObj, bson.E{Key: "number_of_guest", Value: table.Number_of_guest})
		}

		if table.Table_number != nil {
			updateObj = append(updateObj, bson.E{Key: "table_number", Value: table.Table_number})
		}

		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		filter := bson.M{"table_id": tableId}

		result, err := tableCollection.UpdateOne(
			ctx,
			filter, bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)
		if err != nil {
			msg := "table item update failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusAccepted, result)
	}
}
