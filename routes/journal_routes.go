package routes

import (
	"github.com/Hand-TBN1/hand-backend/controller"
	"github.com/Hand-TBN1/hand-backend/middleware"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterJournalRoutes(router *gin.Engine, db *gorm.DB) {
	journalService := &services.JournalService{DB: db}
	journalController := &controller.JournalController{JournalService: journalService}

	api := router.Group("/api")
	{
		journalRoutes := api.Group("/journals", middleware.RoleMiddleware("patient"))
		{
			journalRoutes.GET("", journalController.GetUserJournals)
			journalRoutes.POST("", journalController.CreateJournal)
		}
	}
}
