package models

import (
	"time"

	"github.com/google/uuid"
)

type CheckIn struct {
    ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    PatientID  uuid.UUID `gorm:"type:uuid;not null;foreignKey:PatientID"`
    MoodScore  string
    Notes      string
    CheckInDate time.Time
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
