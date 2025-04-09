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

type ReviewRequest struct {
	UserID     primitive.ObjectID `json:"user_id"  binding:"required"`
	EntityID   primitive.ObjectID `json:"entity_id" binding:"required"`
	EntityType string             `json:"entity_type" binding:"required"` // "hotel", "restaurant", etc.
	Title      string             `json:"title" binding:"required"`
	Content    string             `json:"content" binding:"required"`
	Rating     float64            `json:"rating" binding:"required"`
	Helpful    int                `json:"helpful" binding:"required"`
}

// CreateReview handles the creation of a new review
func CreateReview(c *gin.Context) {
	var req ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	review := model.Review{
		ID:         primitive.NewObjectID(),
		UserID:     req.UserID,
		EntityID:   req.EntityID,
		EntityType: req.EntityType,
		Title:      req.Title,
		Content:    req.Content,
		Rating:     req.Rating,
		Helpful:    req.Helpful,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Save the review to the database
	collection := util.MongoClient.Database(util.DbName).Collection(review.CollectionName())
	result, err := collection.InsertOne(context.TODO(), review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Review created successfully",
		"data":    review,
		"id":      result.InsertedID,
	})
}

// GetReviewByID retrieves a review by its ID
func GetReviewByID(c *gin.Context) {
	id := c.Param("id")

	reviewID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var review model.Review
	collection := util.MongoClient.Database(util.DbName).Collection(model.Review{}.CollectionName())
	err = collection.FindOne(context.TODO(), bson.M{"_id": reviewID}).Decode(&review)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve review"})
		return
	}

	c.JSON(http.StatusOK, review)
}

// GetReviewsByHotel retrieves all reviews for a specific hotel
func GetReviewsByHotel(c *gin.Context) {
	hotelID := c.Param("hotelId")

	objectID, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID format"})
		return
	}

	var reviews []model.Review
	collection := util.MongoClient.Database(util.DbName).Collection(model.Review{}.CollectionName())

	cursor, err := collection.Find(context.TODO(), bson.M{"hotel_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve reviews"})
		return
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &reviews); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// UpdateReview updates an existing review
func UpdateReview(c *gin.Context) {
	id := c.Param("id")

	reviewID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if review exists
	var existingReview model.Review
	err = util.Db.FindOne(context.TODO(), bson.D{{"_id", reviewID}}).Decode(&existingReview)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"title":      req.Title,
			"EntityID":   req.EntityID,
			"EntityType": req.EntityType,
			"Title":      req.Title,
			"updated_at": time.Now(),
		},
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.Review{}.CollectionName())
	result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": reviewID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	// Get the updated review
	var updatedReview model.Review
	err = collection.FindOne(context.TODO(), bson.M{"_id": reviewID}).Decode(&updatedReview)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Review updated successfully",
		"data":    updatedReview,
	})
}

// DeleteReview deletes a review
func DeleteReview(c *gin.Context) {
	id := c.Param("id")

	reviewID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.Review{}.CollectionName())
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": reviewID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}

// GetAverageReviewForHotel calculates the average review for a hotel
func GetAverageReviewForHotel(c *gin.Context) {
	hotelID := c.Param("hotelId")

	objectID, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID format"})
		return
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.Review{}.CollectionName())

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate average review"})
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
