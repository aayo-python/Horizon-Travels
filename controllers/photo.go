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

type PhotoRequest struct {
	URL         string             `json:"url" binding:"required"`
	Title       string             `json:"title"  binding:"required"`
	Description string             `json:"description" binding:"required"`
	EntityID    primitive.ObjectID `json:"entity_id"  `
	EntityType  string             `json:"entity_type"  `
	IsPrimary   bool               `json:"is_primary" `
}

// CreatePhoto handles the creation of a new cuisine
func CreatePhoto(c *gin.Context) {
	var req PhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	cuisine := model.Photo{
		ID:          primitive.NewObjectID(),
		URL:         req.URL,
		Title:       req.Title,
		Description: req.Description,
		EntityID:    req.EntityID,
		EntityType:  req.EntityType,
		IsPrimary:   req.IsPrimary,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Save the cuisine to the database
	collection := util.MongoClient.Database(util.DbName).Collection(cuisine.CollectionName())
	result, err := collection.InsertOne(context.TODO(), cuisine)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cuisine"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Photo created successfully",
		"data":    cuisine,
		"id":      result.InsertedID,
	})
}

// GetPhotoByID retrieves a cuisine by its ID
func GetPhotoByID(c *gin.Context) {
	id := c.Param("id")

	cuisineID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var cuisine model.Photo
	collection := util.MongoClient.Database(util.DbName).Collection(model.Photo{}.CollectionName())
	err = collection.FindOne(context.TODO(), bson.M{"_id": cuisineID}).Decode(&cuisine)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cuisine"})
		return
	}

	c.JSON(http.StatusOK, cuisine)
}

// GetPhotosByHotel retrieves all cuisines for a specific hotel
func GetPhotosByHotel(c *gin.Context) {
	hotelID := c.Param("hotelId")

	objectID, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID format"})
		return
	}

	var cuisines []model.Photo
	collection := util.MongoClient.Database(util.DbName).Collection(model.Photo{}.CollectionName())

	cursor, err := collection.Find(context.TODO(), bson.M{"hotel_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cuisines"})
		return
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &cuisines); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode cuisines"})
		return
	}

	c.JSON(http.StatusOK, cuisines)
}

// UpdatePhoto updates an existing cuisine
func UpdatePhoto(c *gin.Context) {
	id := c.Param("id")

	cuisineID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req PhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if cuisine exists
	var existingPhoto model.Photo
	err = util.Db.FindOne(context.TODO(), bson.D{{"_id", cuisineID}}).Decode(&existingPhoto)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"description": req.Description,
			"isPrimry":    req.IsPrimary,
			"updated_at":  time.Now(),
		},
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.Photo{}.CollectionName())
	result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": cuisineID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cuisine"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	// Get the updated cuisine
	var updatedPhoto model.Photo
	err = collection.FindOne(context.TODO(), bson.M{"_id": cuisineID}).Decode(&updatedPhoto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated cuisine"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Photo updated successfully",
		"data":    updatedPhoto,
	})
}

// DeletePhoto deletes a cuisine
func DeletePhoto(c *gin.Context) {
	id := c.Param("id")

	cuisineID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.Photo{}.CollectionName())
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": cuisineID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cuisine"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}

// GetAveragePhotoForHotel calculates the average cuisine for a hotel
func GetAveragePhotoForHotel(c *gin.Context) {
	hotelID := c.Param("hotelId")

	objectID, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID format"})
		return
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.Photo{}.CollectionName())

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate average cuisine"})
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
