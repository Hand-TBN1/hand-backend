package models

import (
	"time"

	"github.com/google/uuid"
)


type ChatRoomType string

const (
    ConsultationRoom ChatRoomType = "consultation"
    AnonymousRoom    ChatRoomType = "anonymous"
)

type ChatRoom struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	FirstUserID  uuid.UUID `gorm:"type:uuid;not null;foreignKey:UserID"`
	SecondUserID uuid.UUID `gorm:"type:uuid;not null;foreignKey:UserID"`
	IsEnd        bool
	Type        ChatRoomType  `gorm:"type:room_enum;not null"`
	CreatedAt    time.Time

	// Associations
	FirstUser  User `gorm:"foreignKey:FirstUserID"`
	SecondUser User `gorm:"foreignKey:SecondUserID"`
}
