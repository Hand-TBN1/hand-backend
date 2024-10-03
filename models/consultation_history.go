package models

import (
	"time"

	"github.com/google/uuid"
)

type ConsultationHistory struct {
	ID               uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	AppointmentID    uuid.UUID   `gorm:"type:uuid;not null"`  
	Conclusion       string
	ConsultationDate time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Appointment		Appointment `gorm:"foreignKey:AppointmentID"`
	// Associations
	Prescription []Prescription `gorm:"foreignKey:ConsultationHistoryID"`
}
