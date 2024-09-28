package models

import (
	"time"

	"github.com/google/uuid"
)

type MidtransStatus string

var (
	MidtransStatusChallenge MidtransStatus = "challenge"
	MidtransStatusSuccess   MidtransStatus = "success"
	MidtransStatusFailure   MidtransStatus = "failure"
	MidtransStatusPending   MidtransStatus = "pending"
)


type MedicationHistoryTransaction struct {
    ID              uuid.UUID              `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
    UserID          uuid.UUID              `gorm:"type:uuid;not null" json:"user_id"`
    User            User                   `gorm:"foreignKey:UserID" json:"user"`
    TotalPrice      int64                  `gorm:"not null" json:"total_price"`
    PaymentStatus   MidtransStatus         `gorm:"type:midtrans_status;not null" json:"payment_status"`
    TransactionDate time.Time              `gorm:"not null" json:"transaction_date"`
    CreatedAt       time.Time              `json:"created_at"`
    UpdatedAt       time.Time              `json:"updated_at"`
    Items           []MedicationHistoryItem `gorm:"foreignKey:TransactionID" json:"items"`
}

type MedicationHistoryItem struct {
    ID             uuid.UUID                `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
    TransactionID  uuid.UUID                `gorm:"type:uuid;not null" json:"transaction_id"`
    Transaction    MedicationHistoryTransaction `gorm:"foreignKey:TransactionID" json:"-"`
    MedicationID   uuid.UUID                `gorm:"type:uuid;not null" json:"medication_id"`
    Medication     Medication               `gorm:"foreignKey:MedicationID" json:"medication"`
    Name 	string 							`json:"name"`
	Quantity       int                      `gorm:"not null" json:"quantity"`
    PricePerItem   int64                    `gorm:"not null" json:"price_per_item"`
    CreatedAt      time.Time                `json:"created_at"`
}
