package routes

import (
	"github.com/Hand-TBN1/hand-backend/controller"
	"github.com/Hand-TBN1/hand-backend/middleware"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterMedicationTransactionHistoryRoutes(router *gin.Engine, db *gorm.DB) {

	medicationTransactionHistoryService := &services.MedicationTransactionHistoryService{DB: db}

	medicationController := &controller.MedicationTransactionHistoryController{
		MedicationTransactionHistoryService: medicationTransactionHistoryService,
	}

	// Medication transaction history routes
	medicationTransactionHistoryRoutes := router.Group("/api/medication", middleware.RoleMiddleware("patient"))
	{
		medicationTransactionHistoryRoutes.GET("/history/:user_id", medicationController.GetMedicationHistoryByUserID)
		medicationTransactionHistoryRoutes.POST("/transaction", medicationController.PostMedicationTransaction)
	}
}
