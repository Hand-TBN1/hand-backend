package models

import (
	"time"

	"github.com/google/uuid"
)

type MediaType string

const (
	Article  MediaType = "article"
	Video MediaType = "video"
)

type Media struct {
    ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Type			 MediaType `gorm:"type:media_enum"`
    Title            string `gorm:"type:varchar(255);not null"`
	Content          string `gorm:"type:text;not null"`
	ThumbnailURL     string `json:"image_url"`
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
