package routes

import (
	"github.com/Hand-TBN1/hand-backend/controller"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAuthRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize services
	authService := &services.AuthService{DB: db}
	authController := &controller.AuthController{AuthService: authService}

	api := router.Group("/api")
	{
		// Authentication routes
		api.POST("/register", authController.Register)
		api.POST("/login", authController.Login)
	}
}
