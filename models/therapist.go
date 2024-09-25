package models

import (
	"time"

	"github.com/google/uuid"
)

type Therapist struct {
    ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    UserID          uuid.UUID `gorm:"type:uuid;not null;foreignKey:UserID"`
    Specialization  string
    AppointmentRate int64
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
