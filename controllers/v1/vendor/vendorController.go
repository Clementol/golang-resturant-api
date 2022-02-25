package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

var vendorCollection = db.OpenCollection(db.Client, "vendor")
var orderCollection = db.OpenCollection(db.Client, "order")

var validate = validator.New()

// func GetVendors() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

// 		var vendors []primitive.M

// 		result, err := vendorCollection.Find(ctx, bson.M{})
// 		defer cancel()
// 		if err != nil {
// 			msg := "error getting vendors " + err.Error()
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
// 			return
// 		}
// 		if err := result.All(ctx, &vendors); err != nil {
// 			log.Fatal(err.Error())
// 		}
// 		c.JSON(http.StatusOK, vendors)
// 	}
// }

func GetVendor() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var vendor bson.M

		vendorId := c.Param("vendor_id")

		err := vendorCollection.FindOne(ctx, bson.M{"vendor_id": vendorId}).Decode(&vendor)
		defer cancel()
		if err != nil {
			msg := "Unable to get vendor" + err.Error()
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, vendor)

	}
}

func CreateVendor() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var vendor models.Vendor

		fileHeader, err := c.FormFile("file")

		if err != nil {
			msg := "vendor image is required"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		if err := c.Bind(&vendor); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(&vendor)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		if fileHeader != nil {
			imageText := helper.RenameFileName(*fileHeader)
			file, err := fileHeader.Open()
			if err != nil {
				log.Println(err.Error())
			}
			msg, url := helper.UploadToS3(imageText, file, fileHeader)

			if msg != "" {
				log.Println(msg)
			}
			if url != "" {
				vendor.Image = url
			}
		}

		vendor.ID = primitive.NewObjectID()
		vendor.Vendor_id = vendor.ID.Hex()

		vendor.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		vendor.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		result, err := vendorCollection.InsertOne(ctx, vendor)

		if err != nil {
			msg := "Unable to add vendor" + err.Error()
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateVendor() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var vendor models.Vendor

		vendorId := c.Param("vendor_id")
		defer cancel()
		if err := c.Bind(&vendor); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updateObj := bson.M{}

		if vendor.Name != "" {
			updateObj["name"] = vendor.Name
		}

		fileHeader, _ := c.FormFile("file")

		if fileHeader != nil {

			imageText := helper.RenameFileName(*fileHeader)
			file, err := fileHeader.Open()
			if err != nil {
				log.Println(err.Error())
			}
			msg, url := helper.UploadToS3(imageText, file, fileHeader)

			if msg != "" {
				log.Println(msg)
			}
			if url != "" {
				vendor.Image = url
			}
			updateObj["vendor_image"] = vendor.Image
		}

		if vendor.Location != "" {
			updateObj["location"] = vendor.Location
		}
		if vendor.Latitude != nil || vendor.Longitude != nil {
			updateObj["latitude"] = vendor.Latitude
			updateObj["longitude"] = vendor.Longitude

		}
		if vendor.Open_time != nil {
			updateObj["open_time"] = vendor.Open_time
		}
		if vendor.Close_time != nil {
			updateObj["close_time"] = vendor.Close_time
		}
		if vendor.Delivery_fee != nil {
			updateObj["delivery_fee"] = vendor.Delivery_fee
		}

		vendor.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj["updated_at"] = vendor.Updated_at

		filter := bson.M{"vendor_id": vendorId}
		updatedVendor := bson.M{}
		opt := options.FindOneAndUpdate().SetReturnDocument(options.After)

		err := vendorCollection.FindOneAndUpdate(ctx, filter,
			bson.M{"$set": updateObj},
			opt,
		).Decode(&updatedVendor)

		if err != nil {
			msg := "Unable to update vendor" + err.Error()
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusAccepted, updatedVendor)
	}
}

func GetVendorOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		vendorId := c.Param("vendor_id")

		matchVendorStage := bson.D{{Key: "$match", Value: bson.D{
			{Key: "vendor_id", Value: vendorId},
		},
		}}

		lookupCustomerStage := bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "user"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "user_id"},
			{Key: "as", Value: "user"},
		},
		}}

		lookupFoodStage := bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "food"},
			{Key: "localField", Value: "order_items.food_id"},
			{Key: "foreignField", Value: "food_id"},
			{Key: "as", Value: "order_foods"},
		},
		}}

		groupCustomerStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$user_id"},
			{Key: "customer_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		},
		}}

		projectOrderStage := bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "customer_count", Value: "$customer_count"},
			{Key: "customer", Value: "$data"},
		},
		}}

		unsetOrderStage := bson.D{{Key: "$unset", Value: []interface{}{
			"customer.order_items",
			"customer.user_id",
			"customer.user.token",
			"customer.user.refresh_token",
			"customer.user.password",
			"customer.user.created_at",
			"customer.user.updated_at",
			"customer.user.user_d",
			"customer.user._id",
		},
		}}

		var vendorOrders []bson.M

		result, err := orderCollection.Aggregate(ctx, mongo.Pipeline{
			matchVendorStage,
			lookupCustomerStage,
			lookupFoodStage,
			groupCustomerStage,
			projectOrderStage,
			unsetOrderStage,
		})

		if err != nil {
			msg := "unable to get vendor orders"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if err := result.All(ctx, &vendorOrders); err != nil {
			log.Fatal(err.Error())
		}
		defer cancel()

		if vendorOrders == nil {
			msg := "vendor doesn't have orders"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, vendorOrders)
	}
}

func UpdateVendorOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		// var order models.Order

		var updateOrder models.UpdateOrder
		vendorId := c.Param("vendor_id")

		updateObj := bson.M{}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := c.BindJSON(&updateOrder); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validateErr := validate.Struct(&updateOrder)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
			return
		}

		// log.Fatal(updateOrder.Order_ids)

		updateObj["order_status"] = updateOrder.Order_status

		updateObj["updated_at"], _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		idsToUpdate := bson.M{
			"vendor_id": vendorId,
			"order_id": bson.M{
				"$in": updateOrder.Order_ids,
			},
		}

		fieldToUpdate := bson.M{"$set": updateObj}

		result, err := orderCollection.UpdateMany(ctx, idsToUpdate, fieldToUpdate)

		if err != nil {
			msg := "order item update failed " + err.Error()
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		msg := fmt.Sprint(result.ModifiedCount) + " order(s) updated"
		c.JSON(http.StatusOK, msg)
	}
}
