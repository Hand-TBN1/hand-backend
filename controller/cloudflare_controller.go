package controller

import (
	"net/http"

	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
)

type CloudflareController struct {
	CloudflareService *services.CloudflareService
}

// UploadCloudflare handles the image upload to Cloudflare R2
func (ctrl *CloudflareController) UploadCloudflare(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
        return
    }

    imageURL, apiErr := ctrl.CloudflareService.UploadCloudflare(c.Request.Context(), file)
    if apiErr != nil {
        c.JSON(apiErr.HttpStatus, gin.H{"error": apiErr.Message})
        return
    }

    c.JSON(http.StatusOK, gin.H{"image_url": imageURL})
}
