package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/Hand-TBN1/hand-backend/utilities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateJournalDTO struct {
	Content string `json:"content" binding:"required"`
}

type UpdateJournalDTO struct {
	Content string `json:"content" binding:"required"`
}

type JournalResponseDTO struct {
    ID      uuid.UUID `json:"id"`
    Content string    `json:"content"`
    UserID  uuid.UUID `json:"user_id"`
}

type JournalController struct {
	JournalService *services.JournalService
}

func (ctrl *JournalController) GetUserJournals(c *gin.Context) {
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

    dateParam := c.Query("date")
    var date *time.Time
    if dateParam != "" {
        parsedDate, err := time.Parse("2006-01-02", dateParam)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD."})
            return
        }

        loc, _ := time.LoadLocation("Asia/Jakarta") 
        localStartOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, loc)
        startOfDayUTC := localStartOfDay.UTC()
        date = &startOfDayUTC
    }

    journals, apiErr := ctrl.JournalService.GetUserJournals(userUUID, date)
    if apiErr != nil {
        c.JSON(apiErr.HttpStatus, apiErr)
        return
    }

    // Map journals to the response DTO
    var response []JournalResponseDTO
    for _, journal := range journals {
        response = append(response, JournalResponseDTO{
            ID:      journal.ID,
            Content: journal.Content,
            UserID:  journal.UserID,
        })
    }

    // Return the mapped journals
    c.JSON(http.StatusOK, response)
}


func (ctrl *JournalController) CreateJournal(c *gin.Context) {
	var dto CreateJournalDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	log.Printf("Received content: %v", dto.Content)

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
	log.Printf("Parsed user UUID: %v", userUUID)


	journal := models.Journal{
		ID:        uuid.New(),
		Content:   dto.Content,
		UserID:    userUUID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	log.Printf("Generated Journal ID: %v", journal.ID)

	apiErr := ctrl.JournalService.CreateJournal(&journal)
	if apiErr != nil {
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Journal created successfully"})
}
