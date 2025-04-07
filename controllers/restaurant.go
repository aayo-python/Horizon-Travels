package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/janto-pee/Horizon-Travels.git/model"
	"github.com/janto-pee/Horizon-Travels.git/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RestaurantRequest struct {
	Name          string            `json:"name" binding:"required"`
	Description   string            `json:"description" binding:"required"`
	Origin        string            `json:"origin" binding:"required"`
	Category      string            `json:"category" binding:"required"`
	Address       string            `json:"address" `
	Location      string            `json:"location" `
	PriceRange    string            `json:"price_range" `
	Rating        float64           `json:"rating" json:"rating"`
	ContactNumber string            `json:"contact_number"  `
	Email         string            `json:"email" `
	Website       string            `json:"website" json:"website"`
	OpeningHours  map[string]string `json:"opening_hours" json:"opening_hours"`
}

// CreateRestaurant handles the creation of a new restaurant
func CreateRestaurant(c *gin.Context) {
	var req RestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	restaurant := model.Restaurant{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Save the restaurant to the database
	collection := util.MongoClient.Database(util.DbName).Collection(restaurant.CollectionName())
	result, err := collection.InsertOne(context.TODO(), restaurant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create restaurant"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Restaurant created successfully",
		"data":    restaurant,
		"id":      result.InsertedID,
	})
}

// GetRestaurantByID retrieves a restaurant by its ID
func GetRestaurantByID(c *gin.Context) {
	id := c.Param("id")

	restaurantID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var restaurant model.Restaurant
	collection := util.MongoClient.Database(util.DbName).Collection(model.Restaurant{}.CollectionName())
	err = collection.FindOne(context.TODO(), bson.M{"_id": restaurantID}).Decode(&restaurant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve restaurant"})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

// GetRestaurantsByHotel retrieves all restaurants for a specific hotel
func GetRestaurantsByHotel(c *gin.Context) {
	hotelID := c.Param("hotelId")

	objectID, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID format"})
		return
	}

	var restaurants []model.Restaurant
	collection := util.MongoClient.Database(util.DbName).Collection(model.Restaurant{}.CollectionName())

	cursor, err := collection.Find(context.TODO(), bson.M{"hotel_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve restaurants"})
		return
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &restaurants); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode restaurants"})
		return
	}

	c.JSON(http.StatusOK, restaurants)
}

// UpdateRestaurant updates an existing restaurant
func UpdateRestaurant(c *gin.Context) {
	id := c.Param("id")

	restaurantID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req RestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if restaurant exists
	var existingRestaurant model.Restaurant
	err = util.Db.FindOne(context.TODO(), bson.D{{"_id", restaurantID}}).Decode(&existingRestaurant)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"origin":     req.Origin,
			"name":       req.Name,
			"updated_at": time.Now(),
		},
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.Restaurant{}.CollectionName())
	result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": restaurantID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update restaurant"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}

	// Get the updated restaurant
	var updatedRestaurant model.Restaurant
	err = collection.FindOne(context.TODO(), bson.M{"_id": restaurantID}).Decode(&updatedRestaurant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated restaurant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Restaurant updated successfully",
		"data":    updatedRestaurant,
	})
}

// DeleteRestaurant deletes a restaurant
func DeleteRestaurant(c *gin.Context) {
	id := c.Param("id")

	restaurantID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.Restaurant{}.CollectionName())
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": restaurantID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete restaurant"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Restaurant deleted successfully"})
}

// GetAverageRestaurantForHotel calculates the average restaurant for a hotel
func GetAverageRestaurantForHotel(c *gin.Context) {
	hotelID := c.Param("hotelId")

	objectID, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID format"})
		return
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.Restaurant{}.CollectionName())

	// MongoDB aggregation pipeline to calculate average
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"hotel_id": objectID}}},
		{{"$group", bson.M{
			"_id":          "$hotel_id",
			"averageScore": bson.M{"$avg": "$score"},
			"count":        bson.M{"$sum": 1},
		}}},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate average restaurant"})
		return
	}
	defer cursor.Close(context.TODO())

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode results"})
		return
	}

	if len(results) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"hotel_id":      hotelID,
			"average_score": 0,
			"count":         0,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"hotel_id":      hotelID,
		"average_score": results[0]["averageScore"],
		"count":         results[0]["count"],
	})
}
