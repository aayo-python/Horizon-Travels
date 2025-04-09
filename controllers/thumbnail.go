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

type ThumbnailRequest struct {
	UserID     primitive.ObjectID `json:"user_id"  binding:"required"`
	EntityID   primitive.ObjectID `json:"entity_id" binding:"required"`
	EntityType string             `json:"entity_type" binding:"required"` // "hotel", "restaurant", etc.
	URL        string             `json:"url" binding:"required"`
	HEIGHT     string             `json:"height" binding:"required"`
	Rating     float64            `json:"rating" binding:"required"`
	Helpful    int                `json:"helpful" binding:"required"`
}

// CreateThumbnail handles the creation of a new thumbnail
func CreateThumbnail(c *gin.Context) {
	var req ThumbnailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	thumbnail := model.Thumbnail{
		ID:         primitive.NewObjectID(),
		EntityID:   req.EntityID,
		EntityType: req.EntityType,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Save the thumbnail to the database
	collection := util.MongoClient.Database(util.DbName).Collection(thumbnail.CollectionName())
	result, err := collection.InsertOne(context.TODO(), thumbnail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create thumbnail"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Thumbnail created successfully",
		"data":    thumbnail,
		"id":      result.InsertedID,
	})
}

// GetThumbnailByID retrieves a thumbnail by its ID
func GetThumbnailByID(c *gin.Context) {
	id := c.Param("id")

	thumbnailID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var thumbnail model.Thumbnail
	collection := util.MongoClient.Database(util.DbName).Collection(model.Thumbnail{}.CollectionName())
	err = collection.FindOne(context.TODO(), bson.M{"_id": thumbnailID}).Decode(&thumbnail)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Thumbnail not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve thumbnail"})
		return
	}

	c.JSON(http.StatusOK, thumbnail)
}

// GetThumbnailsByHotel retrieves all thumbnails for a specific hotel
func GetThumbnailsByHotel(c *gin.Context) {
	hotelID := c.Param("hotelId")

	objectID, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID format"})
		return
	}

	var thumbnails []model.Thumbnail
	collection := util.MongoClient.Database(util.DbName).Collection(model.Thumbnail{}.CollectionName())

	cursor, err := collection.Find(context.TODO(), bson.M{"hotel_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve thumbnails"})
		return
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &thumbnails); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode thumbnails"})
		return
	}

	c.JSON(http.StatusOK, thumbnails)
}

// UpdateThumbnail updates an existing thumbnail
func UpdateThumbnail(c *gin.Context) {
	id := c.Param("id")

	thumbnailID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req ThumbnailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if thumbnail exists
	var existingThumbnail model.Thumbnail
	err = util.Db.FindOne(context.TODO(), bson.D{{"_id", thumbnailID}}).Decode(&existingThumbnail)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Thumbnail not found"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"URL":        req.URL,
			"EntityID":   req.EntityID,
			"EntityType": req.EntityType,
			"HEIGHT":     req.HEIGHT,
			"updated_at": time.Now(),
		},
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.Thumbnail{}.CollectionName())
	result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": thumbnailID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update thumbnail"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Thumbnail not found"})
		return
	}

	// Get the updated thumbnail
	var updatedThumbnail model.Thumbnail
	err = collection.FindOne(context.TODO(), bson.M{"_id": thumbnailID}).Decode(&updatedThumbnail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated thumbnail"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Thumbnail updated successfully",
		"data":    updatedThumbnail,
	})
}

// DeleteThumbnail deletes a thumbnail
func DeleteThumbnail(c *gin.Context) {
	id := c.Param("id")

	thumbnailID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.Thumbnail{}.CollectionName())
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": thumbnailID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete thumbnail"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Thumbnail not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Thumbnail deleted successfully"})
}

// GetAverageThumbnailForHotel calculates the average thumbnail for a hotel
func GetAverageThumbnailForHotel(c *gin.Context) {
	hotelID := c.Param("hotelId")

	objectID, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID format"})
		return
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.Thumbnail{}.CollectionName())

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate average thumbnail"})
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
