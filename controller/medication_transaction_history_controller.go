package controller

import (
	"net/http"
	"time"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/Hand-TBN1/hand-backend/utilities"
	"github.com/Hand-TBN1/hand-backend/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MedicationTransactionHistoryController struct {
	MedicationTransactionHistoryService *services.MedicationTransactionHistoryService
	PaymentService *services.PaymentService
}


// GetMedicationHistoryByUserID retrieves medication history by user ID
func (ctrl *MedicationTransactionHistoryController) GetMedicationHistoryByUserID(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, apierror.NewApiErrorBuilder().
			WithStatus(http.StatusUnauthorized).
			WithMessage("Unauthorized access").
			Build())
		return
	}

	userClaims := claims.(*utilities.Claims)

	userUUID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, apierror.NewApiErrorBuilder().
			WithStatus(http.StatusBadRequest).
			WithMessage("Invalid user ID in token").
			Build())
		return
	}

	history, err := ctrl.MedicationTransactionHistoryService.GetMedicationHistoryByUserID(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage(err.Error()).
			Build())
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}

// PostMedicationTransaction handles a new medication purchase
func (ctrl *MedicationTransactionHistoryController) PostMedicationTransaction(c *gin.Context) {
    var transactionRequest dto.CheckoutMedicationRequest

    claims, exists := c.Get("claims")
    if !exists {
        c.JSON(http.StatusUnauthorized, apierror.NewApiErrorBuilder().
            WithStatus(http.StatusUnauthorized).
            WithMessage("Unauthorized access").
            Build())
        return
    }

    userClaims := claims.(*utilities.Claims)

    if err := c.ShouldBindJSON(&transactionRequest); err != nil {
        c.JSON(http.StatusBadRequest, apierror.NewApiErrorBuilder().
            WithStatus(http.StatusBadRequest).
            WithMessage(apierror.ErrInvalidInput).
            Build())
        return
    }

    totalPrice := calculateTotalPrice(transactionRequest.AllItem) // Adjusted to use the correct struct field
    userID := uuid.MustParse(userClaims.UserID)
    transaction := models.MedicationHistoryTransaction{
		ID : uuid.New(),
        UserID:          userID,
        TotalPrice:      totalPrice,
        TransactionDate: time.Now(),
        PaymentStatus:   models.MidtransStatusPending, // Default payment status as Pending
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
    }

    // Create transaction in database
    if err := ctrl.MedicationTransactionHistoryService.CreateMedicationTransaction(&transaction, transactionRequest.AllItem); err != nil {
        c.JSON(http.StatusInternalServerError, apierror.NewApiErrorBuilder().
            WithStatus(http.StatusInternalServerError).
            WithMessage(err.Error()).
            Build())
        return
    }

    // Initiate payment after transaction record is created
    paymentResponse, err := ctrl.PaymentService.CreatePayment(transaction.ID.String(), totalPrice)
    if err != nil {
        c.JSON(http.StatusInternalServerError, apierror.NewApiErrorBuilder().
            WithStatus(http.StatusInternalServerError).
            WithMessage(err.Error()).
            Build())
        return
    }

    // Send the payment URL to the frontend
    c.JSON(http.StatusOK, gin.H{
        "message":       "Transaction successful, proceed to payment",
        "payment_url":   paymentResponse.RedirectURL,
    })
}

// Helper function to calculate the total price
func calculateTotalPrice(items []dto.CheckoutItem) int64 {
	var total int64
	for _, item := range items {
		total += item.Price * int64(item.Quantity)
	}
	return total
}
