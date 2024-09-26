package models

import (
	"time"

	"github.com/google/uuid"
)

type Availability struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	TherapistID     uuid.UUID `gorm:"type:uuid;not null;foreignKey:UserID"`
	Date            time.Time `gorm:"not null"`
	IsAvailable     bool      `gorm:"default:false"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
