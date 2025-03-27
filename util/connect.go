package util

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define your MongoDB connection string
const uri = "mongodb://localhost:27017"

// Create a global variable to hold our MongoDB connection
var MongoClient *mongo.Client

var Db *mongo.Collection
var ctx context.Context

// This function runs before we call our main function and connects to our MongoDB database. If it cannot connect, the application stops.
func init() {
	if err := connect_to_mongodb(); err != nil {
		log.Fatal("Could not connect to MongoDB")
	}
	Db = MongoClient.Database("travel").Collection("hotel")
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
}

// Our implementation code to connect to MongoDB at startup
func connect_to_mongodb() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	MongoClient = client
	return err
}
