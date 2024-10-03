package services

import (


	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatService struct {
	DB *gorm.DB
}

func (s *ChatService) GetMessagesByRoomID(roomID string) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage

	parsedID, err := uuid.Parse(roomID)
	if err != nil {
		return nil, err
	}
	// Perform the query to retrieve all messages from the specified chat room.
	err = s.DB.Where("chat_room_id = ?", parsedID).Order("sent_at asc").Find(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}


func (s *ChatService) IsUserInRoom(userID, roomID uuid.UUID) bool {
	var chatRoom models.ChatRoom
	if err := s.DB.Where("id = ? AND (first_user_id = ? OR second_user_id = ?)", roomID, userID, userID).First(&chatRoom).Error; err != nil {
		return false
	}
	return true
}

func (s *ChatService) GetChatRoomsWithMessages(userUUID uuid.UUID) ([]models.ChatRoom, error) {
	var chatRooms []models.ChatRoom

	err := s.DB.
		Model(&models.ChatRoom{}).
		Joins("JOIN chat_messages ON chat_messages.chat_room_id = chat_rooms.id").
		Where("chat_rooms.first_user_id = ? OR chat_rooms.second_user_id = ?", userUUID, userUUID). // Ensure user is part of the chat room
		Where("chat_rooms.type = ?", "consultation"). // Filter for rooms of type 'consultation'
		Preload("FirstUser").
		Preload("SecondUser").
		Group("chat_rooms.id").
		Having("COUNT(chat_messages.id) > 0"). // Ensure rooms with at least one message
		Find(&chatRooms).Error

	if err != nil {
		return nil, err
	}

	return chatRooms, nil
}