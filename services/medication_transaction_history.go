package services

import (


	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/Hand-TBN1/hand-backend/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MedicationTransactionHistoryService struct {
	DB *gorm.DB
}

type CheckoutItem struct {
	MedicationID uuid.UUID
	Price int64
	Quantity  int 
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
func (service *MedicationTransactionHistoryService) CreateMedicationTransaction(transaction *models.MedicationHistoryTransaction, items []dto.CheckoutItem) error {
    // Start transaction
    db := service.DB.Begin()
    if err := db.Create(&transaction).Error; err != nil {
        db.Rollback()
        return err
    }

    // Create each medication history item
    for _, item := range items {
        mhItem := models.MedicationHistoryItem{
            TransactionID:  transaction.ID,
            MedicationID:   item.MedicationID,
            Quantity:       item.Quantity,
            PricePerItem:   item.Price,
        }
        if err := db.Create(&mhItem).Error; err != nil {
            db.Rollback()
            return err
        }
    }

    db.Commit()
    return nil
}

