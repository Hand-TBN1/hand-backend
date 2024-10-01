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
		apiErr := apierror.NewApiErrorBuilder().
			WithStatus(http.StatusBadRequest).
			WithMessage(apierror.ErrInvalidInput).
			Build()
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	// Convert date to time.Time
	appointmentDate, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		apiErr := apierror.NewApiErrorBuilder().
			WithStatus(http.StatusBadRequest).
			WithMessage("Invalid date format").
			Build()
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	claims, exists := c.Get("claims")
	if !exists {
		apiErr := apierror.NewApiErrorBuilder().
			WithStatus(http.StatusUnauthorized).
			WithMessage(apierror.ErrUnauthorized).
			Build()
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	userClaims := claims.(*utilities.Claims)

	therapist, err := ctrl.TherapistService.GetTherapistDetails(req.TherapistID)
	if err != nil {
		apiErr := apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage(apierror.ErrInternalServerError).
			Build()
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	appointment := models.Appointment{
		ID:              uuid.New(),
		UserID:          uuid.MustParse(userClaims.UserID),
		TherapistID:     uuid.MustParse(req.TherapistID),
		AppointmentDate: appointmentDate,
		Price:           therapist.AppointmentRate,
		PaymentStatus:   models.MidtransStatusPending,
		Type:            models.ConsultationType(req.ConsultationType),
		Status:          models.Success, 
		CreatedAt:       time.Now(),
	}

	// Save appointment to the database
	if err := ctrl.AppointmentService.CreateAppointment(&appointment); err != nil {
		apiErr := apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to create appointment").
			Build()
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	// Create payment for the appointment
	paymentResponse, err := ctrl.PaymentService.CreatePayment(appointment.ID.String(), appointment.Price)
	if err != nil {
		apiErr := apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage(err.Error()).
			Build()
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	// Respond with success and payment redirect URL
	c.JSON(http.StatusOK, gin.H{
		"message":      "Appointment created successfully",
		"appointment_id": appointment.ID,
		"payment_status": appointment.PaymentStatus,
		"redirect_url": paymentResponse.RedirectURL, 
	})
}


func (ctrl *AppointmentController) GetAppointmentHistory(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		apiErr := apierror.NewApiErrorBuilder().
			WithStatus(http.StatusUnauthorized).
			WithMessage(apierror.ErrUnauthorized).
			Build()
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	userClaims := claims.(*utilities.Claims)

	status := c.Query("status")

	appointments, err := ctrl.AppointmentService.GetAppointmentsByUserID(userClaims.UserID, status)
	if err != nil {
		apiErr := apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to fetch appointments").
			Build()
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	var result []gin.H
	for _, appointment := range appointments {
		therapist := appointment.Therapist
		therapistDetails := gin.H{
			"name":      therapist.Name,
			"image_url": therapist.ImageURL,
			"location":therapist.Therapist.Location,
		}

		result = append(result, gin.H{
			"appointment_id":   appointment.ID,
			"therapist":        therapistDetails,
			"price":            appointment.Price,
			"appointment_date": appointment.AppointmentDate,
			"type":             appointment.Type,
			"status":           appointment.PaymentStatus,
			"payment_status":   appointment.PaymentStatus,
		})
	}

	c.JSON(http.StatusOK, result)
}
