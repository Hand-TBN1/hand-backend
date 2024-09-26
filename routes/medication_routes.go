package routes

import (
	"github.com/Hand-TBN1/hand-backend/controller"
	"github.com/Hand-TBN1/hand-backend/middleware"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterMedicationRoutes(router *gin.Engine, db *gorm.DB) {
	medicationService := &services.MedicationService{DB: db}
	medicationController := &controller.MedicationController{MedicationService: medicationService}

	api := router.Group("/api")
	{
		adminTherapistRoutes := api.Group("/medications")
		adminTherapistRoutes.Use(middleware.RoleMiddleware("admin", "therapist"))
		{
			adminTherapistRoutes.POST("/create", medicationController.AddMedication)
			adminTherapistRoutes.PUT("/:id", medicationController.UpdateMedication)
			adminTherapistRoutes.DELETE("/:id", medicationController.DeleteMedication)
		}
		allRolesRoutes := api.Group("/medications")
		allRolesRoutes.Use(middleware.RoleMiddleware("patient", "therapist", "admin"))
		{
			allRolesRoutes.GET("", medicationController.GetMedications)  // Search by name /medications?name=panadol
		}
	}
}
