package api

import (
	"context"
	"log"
	"time"

	// Add the MongoDB driver packages

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Your MongoDB Atlas Connection String
const uri = "mongodb://localhost:27017/travel"

// A global variable that will hold a reference to the MongoDB client
var mongoClient *mongo.Client

var db *mongo.Collection
var ctx context.Context

// The init function will run before our main function to establish a connection to MongoDB. If it cannot connect it will fail and the program will exit.
func init() {
	if err := connect_to_mongodb(); err != nil {
		log.Fatal("Could not connect to MongoDB")
	}
	db = mongoClient.Database("travel").Collection("hotel")
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

}

// Our implementation logic for connecting to MongoDB
func connect_to_mongodb() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	mongoClient = client
	return err
}
