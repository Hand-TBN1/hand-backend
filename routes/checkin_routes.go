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
	apiPatients := api.Group("/checkins", middleware.RoleMiddleware("patient"))
	{
		apiPatients.POST("/create", checkInController.CreateCheckIn)
		apiPatients.GET("/:id", checkInController.GetCheckIn)
		apiPatients.GET("", checkInController.GetAllCheckIns)
		apiPatients.PUT("/:id", checkInController.UpdateCheckIn)
	}
}
