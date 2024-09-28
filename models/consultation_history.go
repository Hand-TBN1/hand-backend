package models

import (
	"time"

	"github.com/google/uuid"
)

type ConsultationHistory struct {
    ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    AppointmentID    uuid.UUID `gorm:"type:uuid;not null"`
    Conclusion       string
    Price            int
    Prescription     string
    ConsultationDate time.Time
    CreatedAt        time.Time
    UpdatedAt        time.Time

}
