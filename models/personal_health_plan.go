package models

import (
	"time"

	"github.com/google/uuid"
)

type PersonalHealthPlan struct {
    ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    PatientID       uuid.UUID `gorm:"type:uuid;not null;foreignKey:PatientID"`
    PlanName        string
    PlanDescription string
    ReminderSchedule string
    CreatedAt       time.Time
    UpdatedAt       time.Time

    // Associations
    Patient          User              `gorm:"foreignKey:PatientID"`
}
