package util

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Define MongoDB connection string - better to use environment variables
var uri = getEnv("MONGODB_URI", "mongodb://localhost:27017")
var DbName = getEnv("MONGODB_DB_NAME", "travel")

// Create global variables to hold MongoDB connection
var (
	MongoClient *mongo.Client
	Db          *mongo.Collection
	ctx         context.Context
	cancel      context.CancelFunc
)

// Initialize MongoDB connection
func init() {
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)

	if err := ConnectToMongoDB(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	Db = MongoClient.Database(DbName).Collection("hotel")

	// Add a clean shutdown hook
	go func() {
		c := make(chan os.Signal, 1)
		<-c
		DisconnectMongoDB()
		os.Exit(0)
	}()
}

// ConnectToMongoDB establishes connection to MongoDB
func ConnectToMongoDB() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPI).
		SetTimeout(10 * time.Second)

	var err error
	MongoClient, err = mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}

	// Verify connection with ping
	if err = MongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	log.Println("Successfully connected to MongoDB")
	return nil
}

// DisconnectMongoDB closes the MongoDB connection
func DisconnectMongoDB() {
	if MongoClient != nil {
		if err := MongoClient.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}
	cancel()
	log.Println("MongoDB connection closed")
}

// GetCollection returns a specific collection from the database
func GetCollection(collectionName string) *mongo.Collection {
	return MongoClient.Database(DbName).Collection(collectionName)
}

// Helper function to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// package util

// import (
// 	"context"
// 	"log"
// 	"time"

// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// // Define your MongoDB connection string
// const uri = "mongodb://localhost:27017"

// // Create a global variable to hold our MongoDB connection
// var MongoClient *mongo.Client

// var Db *mongo.Collection
// var ctx context.Context

// // This function runs before we call our main function and connects to our MongoDB database. If it cannot connect, the application stops.
// func init() {
// 	if err := connect_to_mongodb(); err != nil {
// 		log.Fatal("Could not connect to MongoDB")
// 	}
// 	Db = MongoClient.Database("travel").Collection("hotel")
// 	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
// }

// // Our implementation code to connect to MongoDB at startup
// func connect_to_mongodb() error {
// 	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
// 	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

// 	client, err := mongo.Connect(context.TODO(), opts)
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = client.Ping(context.TODO(), nil)
// 	MongoClient = client
// 	return err
// }
