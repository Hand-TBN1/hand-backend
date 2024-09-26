package controller

import (
	"net/http"
	"time"

	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateMediaDTO uses the enum MediaType instead of string
type CreateMediaDTO struct {
	Type        models.MediaType `json:"type" binding:"required"`
	Title       string           `json:"title" binding:"required"`
	Content     string           `json:"content" binding:"required"`
	ThumbnailURL string          `json:"image_url"`
}

type UpdateMediaDTO struct {
	Type        models.MediaType `json:"type" binding:"required"`
	Title       string           `json:"title" binding:"required"`
	Content     string           `json:"content" binding:"required"`
	ThumbnailURL string          `json:"image_url"`
}

type MediaController struct {
	MediaService *services.MediaService
}

func (ctrl *MediaController) CreateMedia(c *gin.Context) {
	var dto CreateMediaDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	media := models.Media{
		ID:           uuid.New(),
		Type:         dto.Type,
		Title:        dto.Title,
		Content:      dto.Content,
		ThumbnailURL: dto.ThumbnailURL,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	apiErr := ctrl.MediaService.AddMedia(&media)
	if apiErr != nil {
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Media created successfully"})
}

func (ctrl *MediaController) GetAllMedia(c *gin.Context) {
	mediaList, apiErr := ctrl.MediaService.GetAllMedia()
	if apiErr != nil {
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	c.JSON(http.StatusOK, mediaList)
}

func (ctrl *MediaController) GetMedia(c *gin.Context) {
	id := c.Param("id")
	media, apiErr := ctrl.MediaService.GetMedia(id)
	if apiErr != nil {
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	c.JSON(http.StatusOK, media)
}

func (ctrl *MediaController) UpdateMedia(c *gin.Context) {
	id := c.Param("id")
	var dto UpdateMediaDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	media, apiErr := ctrl.MediaService.GetMedia(id)
	if apiErr != nil {
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	media.Type = dto.Type
	media.Title = dto.Title
	media.Content = dto.Content
	media.ThumbnailURL = dto.ThumbnailURL
	media.UpdatedAt = time.Now()

	apiErr = ctrl.MediaService.UpdateMedia(media)
	if apiErr != nil {
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Media updated successfully"})
}

func (ctrl *MediaController) DeleteMedia(c *gin.Context) {
	id := c.Param("id")
	apiErr := ctrl.MediaService.DeleteMedia(id)
	if apiErr != nil {
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	c.Status(http.StatusNoContent)
}
