package controllers

import (
	"context"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	db "github.com/Clementol/restur-manag/database"
	helper "github.com/Clementol/restur-manag/helpers"
	"github.com/Clementol/restur-manag/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection *mongo.Collection = db.OpenCollection(db.Client, "food")
var menuCollection *mongo.Collection = db.OpenCollection(db.Client, "menu")
var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}
		startIndex := (page - 1) * recordPerPage
		// startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		},
		}}

		projectStage := bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "food_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
		},
		}}

		result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})
		defer cancel()
		if err != nil {
			msg := "error occured while listing food items"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		var allFoods []bson.M
		if err = result.All(ctx, &allFoods); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allFoods[0])
	}
}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		foodId := c.Param("food_id")
		var food models.Food

		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": `error getting food item`})
			return
		}
		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var menu models.Menu
		var food models.Food

		file, err := c.FormFile("file")
		if err != nil {
			msg := "food image is required " + err.Error()
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		if err := c.Bind(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(food)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		err = menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
		defer cancel()
		if err != nil {
			msg := "menu was not found"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		imageFIle := helper.RenameFileName(*file)
		food.Food_image = imageFIle
		food.ID = primitive.NewObjectID()
		food.Food_id = food.ID.Hex()
		var num = toFixed(*food.Price, 2)
		food.Price = &num
		food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		result, insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr != nil {
			msg := "food item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated, result)
	}
}

func round(num float64) int {
	return int(num + math.Copysign(.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var food models.Food
		var menu models.Menu

		foodId := c.Param("food_id")
		file, _ := c.FormFile("file")

		if err := c.Bind(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var updateObj primitive.D

		if food.Name != nil {
			updateObj = append(updateObj, bson.E{Key: "name", Value: food.Name})
		}

		if food.Price != nil {
			updateObj = append(updateObj, bson.E{Key: "price", Value: food.Price})

		}

		if file != nil {
			imageFile := helper.RenameFileName(*file)
			updateObj = append(updateObj, bson.E{Key: "food_image", Value: imageFile})

		}

		if food.Menu_id != nil {
			err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
			defer cancel()
			if err != nil {
				msg := "message: Menu was not found"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "menu_id", Value: food.Menu_id})
		}

		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: food.Updated_at})

		filter := bson.M{"food_id": foodId}
		updatedFood := bson.M{}
		opt := options.FindOneAndUpdate().SetReturnDocument(options.After)

		err := foodCollection.FindOneAndUpdate(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			opt,
		).Decode(&updatedFood)

		if err != nil {
			msg := "food not found" + err.Error()
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, updatedFood)
	}
}
