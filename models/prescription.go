package models

import (
	"github.com/google/uuid"
)

type Prescription struct {
    ID                   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    ConsultationHistoryID uuid.UUID `gorm:"type:uuid;not null;foreignKey:ConsultationHistoryID"`
    MedicationID          uuid.UUID `gorm:"type:uuid;not null;foreignKey:MedicationID"`
    Dosage                string
}
