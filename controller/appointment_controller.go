package controller

import (
	"net/http"
	"time"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/Hand-TBN1/hand-backend/utilities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppointmentController struct {
	AppointmentService *services.AppointmentService
	PaymentService *services.PaymentService
	TherapistService *services.TherapistService
}

// CreateAppointment - Book an appointment with the therapist
func (ctrl *AppointmentController) CreateAppointment(c *gin.Context) {
	var req struct {
		TherapistID      string `json:"therapist_id"`
		Date             string `json:"date"` // "2024-09-29T15:00:00"
		ConsultationType string `json:"consultation_type"`
	}

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

	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userClaims := claims.(*utilities.Claims)
	defaultStatus := models.Success

	therapis , err := ctrl.TherapistService.GetTherapistDetails(userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Internal Server Error"})
		return
	}
	// Create Appointment
	appointment := models.Appointment{
		ID:              uuid.New(),
		UserID:          uuid.MustParse(userClaims.UserID),  // Get user_id directly from claims
		TherapistID:     uuid.MustParse(req.TherapistID),
		AppointmentDate: appointmentDate,
		Price:  therapis.AppointmentRate,
		Type:            models.ConsultationType(req.ConsultationType),
		Status:          defaultStatus,
		CreatedAt:       time.Now(),
	}

	if err := ctrl.AppointmentService.CreateAppointment(&appointment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create appointment"})
		return
	}

	paymentResponse, err := ctrl.PaymentService.CreatePayment(appointment.ID.String(), appointment.Price)
    if err != nil {
        c.JSON(http.StatusInternalServerError, apierror.NewApiErrorBuilder().
            WithStatus(http.StatusInternalServerError).
            WithMessage(err.Error()).
            Build())
        return
    }



	c.JSON(http.StatusOK, gin.H{"message": "Appointment created successfully", "redirect_url" : paymentResponse.RedirectURL})
}

