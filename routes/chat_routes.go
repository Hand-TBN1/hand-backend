package routes

import (
	"github.com/Hand-TBN1/hand-backend/controller"
	"github.com/Hand-TBN1/hand-backend/middleware"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterChatRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize services
	chatService := &services.ChatService{DB: db}
	chatController := &controller.ChatController{ChatService: chatService}

	api := router.Group("/api")
	apiPatients := api.Group("/room", middleware.RoleMiddleware())
	{
		apiPatients.GET("/chat",chatController.GetChatRoomsWithMessagesHandler);
		apiPatients.GET("/message/:roomId", chatController.GetMessageInRoom);
	

	}
}


