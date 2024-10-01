package controller

import (
	"net/http"

	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	PaymentService *services.PaymentService
	AppointmentService *services.AppointmentService
}

// CreatePaymentTransaction handles the payment transaction
func (ctrl *PaymentController) CreatePaymentTransaction(c *gin.Context) {
	var req struct {
		OrderID     string `json:"order_id"`
		GrossAmount int64  `json:"gross_amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	
	resp, err := ctrl.PaymentService.CreatePayment(req.OrderID, req.GrossAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": resp.Token, "redirect_url": resp.RedirectURL})
}


// HandlePaymentNotification processes the Midtrans notification webhook
func (ctrl *PaymentController) HandlePaymentNotification(c *gin.Context) {
	var notification map[string]interface{}
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification payload"})
		return
	}

	orderID := notification["order_id"].(string)
	transactionStatus := notification["transaction_status"].(string)

	// Update payment and appointment status based on the notification
	err := ctrl.AppointmentService.UpdatePaymentAndAppointmentStatus(orderID, transactionStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment status updated"})
}
