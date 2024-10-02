package routes

import (
	"github.com/Hand-TBN1/hand-backend/controller"
	"github.com/Hand-TBN1/hand-backend/middleware"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
)

func RegisterCloudflareRoutes(router *gin.Engine) {
	cloudflareService := &services.CloudflareService{}
	cloudflareController := &controller.CloudflareController{CloudflareService: cloudflareService}

	api := router.Group("/api")
	{
		api.POST("/upload-image", middleware.RoleMiddleware(), cloudflareController.UploadCloudflare) 
	}
}
