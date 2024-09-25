package models

import (
	"time"

	"github.com/google/uuid"
)

type CheckIn struct {
    ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    UserID      uuid.UUID `gorm:"type:uuid;not null"`
    MoodScore  string
    Notes      string
    CheckInDate time.Time
    CreatedAt  time.Time
    UpdatedAt  time.Time

    // Associations
    User       User  `gorm:"foreignKey:UserID"`
}
