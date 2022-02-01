package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {
	godotenv.Load()
	mongoDbUri := os.Getenv("MONGO_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDbUri))
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to mongodb")
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection = client.Database("resturant").Collection(collectionName)

	return collection
}
