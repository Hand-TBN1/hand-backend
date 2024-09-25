package models

import (
	"time"

	"github.com/google/uuid"
)

type BookedSchedule struct {
    ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    UserID          uuid.UUID `gorm:"type:uuid;not null"`
    TherapistID     uuid.UUID `gorm:"type:uuid;not null"`
    AppointmentDate time.Time
    Status          string
    UpdatedAt       time.Time

    // Associations
    User       User       `gorm:"foreignKey:UserID"`
    Therapist  User       `gorm:"foreignKey:TherapistID"`
}
