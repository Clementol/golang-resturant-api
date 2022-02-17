package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	db "github.com/Clementol/restur-manag/database"
	helper "github.com/Clementol/restur-manag/helpers"
	"github.com/Clementol/restur-manag/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection = db.OpenCollection(db.Client, "user")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))

		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		// startIndex := (page - 1) * recordPerPage
		// startIndex, err = strconv.Atoi(c.Query("startIndex"))

		// matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		projectStage := []bson.M{
			{"$project": bson.M{
				"_id": 0,
				// "total_count": 1,
				// "users_items": bson.M{"$slice": []interface{}{ startIndex, recordPerPage}},

			}}}

		result, err := userCollection.Aggregate(ctx,
			projectStage,
		)
		if err != nil {
			msg := "error occurred while listing user items" + err.Error()
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		var allUsers []bson.M

		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allUsers)

	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		userId := c.MustGet("user_id").(string)

		var user bson.M

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			msg := "error occured while getting user" + err.Error()
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, user)
	}
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		file, _ := c.FormFile("file")

		if err := c.Bind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// validate the data based on user struct

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			msg := "error occurred while checking for email"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		password := helper.HashPassword(*user.Password)
		user.Password = &password

		phoneCount, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			msg := "error occured while checking for phone number"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		if emailCount > 0 || phoneCount > 0 {
			msg := "this email or phone number already exist"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_Name, *user.Last_Name, user.User_id)
		user.Token = &token
		user.Refresh_Token = &refreshToken

		if file != nil {
			imageText := helper.RenameFileName(*file)
			user.Avatar = &imageText
		}

		result, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := "user was not registered"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": &user.Email}).Decode(&foundUser)
		if err != nil {
			msg := "invalid user credentials"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		defer cancel()
		passwordInvalid, msg := helper.VerifyPassword(*user.Password, *foundUser.Password)
		if !passwordInvalid {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		defer cancel()

		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, foundUser.User_id)

		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		c.JSON(http.StatusOK, foundUser)

	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("user_id")

		file, _ := c.FormFile("file")

		var user models.User

		if err := c.Bind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updateObj := bson.M{}

		if file != nil {
			imageText := helper.RenameFileName(*file)
			updateObj["avatar"] = imageText
		}

		if user.Phone != nil {
			updateObj["phone"] = user.Phone
		}

		if user.First_Name != nil {
			updateObj["first_name"] = user.First_Name
		}

		if user.Last_Name != nil {

			updateObj["last_name"] = user.Last_Name
		}

		updateObj["updated_at"], _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		updatedUser := bson.M{}
		filter := bson.M{"user_id": userId}
		opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
		err := userCollection.FindOneAndUpdate(ctx, filter,
			bson.M{"$set": updateObj}, opt).Decode(&updatedUser)

		if err != nil {
			msg := "invalid credentials " + err.Error()
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusAccepted, updatedUser)

	}
}
