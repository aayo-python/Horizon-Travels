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

// CuisineRequest represents the structure for a cuisine request
// CreatedAt   time.Time          `json:"created_at" binding:"required"`
// UpdatedAt   time.Time          `json:"updated_at" binding:"required"`
type CuisineRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Origin      string `json:"origin" binding:"required"`
	Category    string `json:"category" binding:"required"`
}

// CreateCuisine handles the creation of a new cuisine
func CreateCuisine(c *gin.Context) {
	var req CuisineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	cuisine := model.Cuisine{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Description: req.Description,
		Origin:      req.Origin,
		Category:    req.Category,
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
		"message": "Cuisine created successfully",
		"data":    cuisine,
		"id":      result.InsertedID,
	})
}

// GetCuisineByID retrieves a cuisine by its ID
func GetCuisineByID(c *gin.Context) {
	id := c.Param("id")

	cuisineID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var cuisine model.Cuisine
	collection := util.MongoClient.Database(util.DbName).Collection(model.Cuisine{}.CollectionName())
	err = collection.FindOne(context.TODO(), bson.M{"_id": cuisineID}).Decode(&cuisine)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cuisine not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cuisine"})
		return
	}

	c.JSON(http.StatusOK, cuisine)
}

// GetCuisinesByHotel retrieves all cuisines for a specific hotel
func GetCuisinesByHotel(c *gin.Context) {
	hotelID := c.Param("hotelId")

	objectID, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID format"})
		return
	}

	var cuisines []model.Cuisine
	collection := util.MongoClient.Database(util.DbName).Collection(model.Cuisine{}.CollectionName())

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

// UpdateCuisine updates an existing cuisine
func UpdateCuisine(c *gin.Context) {
	id := c.Param("id")

	// cuisineID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req CuisineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if cuisine exists
	var existingHotel model.Cuisine
	err = util.Db.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&existingCuisine)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cuisine not found"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"score":      req,
			"comment":    req.Comment,
			"updated_at": time.Now(),
		},
	}

	result, err := util.Db.UpdateByID(context.TODO(), id, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No changes applied"})
		return
	}

	// Get updated hotel
	var updatedCuisine model.CupdatedCuisine
	err = util.Db.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&updatedCuisine)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, updatedCuisine)

	// collection := util.MongoClient.Database(util.DbName).Collection(model.Cuisine{}.CollectionName())
	// result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": cuisineID}, update)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cuisine"})
	// 	return
	// }

	// if result.MatchedCount == 0 {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Cuisine not found"})
	// 	return
	// }

	// // Get the updated cuisine
	// var updatedCuisine model.Cuisine
	// err = collection.FindOne(context.TODO(), bson.M{"_id": cuisineID}).Decode(&updatedCuisine)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated cuisine"})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "Cuisine updated successfully",
	// 	"data":    updatedCuisine,
	// })
}

// DeleteCuisine deletes a cuisine
func DeleteCuisine(c *gin.Context) {
	id := c.Param("id")

	cuisineID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.Cuisine{}.CollectionName())
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": cuisineID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cuisine"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cuisine not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cuisine deleted successfully"})
}

// GetAverageCuisineForHotel calculates the average cuisine for a hotel
func GetAverageCuisineForHotel(c *gin.Context) {
	hotelID := c.Param("hotelId")

	objectID, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID format"})
		return
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.Cuisine{}.CollectionName())

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
