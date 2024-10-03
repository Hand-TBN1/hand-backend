package models

import (
	"github.com/google/uuid"
)

type Prescription struct {
	ID                  uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	ConsultationHistoryID uuid.UUID   `gorm:"type:uuid;not null"`  
	MedicationID        uuid.UUID   `gorm:"type:uuid;not null"`
	Dosage              string      `json:"dosage"`
	Quantity 			string 		`json:"quantity"`

	// Associations
	ConsultationHistory ConsultationHistory `gorm:"foreignKey:ConsultationHistoryID"`
	Medication          Medication          `gorm:"foreignKey:MedicationID"`
}
