package models

import (
	"time"

	"github.com/google/uuid"
)

type AppointmentScheduleStatus string

const (
    Success  AppointmentScheduleStatus = "success"
    Canceled AppointmentScheduleStatus = "canceled"
)


type Appointment struct {
    ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    UserID          uuid.UUID `gorm:"type:uuid;not null"`
    TherapistID     uuid.UUID `gorm:"type:uuid;not null"`
    Type            ConsultationType `gorm:"type:consultation_enum"`
    Status         AppointmentScheduleStatus `gorm:"type:appointment_schedule_status_enum"`
    AppointmentDate time.Time
    CreatedAt       time.Time
    UpdatedAt       time.Time

    // Associations
    User       User  `gorm:"foreignKey:UserID"`
    Therapist  User  `gorm:"foreignKey:TherapistID"`
    ConsultationHistory  ConsultationHistory  `gorm:"foreignKey:AppointmentID"`
}