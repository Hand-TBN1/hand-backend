package routes

import (
	"github.com/Hand-TBN1/hand-backend/controller"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterConsultationRoutes(router *gin.Engine, db *gorm.DB) {
	consultationService := &services.ConsultationHistoryService{DB: db}
	consultationController := &controller.ConsultationHistoryController{ConsultationHistoryService: consultationService}
	api := router.Group("/api")
	{	
		api.GET("/consultations/:user_id", consultationController.GetAllUserConsultationHistory)
	}
}
