package models

import (
	"time"

	"github.com/google/uuid"
)

type ConsultationHistory struct {
    ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    PatientID        uuid.UUID `gorm:"type:uuid;not null;foreignKey:PatientID"`
    TherapistID      uuid.UUID `gorm:"type:uuid;not null;foreignKey:TherapistID"`
    Conclusion       string
    Price            int
    Prescription     string
    ConsultationDate time.Time
    CreatedAt        time.Time
    UpdatedAt        time.Time

    // Associations
    Patient          User       `gorm:"foreignKey:PatientID"`
    Therapist        Therapist  `gorm:"foreignKey:TherapistID"`
}
