package routes

import (
	"github.com/Hand-TBN1/hand-backend/controller"
	"github.com/Hand-TBN1/hand-backend/middleware"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterTherapistRoutes(router *gin.Engine, db *gorm.DB) {
	therapistService := &services.TherapistService{DB: db}
	therapistController := &controller.TherapistController{TherapistService: therapistService}

	api := router.Group("/api")
	{
		therapistRoutesAdmin := api.Group("/therapists")
		therapistRoutesAdmin.Use(middleware.RoleMiddleware("admin"))
		{
			therapistRoutesAdmin.POST("/create", therapistController.CreateTherapist)
		}
		therapistRoutes := api.Group("/therapists")
		therapistRoutes.Use(middleware.RoleMiddleware("therapist", "admin")) 
		{
			therapistRoutes.PATCH("/availability", therapistController.UpdateAvailability)
			therapistRoutes.GET("/appointments", therapistController.GetTherapistAppointments)
		}

		api.GET("/therapists", therapistController.GetTherapistsFiltered)
		api.GET("/therapist/:id/details", therapistController.GetTherapistDetails)
		api.GET("/therapist/:id/schedule", therapistController.GetTherapistSchedule)
	}
}
