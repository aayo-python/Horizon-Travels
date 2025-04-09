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

type VacationRentalRequest struct {
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

// CreateVacationRental handles the creation of a new vacationRental
func CreateVacationRental(c *gin.Context) {
	var req VacationRentalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	vacationRental := model.VacationRental{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Save the vacationRental to the database
	collection := util.MongoClient.Database(util.DbName).Collection(vacationRental.CollectionName())
	result, err := collection.InsertOne(context.TODO(), vacationRental)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vacationRental"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "VacationRental created successfully",
		"data":    vacationRental,
		"id":      result.InsertedID,
	})
}

// GetVacationRentalByID retrieves a vacationRental by its ID
func GetVacationRentalByID(c *gin.Context) {
	id := c.Param("id")

	vacationRentalID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var vacationRental model.VacationRental
	collection := util.MongoClient.Database(util.DbName).Collection(model.VacationRental{}.CollectionName())
	err = collection.FindOne(context.TODO(), bson.M{"_id": vacationRentalID}).Decode(&vacationRental)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "VacationRental not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve vacationRental"})
		return
	}

	c.JSON(http.StatusOK, vacationRental)
}

// GetVacationRentalsByHotel retrieves all vacationRentals for a specific hotel
func GetVacationRentalsByHotel(c *gin.Context) {
	hotelID := c.Param("hotelId")

	objectID, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID format"})
		return
	}

	var vacationRentals []model.VacationRental
	collection := util.MongoClient.Database(util.DbName).Collection(model.VacationRental{}.CollectionName())

	cursor, err := collection.Find(context.TODO(), bson.M{"hotel_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve vacationRentals"})
		return
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &vacationRentals); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode vacationRentals"})
		return
	}

	c.JSON(http.StatusOK, vacationRentals)
}

// UpdateVacationRental updates an existing vacationRental
func UpdateVacationRental(c *gin.Context) {
	id := c.Param("id")

	vacationRentalID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req VacationRentalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if vacationRental exists
	var existingVacationRental model.VacationRental
	err = util.Db.FindOne(context.TODO(), bson.D{{"_id", vacationRentalID}}).Decode(&existingVacationRental)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "VacationRental not found"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"origin":     req.Origin,
			"name":       req.Name,
			"updated_at": time.Now(),
		},
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.VacationRental{}.CollectionName())
	result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": vacationRentalID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vacationRental"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "VacationRental not found"})
		return
	}

	// Get the updated vacationRental
	var updatedVacationRental model.VacationRental
	err = collection.FindOne(context.TODO(), bson.M{"_id": vacationRentalID}).Decode(&updatedVacationRental)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated vacationRental"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "VacationRental updated successfully",
		"data":    updatedVacationRental,
	})
}

// DeleteVacationRental deletes a vacationRental
func DeleteVacationRental(c *gin.Context) {
	id := c.Param("id")

	vacationRentalID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.VacationRental{}.CollectionName())
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": vacationRentalID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vacationRental"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "VacationRental not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "VacationRental deleted successfully"})
}

// GetAverageVacationRentalForHotel calculates the average vacationRental for a hotel
func GetAverageVacationRentalForHotel(c *gin.Context) {
	hotelID := c.Param("hotelId")

	objectID, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID format"})
		return
	}

	collection := util.MongoClient.Database(util.DbName).Collection(model.VacationRental{}.CollectionName())

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate average vacationRental"})
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
