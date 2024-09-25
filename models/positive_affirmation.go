package models

import (
	"time"

	"github.com/google/uuid"
)

type PositiveAffirmation struct {
    ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    PatientID         uuid.UUID `gorm:"type:uuid;not null;foreignKey:PatientID"`
    AffirmationContent string
    SentAt            time.Time

    // Associations
    Patient            User              `gorm:"foreignKey:PatientID"`
}
