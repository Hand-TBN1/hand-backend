package models

import (
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
    ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    SenderID       uuid.UUID `gorm:"type:uuid;not null;foreignKey:SenderID"`
    ChatRoomID     uuid.UUID `gorm:"type:uuid;not null;foreignKey:ChatRoomID"`
    MessageContent string
    SentAt         time.Time
}
