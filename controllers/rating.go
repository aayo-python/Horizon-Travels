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

// RatingRequest represents the structure for a rating request
type RatingRequest struct {
	UserID  string  `json:"user_id" binding:"required"`
	HotelID string  `json:"hotel_id" binding:"required"`
	Score   float64 `json:"score" binding:"required,min=1,max=5"`
	Comment string  `json:"comment"`
}

// CreateRating handles the creation of a new rating
func CreateRating(c *gin.Context) {
	var req RatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert string IDs to ObjectIDs
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	hotelID, err := primitive.ObjectIDFromHex(req.HotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID format"})
		return
	}

	now := time.Now()
	rating := model.Rating{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		HotelID:   hotelID,
		Score:     req.Score,
		Comment:   req.Comment,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save the rating to the database
	collection := util.MongoClient.Database(util.DbName).Collection(rating.CollectionName())
	result, err := collection.InsertOne(context.TODO(), rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rating"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Rating created successfully",
		"data":    rating,
		"id":      result.InsertedID,
	})
}

// GetRatingByID retrieves a rating by its ID
func GetRatingByID(c *gin.Context) {
	id := c.Param("id")

	ratingID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var rating model.Rating
	collection := util.MongoClient.Database(util.DbName).Collection(model.Rating{}.CollectionName())
	err = collection.FindOne(context.TODO(), bson.M{"_id": ratingID}).Decode(&rating)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rating not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rating"})
		return
	}

	c.JSON(http.StatusOK, rating)
}

// GetRatingsByHotel retrieves all ratings for a specific hotel
func GetRatingsByHotel(c *gin.Context) {
	hotelID := c.Param("hotelId")

	objectID, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID format"})
		return
	}

	var ratings []model.Rating
	collection := util.MongoClient.Database(util.DbName).Collection(model.Rating{}.CollectionName())

	cursor, err := collection.Find(context.TODO(), bson.M{"hotel_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ratings"})
		return
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &ratings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode ratings"})
		return
	}

	c.JSON(http.StatusOK, ratings)
}

// UpdateRating updates an existing rating
func UpdateRating(c *gin.Context) {
	id := c.Param("id")

	ratingID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req RatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"score":      req.Score,
			"comment":    req.Comment,
			"updated_at": time.Now(),
		},
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.Rating{}.CollectionName())
	result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": ratingID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update rating"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rating not found"})
		return
	}

	// Get the updated rating
	var updatedRating model.Rating
	err = collection.FindOne(context.TODO(), bson.M{"_id": ratingID}).Decode(&updatedRating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated rating"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rating updated successfully",
		"data":    updatedRating,
	})
}

// DeleteRating deletes a rating
func DeleteRating(c *gin.Context) {
	id := c.Param("id")

	ratingID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.Rating{}.CollectionName())
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": ratingID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete rating"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rating not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rating deleted successfully"})
}

// GetAverageRatingForHotel calculates the average rating for a hotel
func GetAverageRatingForHotel(c *gin.Context) {
	hotelID := c.Param("hotelId")

	objectID, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID format"})
		return
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.Rating{}.CollectionName())

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate average rating"})
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
