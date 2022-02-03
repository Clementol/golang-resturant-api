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

type InvoiceViewFormat struct {
	Invoice_id       string
	Payment_method   string
	Order_id         string
	Payment_status   *string
	Payment_due      interface{}
	Table_number     interface{}
	Payment_due_date time.Time
	Order_details    interface{}
}

// var validate = validator.New()
var invoiceCollection *mongo.Collection = db.OpenCollection(db.Client, "invoice")

// var orderCollection *mongo.Collection = db.OpenCollection(db.Client, "order")

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		result, err := invoiceCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			msg := "order occured while listing order items"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		var allInvoices []primitive.M
		if err = result.All(ctx, &allInvoices); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allInvoices)
	}
}

func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		invoiceId := c.Param("invoice_id")

		var invoice models.Invoice

		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
		defer cancel()
		if err != nil {
			msg := "error occures while listing invoice item"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		var invoiceView InvoiceViewFormat

		allOrderItems, err1 := ItemsByOrder(invoice.Order_id)
		if err1 != nil {
			msg := "error getting order items"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		invoiceView.Order_id = invoice.Order_id
		invoiceView.Payment_due_date = invoice.Payment_due_date
		invoiceView.Payment_method = "null"
		if invoice.Payment_method != nil {
			invoiceView.Payment_method = *invoice.Payment_method
		}

		invoiceView.Invoice_id = invoice.Invoice_id
		invoiceView.Payment_status = *&invoice.Payment_status
		invoiceView.Payment_due = allOrderItems[0]["payment_due"]
		invoiceView.Table_number = allOrderItems[0]["table_number"] 
		invoiceView.Order_details = allOrderItems[0]["order_items"]

		c.JSON(http.StatusOK, invoiceView)
	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var invoice models.Invoice
		var order models.Order

		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(invoice)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := orderCollection.FindOne(ctx, bson.M{"order_id": order.Order_id}).Decode(&order)
		defer cancel()
		if err != nil {
			msg := "message: order was not found"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		status := "PENDING"
		if invoice.Payment_status == nil {
			invoice.Payment_status = &status
		}
		invoice.Payment_due_date, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
		invoice.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.ID = primitive.NewObjectID()
		invoice.Invoice_id = invoice.ID.Hex()

		resulst, insertErr := invoiceCollection.InsertOne(ctx, invoice)

		if insertErr != nil {
			msg := "invoice item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated, resulst)

	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var invoice models.Invoice
		invoiceId := c.Param("invoice_id")

		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"invoice_id": invoiceId}

		var updateObj primitive.D

		if invoice.Payment_method != nil {
			updateObj = append(updateObj, bson.E{Key: "payment_method", Value: invoice.Payment_method})
		}

		if invoice.Payment_status != nil {
			updateObj = append(updateObj, bson.E{Key: "payment_status", Value: invoice.Payment_status})
		}

		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: invoice.Updated_at})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		status := "PENDING"
		if invoice.Payment_status == nil {
			invoice.Payment_status = &status
		}

		result, err := invoiceCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := "invoice item update failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		defer cancel()
		c.JSON(http.StatusAccepted, result)

	}
}
