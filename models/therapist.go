package models

import (
	"time"

	"github.com/google/uuid"
)

type ConsultationType string

const (
	Online  ConsultationType = "online"
	Offline ConsultationType = "offline"
	Hybrid  ConsultationType = "hybrid"
)

type Therapist struct {
	ID              uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID          uuid.UUID        `gorm:"type:uuid;not null;foreignKey:UserID"`
	Location        string
	Specialization  string
	Consultation    ConsultationType `gorm:"type:consultation_enum"`
	AppointmentRate int64
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Appointments    []Appointment    `gorm:"foreignKey:TherapistID"`
}
