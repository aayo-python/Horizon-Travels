package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/janto-pee/Horizon-Travels.git/model"
	"github.com/janto-pee/Horizon-Travels.git/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// List Hotel
type listHotelsRequest struct {
	PageID   int64
	PageSize int64
}

func ListHotels(c *gin.Context) {
	var req listHotelsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	skip := (req.PageID - 1) * req.PageSize
	fmt.Println(req.PageID, req.PageSize, skip)
	opts := options.Find().SetLimit(int64(req.PageSize)).SetSkip(skip)
	cursor, err := util.Db.Find(context.TODO(), bson.D{{}}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var hotels []bson.M
	if err = cursor.All(context.TODO(), &hotels); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hotels)
}

// Get Hotel Details
func GetHotelByID(c *gin.Context) {

	idStr := c.Param("id")
	fmt.Println(idStr)

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var hotel bson.M
	err = util.Db.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&hotel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hotel)
}

// Search Hotel By Location &
// Search Hotel
type SearchHotelsRequest struct {
	PageID   int
	PageSize int
	Keyword  string
}

func SearchHotels(c *gin.Context) {
	PageID, err := strconv.Atoi(c.Query("PageID"))
	PageSize, err := strconv.Atoi(c.Query("PageSize"))
	Keyword := c.Query("Keyword")

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	filter := bson.D{{"$text", bson.D{{"$search", Keyword}}}}
	skip := (PageID - 1) * PageSize

	opts := options.Find().SetLimit(int64(PageSize)).SetSkip(int64(skip))
	cursor, err := util.Db.Find(context.TODO(), filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var hotels []bson.M
	if err = cursor.All(context.TODO(), &hotels); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hotels)
}

// Get Hotel Filters
type FilterHotelRequest struct {
	Title         string
	Content       string
	PrimaryInfo   string
	SecondaryInfo string
	AccentedLabel string
	Provider      string
	PriceDetails  string
	PriceSummary  string
}

func FilterHotels(c *gin.Context) {
	var req listHotelsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	filter := bson.D{{"title", bson.D{{"$regex", "^E"}}}, {"provider", bson.D{{"$regex", "^E"}}}}
	skip := (req.PageID - 1) * req.PageSize
	fmt.Println(req.PageID, req.PageSize, skip)
	opts := options.Find().SetLimit(int64(req.PageSize)).SetSkip(skip)
	cursor, err := util.Db.Find(context.TODO(), filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var hotels []bson.M
	if err = cursor.All(context.TODO(), &hotels); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hotels)
}

func AggregateHotels(c *gin.Context) {
	var pipeline interface{}
	if err := c.ShouldBindJSON(&pipeline); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cursor, err := util.Db.Aggregate(context.TODO(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []bson.M
	if err = cursor.All(context.TODO(), &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

/*
*
*	MUTATIONS
*
 */

type createHotelRequest struct {
	Title         string
	Content       string
	PrimaryInfo   string
	SecondaryInfo string
	AccentedLabel string
	Provider      string
	PriceDetails  string
	PriceSummary  string
}

func CreateHotels(c *gin.Context) {
	var req createHotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := model.Hotel{
		Title:         req.Title,
		Content:       req.Content,
		PrimaryInfo:   req.PrimaryInfo,
		SecondaryInfo: req.SecondaryInfo,
		AccentedLabel: req.AccentedLabel,
		Provider:      req.Provider,
		PriceDetails:  req.PriceDetails,
		PriceSummary:  req.PriceSummary,
	}

	hotel, err := util.Db.InsertOne(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, hotel)

}

type updateHotelParam struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func updateHotel(c *gin.Context) {
	var id updateHotelParam
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req createHotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := model.Hotel{
		Title:         req.Title,
		Content:       req.Content,
		PrimaryInfo:   req.PrimaryInfo,
		SecondaryInfo: req.SecondaryInfo,
		AccentedLabel: req.AccentedLabel,
		Provider:      req.Provider,
		PriceDetails:  req.PriceDetails,
		PriceSummary:  req.PriceSummary,
	}
	if err := util.Db.FindOne(c, bson.M{"id": id}).Decode(&arg); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err := c.BindJSON(&arg); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	hotel, err := util.Db.UpdateByID(c, bson.M{"id": id}, bson.M{"$set": arg})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, hotel)
}

func deleteHotel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := util.Db.DeleteOne(c, bson.M{"id": id})
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}
