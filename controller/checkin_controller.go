package controller

import (
	"net/http"
	"time"

	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/Hand-TBN1/hand-backend/utilities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateCheckInDTO struct {
	MoodScore int       `json:"mood_score" binding:"required"` 
	Notes     string    `json:"notes"`
}

type UpdateCheckInDTO struct {
	MoodScore int    `json:"mood_score"` 
	Notes     string `json:"notes"`
}

type CheckInResponseDTO struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	MoodScore  int       `json:"mood_score"` 
	Notes      string    `json:"notes"`
	CheckInDate time.Time `json:"check_in_date"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type CheckInController struct {
	CheckInService *services.CheckInService
}

func (ctrl *CheckInController) CreateCheckIn(c *gin.Context) {
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

	var createDTO CreateCheckInDTO
	if err := c.ShouldBindJSON(&createDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	checkIn := models.CheckIn{
		ID:         uuid.New(),
		UserID:     userUUID, 
		MoodScore:  createDTO.MoodScore,
		Notes:      createDTO.Notes,
		CheckInDate: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := ctrl.CheckInService.CreateCheckIn(&checkIn); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CheckInResponseDTO{
		ID:         checkIn.ID,
		UserID:     checkIn.UserID,
		MoodScore:  checkIn.MoodScore,
		Notes:      checkIn.Notes,
		CheckInDate: checkIn.CheckInDate,
		CreatedAt:  checkIn.CreatedAt,
		UpdatedAt:  checkIn.UpdatedAt,
	})
}

func (ctrl *CheckInController) GetCheckIn(c *gin.Context) {
	id := c.Param("id")
	checkIn, err := ctrl.CheckInService.GetCheckIn(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Check-in not found"})
		return
	}
	c.JSON(http.StatusOK, CheckInResponseDTO{
		ID:         checkIn.ID,
		UserID:     checkIn.UserID,
		MoodScore:  checkIn.MoodScore, 
		Notes:      checkIn.Notes,
		CheckInDate: checkIn.CheckInDate,
		CreatedAt:  checkIn.CreatedAt,
		UpdatedAt:  checkIn.UpdatedAt,
	})
}

func (ctrl *CheckInController) GetAllCheckIns(c *gin.Context) {
	checkIns, err := ctrl.CheckInService.GetAllCheckIns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var response []CheckInResponseDTO
	for _, checkIn := range checkIns {
		response = append(response, CheckInResponseDTO{
			ID:         checkIn.ID,
			UserID:     checkIn.UserID,
			MoodScore:  checkIn.MoodScore,
			Notes:      checkIn.Notes,
			CheckInDate: checkIn.CheckInDate,
			CreatedAt:  checkIn.CreatedAt,
			UpdatedAt:  checkIn.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK, response)
}

func (ctrl *CheckInController) UpdateCheckIn(c *gin.Context) {
	// Extract userID from JWT claims
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

	// Get today's date
	today := time.Now().UTC().Format("2006-01-02")

	var updateDTO UpdateCheckInDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Step 1: Find the check-in by userID and today's date in UTC
	checkIn, err := ctrl.CheckInService.FindCheckInByUserIDAndDate(userUUID, today)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No check-in found for today"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Step 2: Update the check-in's mood score and notes
	checkIn.MoodScore = updateDTO.MoodScore
	checkIn.Notes = updateDTO.Notes
	checkIn.UpdatedAt = time.Now().UTC()

	// Step 3: Save the updated check-in
	if err := ctrl.CheckInService.UpdateCheckInByUserIDAndDate(userUUID, today, *checkIn); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Check-in updated successfully"})
}
