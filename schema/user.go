package schema

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"type:varchar(320);unique;not null;index:,type:hash"`
	Name      string         `json:"name" gorm:"type:varchar(50);not null"`
	Password  string         `json:"-" gorm:"type:char(60);not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
