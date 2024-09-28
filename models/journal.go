package models

import (
	"time"

	"github.com/google/uuid"
)

type Journal struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"` 
	User      User      `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
