package routes

import (
	"github.com/Hand-TBN1/hand-backend/controller"
	"github.com/Hand-TBN1/hand-backend/middleware"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterCheckInRoutes(router *gin.Engine, db *gorm.DB) {

	checkInService := &services.CheckInService{DB: db}
	checkInController := &controller.CheckInController{CheckInService: checkInService}

	api := router.Group("/api")
	{
		api.POST("/checkins/create", middleware.RoleMiddleware("patient"), checkInController.CreateCheckIn)
		api.GET("/checkins/:id", middleware.RoleMiddleware("patient"), checkInController.GetCheckIn)
		api.GET("/checkins", middleware.RoleMiddleware("patient"), checkInController.GetAllCheckIns)
		api.PUT("/checkins/:id", middleware.RoleMiddleware("patient"), checkInController.UpdateCheckIn)
	}
}
