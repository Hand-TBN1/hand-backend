package models

import (
	"time"

	"github.com/google/uuid"
)

type BookedScheduleStatus string

const (
    Success  BookedScheduleStatus = "success"
    Canceled BookedScheduleStatus = "canceled"
)

type BookedSchedule struct {
    ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    UserID          uuid.UUID `gorm:"type:uuid;not null"`
    TherapistID     uuid.UUID `gorm:"type:uuid;not null"`
    AppointmentDate time.Time
    Status         BookedScheduleStatus `gorm:"type:booked_schedule_status_enum"`
    UpdatedAt       time.Time

    // Associations
    User       User       `gorm:"foreignKey:UserID"`
    Therapist  User       `gorm:"foreignKey:TherapistID"`
}
