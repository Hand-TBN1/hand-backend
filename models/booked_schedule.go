package models

import (
	"time"

	"github.com/google/uuid"
)

type BookedScheduleStatus string

const (
    Success  BookedScheduleStatus = "success"
    Canceled BookedScheduleStatus = "cancelled"
)

type BookedSchedule struct {
    ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    PatientID      uuid.UUID `gorm:"type:uuid;not null;foreignKey:PatientID"`
    TherapistID    uuid.UUID `gorm:"type:uuid;not null;foreignKey:TherapistID"`
    AppointmentDate time.Time
    Status         BookedScheduleStatus `gorm:"type:enum('success', 'canceled')"`
    UpdatedAt      time.Time
}
