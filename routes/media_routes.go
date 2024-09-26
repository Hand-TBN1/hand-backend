package routes

import (
	"github.com/Hand-TBN1/hand-backend/controller"
	"github.com/Hand-TBN1/hand-backend/middleware"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterMediaRoutes(router *gin.Engine, db *gorm.DB) {
	mediaService := &services.MediaService{DB: db}
	mediaController := &controller.MediaController{MediaService: mediaService}

	api := router.Group("/api")
	{
		adminTherapistRoutes := api.Group("/media")
		adminTherapistRoutes.Use(middleware.RoleMiddleware("admin", "therapist"))
		{
			adminTherapistRoutes.POST("", mediaController.CreateMedia)
			adminTherapistRoutes.PUT("/:id", mediaController.UpdateMedia)
			adminTherapistRoutes.DELETE("/:id", mediaController.DeleteMedia)
		}
		allRolesRoutes := api.Group("/media")
		allRolesRoutes.Use(middleware.RoleMiddleware("patient", "therapist", "admin"))
		{
			allRolesRoutes.GET("", mediaController.GetAllMedia)
			allRolesRoutes.GET("/:id", mediaController.GetMedia)
		}
	}
}
