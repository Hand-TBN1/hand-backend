package routes

import (
	"github.com/Hand-TBN1/hand-backend/controller"
	"github.com/Hand-TBN1/hand-backend/middleware"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(router *gin.Engine, db *gorm.DB) {
	userService := &services.UserService{DB: db}
	userController := &controller.UserController{UserService: userService}

	api := router.Group("/api")
	{
		api.GET("/profile", middleware.RoleMiddleware(), userController.GetProfile)   // Get user profile
		api.PUT("/edit-profile", middleware.RoleMiddleware(), userController.EditProfile)  // Edit user profile
	}
}
