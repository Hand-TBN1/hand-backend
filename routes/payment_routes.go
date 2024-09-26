package routes

import (
	"github.com/Hand-TBN1/hand-backend/controller"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
)

func SetupPaymentRoutes(router *gin.Engine, paymentService *services.PaymentService) {
	paymentController := &controller.PaymentController{
		PaymentService: paymentService,
	}

	api := router.Group("/api")
	{
		api.POST("/payment", paymentController.CreatePaymentTransaction)
	}
}
