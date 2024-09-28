package models

import (
	"time"

	"github.com/google/uuid"
)


type MedicationHistoryTransaction struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID          uuid.UUID `gorm:"type:uuid;not null"` 
	User            User      `gorm:"foreignKey:UserID"` 
	TotalPrice      int64     `gorm:"not null"`    
	TransactionDate time.Time `gorm:"not null"` 
	CreatedAt       time.Time
	UpdatedAt       time.Time

	// Associations
	Items []MedicationHistoryItem `gorm:"foreignKey:TransactionID"` 
}

type MedicationHistoryItem struct {
	ID             uuid.UUID                   `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	TransactionID  uuid.UUID                   `gorm:"type:uuid;not null"` 
	Transaction    MedicationHistoryTransaction `gorm:"foreignKey:TransactionID"`  
	MedicationID   uuid.UUID                   `gorm:"type:uuid;not null"`
	Medication     Medication                  `gorm:"foreignKey:MedicationID"`
	Quantity       int                         `gorm:"not null"`
	PricePerItem   int64                       `gorm:"not null"`
	CreatedAt      time.Time
}
