package controller

import (
	"net/http"
	"time"

	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppointmentController struct {
	AppointmentService *services.AppointmentService
}

// CreateAppointment - Book an appointment with the therapist
func (ctrl *AppointmentController) CreateAppointment(c *gin.Context) {
	var req struct {
		TherapistID      string `json:"therapist_id"`
		Date             string `json:"date"` // "2024-09-29T15:00:00"
		ConsultationType string `json:"consultation_type"`
	}

	// Bind JSON to request struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Convert date to time.Time
	appointmentDate, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	userID := c.MustGet("user_id").(string) // Get user ID from JWT token or session

	// Create Appointment
	appointment := models.Appointment{
		ID:              uuid.New(),
		UserID:          uuid.MustParse(userID),
		TherapistID:     uuid.MustParse(req.TherapistID),
		AppointmentDate: appointmentDate,
		Type: models.ConsultationType(req.ConsultationType),
		CreatedAt:       time.Now(),
	}

	if err := ctrl.AppointmentService.CreateAppointment(&appointment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create appointment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment created successfully"})
}
