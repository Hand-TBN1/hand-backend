package routes

import (
	"github.com/Hand-TBN1/hand-backend/controller"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupPaymentRoutes(router *gin.Engine, db *gorm.DB) {
	paymentService := &services.PaymentService{}
	appointmentService := &services.AppointmentService{DB: db}

	paymentController := &controller.PaymentController{
		PaymentService:    paymentService,
		AppointmentService: appointmentService,
	}

	api := router.Group("/api")
	{
		api.POST("/payment", paymentController.CreatePaymentTransaction)

		// Route to handle Midtrans payment notifications (webhook)
		api.POST("/payment-notification", paymentController.HandlePaymentNotification)
	}
}
