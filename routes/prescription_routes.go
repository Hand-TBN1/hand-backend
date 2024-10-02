package routes

import (
	"github.com/Hand-TBN1/hand-backend/controller"
	"github.com/Hand-TBN1/hand-backend/middleware"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterPrescriptionRoutes(router *gin.Engine, db *gorm.DB) {
	// Create the necessary services
	therapistService := &services.TherapistService{DB: db}
	appointmentService := &services.AppointmentService{DB: db}
	consultationHistoryService := &services.ConsultationHistoryService{DB: db}
	prescriptionService := &services.PrescriptionService{DB: db}

	therapistController := &controller.TherapistController{
		TherapistService:             therapistService,
		AppointmentService:           appointmentService,
		ConsultationHistoryService:   consultationHistoryService,
		PrescriptionService:          prescriptionService,
	}

	// Define API routes
	api := router.Group("/api")
	{
		therapistRoutes := api.Group("/api/therapists")
		therapistRoutes.Use(middleware.RoleMiddleware("therapist"))
		{
			therapistRoutes.POST("/appointment/:appointmentID/prescription", therapistController.AddPrescriptionAndMedication)
		}
	}
}
