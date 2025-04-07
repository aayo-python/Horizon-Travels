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

// CreatePhoto handles the creation of a new photo
func CreatePhoto(c *gin.Context) {
	var req PhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	photo := model.Photo{
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

	// Save the photo to the database
	collection := util.MongoClient.Database(util.DbName).Collection(photo.CollectionName())
	result, err := collection.InsertOne(context.TODO(), photo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create photo"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Photo created successfully",
		"data":    photo,
		"id":      result.InsertedID,
	})
}

// GetPhotoByID retrieves a photo by its ID
func GetPhotoByID(c *gin.Context) {
	id := c.Param("id")

	photoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var photo model.Photo
	collection := util.MongoClient.Database(util.DbName).Collection(model.Photo{}.CollectionName())
	err = collection.FindOne(context.TODO(), bson.M{"_id": photoID}).Decode(&photo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve photo"})
		return
	}

	c.JSON(http.StatusOK, photo)
}

// GetPhotosByHotel retrieves all photos for a specific hotel
func GetPhotosByHotel(c *gin.Context) {
	hotelID := c.Param("hotelId")

	objectID, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID format"})
		return
	}

	var photos []model.Photo
	collection := util.MongoClient.Database(util.DbName).Collection(model.Photo{}.CollectionName())

	cursor, err := collection.Find(context.TODO(), bson.M{"hotel_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve photos"})
		return
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &photos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode photos"})
		return
	}

	c.JSON(http.StatusOK, photos)
}

// UpdatePhoto updates an existing photo
func UpdatePhoto(c *gin.Context) {
	id := c.Param("id")

	photoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req PhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if photo exists
	var existingPhoto model.Photo
	err = util.Db.FindOne(context.TODO(), bson.D{{"_id", photoID}}).Decode(&existingPhoto)
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
	result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": photoID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update photo"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	// Get the updated photo
	var updatedPhoto model.Photo
	err = collection.FindOne(context.TODO(), bson.M{"_id": photoID}).Decode(&updatedPhoto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated photo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Photo updated successfully",
		"data":    updatedPhoto,
	})
}

// DeletePhoto deletes a photo
func DeletePhoto(c *gin.Context) {
	id := c.Param("id")

	photoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.Photo{}.CollectionName())
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": photoID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photo"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}

// GetAveragePhotoForHotel calculates the average photo for a hotel
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate average photo"})
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
