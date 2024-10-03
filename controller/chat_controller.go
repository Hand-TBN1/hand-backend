package controller

import (
	"net/http"

	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/Hand-TBN1/hand-backend/utilities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type ChatController struct {
	ChatService *services.ChatService
}

func (ctrl *ChatController) GetMessageInRoom (c *gin.Context){
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	userClaims := claims.(*utilities.Claims) 
	userUUID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in token"})
		return
	}

	
	roomId := c.Param("roomId")

	roomUUID, err := uuid.Parse(roomId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	if !ctrl.ChatService.IsUserInRoom(userUUID, roomUUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: User is not a participant in the chat room"})
		return
	}



	messages, err := ctrl.ChatService.GetMessagesByRoomID(roomId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages);



}




func (ctrl *ChatController) GetChatRoomsWithMessagesHandler(c *gin.Context) {
	// Check if claims exist
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Extract user ID from the claims
	userClaims := claims.(*utilities.Claims)
	userUUID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in token"})
		return
	}

	// Get chat rooms where the user is involved and there is at least one message
	chatRooms, err := ctrl.ChatService.GetChatRoomsWithMessages(userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Respond with the found chat rooms
	c.JSON(http.StatusOK, chatRooms)
}