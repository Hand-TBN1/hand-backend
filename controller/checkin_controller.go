package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/Hand-TBN1/hand-backend/services"
)

type CreateCheckInDTO struct {
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	MoodScore string    `json:"mood_score" binding:"required"`
	Notes     string    `json:"notes"`
}

type UpdateCheckInDTO struct {
	MoodScore string `json:"mood_score"`
	Notes     string `json:"notes"`
}

type CheckInResponseDTO struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	MoodScore  string    `json:"mood_score"`
	Notes      string    `json:"notes"`
	CheckInDate time.Time `json:"check_in_date"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type CheckInController struct {
	CheckInService *services.CheckInService
}

func (ctrl *CheckInController) CreateCheckIn(c *gin.Context) {
	var createDTO CreateCheckInDTO
	if err := c.ShouldBindJSON(&createDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	checkIn := models.CheckIn{
		ID:         uuid.New(),
		UserID:     createDTO.UserID,
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
	id := c.Param("id")
	var updateDTO UpdateCheckInDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	checkIn := models.CheckIn{
		MoodScore:  updateDTO.MoodScore,
		Notes:      updateDTO.Notes,
		UpdatedAt:  time.Now(),
	}

	if err := ctrl.CheckInService.UpdateCheckIn(id, checkIn); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
