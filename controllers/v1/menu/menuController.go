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

var menuCollection *mongo.Collection = db.OpenCollection(db.Client, "menu")
var validate = validator.New()

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		result, err := menuCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			msg := "error occur while listing the menu item"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		var allMenus []bson.M
		if err = result.All(ctx, &allMenus); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenus)
	}
}

func GetMenusWithFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		// defer cancel()
		vendorId := c.Param("vendor_id")

		matchMenuStage := bson.D{{Key: "$match", Value: bson.D{
			{Key: "vendor_id", Value: vendorId},
		},
		}}

		lookupFoodStage := bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "food"},
			{Key: "localField", Value: "menu_id"},
			{Key: "foreignField", Value: "menu_id"},
			{Key: "as", Value: "foods"},
		},
		}}

		groupMenuStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		},
		}}

		projectMenuStage := bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "menus", Value: "$data"},
		},
		}}

		foodMenus := []bson.M{}

		menuItems, _ := menuCollection.CountDocuments(ctx, bson.M{"vendor_id": vendorId})

		defer cancel()
		if menuItems == 0 {
			msg := "menu items not found "
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		result, err := menuCollection.Aggregate(
			ctx, mongo.Pipeline{
				matchMenuStage,
				lookupFoodStage,
				groupMenuStage,
				projectMenuStage,
			})
		defer cancel()

		if err != nil {
			msg := "could not get vendor food menus" + err.Error()
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		if err := result.All(ctx, &foodMenus); err != nil {
			log.Fatal(err.Error())
		}
		defer cancel()

		c.JSON(http.StatusOK, foodMenus)
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		menuId := c.Param("menu_id")
		var menu bson.M

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		defer cancel()
		if err != nil {
			msg := `error getting menu item ` + err.Error()
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(menu)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result, insertErr := menuCollection.InsertOne(ctx, menu)

		if insertErr != nil {
			msg := "menu item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated, result)
		// defer cancel()
	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return start.After(time.Now()) && end.After(start)
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		menuId := c.Param("menu_id")
		filter := bson.M{"menu_id": menuId}

		var updateObj primitive.D

		if menu.Start_Date != nil && menu.End_Date != nil {
			if !inTimeSpan(*menu.Start_Date, *menu.End_Date, time.Now()) {
				msg := "Kindly retype the time"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				defer cancel()
				return
			}
			updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.Start_Date})
			updateObj = append(updateObj, bson.E{Key: "end_date", Value: menu.End_Date})

		}
		if menu.Name != "" {
			updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
		}

		// if menu.Category != "" {
		// 	updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Category})
		// }

		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.Updated_at})

		opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
		updatedMenu := bson.M{}

		err := menuCollection.FindOneAndUpdate(
			ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, opt,
		).Decode(&updatedMenu)

		if err != nil {
			msg := "menu update failed" + err.Error()
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		defer cancel()
		c.JSON(http.StatusAccepted, updatedMenu)
	}
}
