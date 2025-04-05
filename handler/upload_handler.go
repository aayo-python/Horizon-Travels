package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/janto-pee/fintech-platform/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UploadHandler handles image uploads
type UploadHandler struct {
	s3Uploader *utils.S3Uploader
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(bucketName, region string) (*UploadHandler, error) {
	uploader, err := utils.NewS3Uploader(bucketName, region)
	if err != nil {
		return nil, err
	}

	return &UploadHandler{
		s3Uploader: uploader,
	}, nil
}

// UploadImage handles image upload requests
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// Get entity information
	entityID := c.PostForm("entity_id")
	entityType := c.PostForm("entity_type")
	title := c.PostForm("title")
	description := c.PostForm("description")

	// Validate entity ID
	objID, err := primitive.ObjectIDFromHex(entityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	// Get file from request
	file, fileHeader, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Upload file to S3
	fileURL, err := h.s3Uploader.UploadFile(file, fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	// Create photo record
	photo := model.Photo{
		ID:          primitive.NewObjectID(),
		URL:         fileURL,
		Title:       title,
		Description: description,
		EntityID:    objID,
		EntityType:  entityType,
		IsPrimary:   false,
	}

	// Here you would save the photo to your database
	// For example: photoRepository.Create(photo)

	c.JSON(http.StatusOK, gin.H{
		"message": "Image uploaded successfully",
		"photo":   photo,
	})
}
