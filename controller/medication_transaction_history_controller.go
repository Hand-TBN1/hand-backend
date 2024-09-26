package controller

import (
	"net/http"
	"time"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/Hand-TBN1/hand-backend/utilities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MedicationTransactionHistoryController struct {
	MedicationTransactionHistoryService *services.MedicationTransactionHistoryService
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
	var transactionRequest struct {
		Items []models.MedicationHistoryItem `json:"items"`
	}

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

	transaction := models.MedicationHistoryTransaction{
		UserID:         uuid.MustParse(userClaims.UserID),
		TotalPrice:     calculateTotalPrice(transactionRequest.Items),
		TransactionDate: time.Now(),
		Items:          transactionRequest.Items,
	}

	if err := ctrl.MedicationTransactionHistoryService.CreateMedicationTransaction(&transaction); err != nil {
		c.JSON(http.StatusInternalServerError, apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage(err.Error()).
			Build())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction successful"})
}

// Helper function to calculate the total price
func calculateTotalPrice(items []models.MedicationHistoryItem) int64 {
	var total int64
	for _, item := range items {
		total += item.PricePerItem * int64(item.Quantity)
	}
	return total
}
