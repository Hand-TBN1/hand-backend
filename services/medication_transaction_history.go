package services

import (
	"errors"
	"time"

	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MedicationTransactionHistoryService struct {
	DB *gorm.DB
}

// GetMedicationHistoryByUserID fetches the medication transaction history for a specific user
func (service *MedicationTransactionHistoryService) GetMedicationHistoryByUserID(userID uuid.UUID) ([]models.MedicationHistoryTransaction, error) {
	var history []models.MedicationHistoryTransaction
	if err := service.DB.Preload("Items.Medication").Where("user_id = ?", userID).Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

// CreateMedicationTransaction creates a new medication transaction for a user
func (service *MedicationTransactionHistoryService) CreateMedicationTransaction(transaction *models.MedicationHistoryTransaction) error {
	// Ensure user exists
	var user models.User
	if err := service.DB.First(&user, "id = ?", transaction.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("usercvxvxcv not found")
		}
		return err
	}

	// Set creation and update timestamps
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()

	// Save transaction and related items in one atomic transaction
	return service.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
		for _, item := range transaction.Items {
			item.TransactionID = transaction.ID
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
