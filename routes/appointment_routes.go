package routes

import (
	"github.com/Hand-TBN1/hand-backend/controller"
	"github.com/Hand-TBN1/hand-backend/middleware"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterAppointmentRoutes sets up the routes for managing appointments.
func RegisterAppointmentRoutes(router *gin.Engine, db *gorm.DB) {
	appointmentService := &services.AppointmentService{DB: db}
	appointmentController := &controller.AppointmentController{AppointmentService: appointmentService}

	api := router.Group("/api/appointment")
	{
		// Ensure the user is authenticated to create appointments
		api.POST("/create-appointment", middleware.RoleMiddleware("patient"), appointmentController.CreateAppointment)
	}
}
