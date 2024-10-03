package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/Hand-TBN1/hand-backend/utilities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TherapistController struct {
	TherapistService           *services.TherapistService
	AppointmentService         *services.AppointmentService
	ConsultationHistoryService *services.ConsultationHistoryService
	PrescriptionService        *services.PrescriptionService
}

type CreateTherapistDTO struct {
	Name             string               `json:"name" binding:"required"`
	Email            string               `json:"email" binding:"required"`
	PhoneNumber      string               `json:"phone_number" binding:"required"`
	Password         string               `json:"password" binding:"required"`
	Location         string               `json:"location" binding:"required"`
	Specialization   string               `json:"specialization" binding:"required"`
	Consultation     models.ConsultationType `json:"consultation" binding:"required"`
	AppointmentRate  int64                `json:"appointment_rate" binding:"required"`
}

func (ctrl *TherapistController) GetTherapistsFiltered(c *gin.Context) {
	consultationType := c.Query("consultation")
	location := c.Query("location")
	dateStr := c.Query("date")

	var date time.Time
	if dateStr != "" {
		var err error
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
	}

	therapists, apiErr := ctrl.TherapistService.GetTherapistsFiltered(consultationType, location, date)
	if apiErr != nil {
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	c.JSON(http.StatusOK, therapists)
}

func (ctrl *TherapistController) CreateTherapist(c *gin.Context) {
	var createTherapistDTO CreateTherapistDTO

	if err := c.ShouldBindJSON(&createTherapistDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	var existingUser models.User
	if err := ctrl.TherapistService.DB.Where("email = ?", createTherapistDTO.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	hashedPassword, err := utilities.HashPassword(createTherapistDTO.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		ID:               uuid.New(),
		Name:             createTherapistDTO.Name,
		Email:            createTherapistDTO.Email,
		PhoneNumber:      createTherapistDTO.PhoneNumber,
		Password:         hashedPassword, 
		Role:             models.RoleTherapist,
		IsMobileVerified: false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	apiErr := ctrl.TherapistService.AddUser(&user)
	if apiErr != nil {
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	therapist := models.Therapist{
		ID:              uuid.New(),
		UserID:          user.ID, 
		Location:        createTherapistDTO.Location,
		Specialization:  createTherapistDTO.Specialization,
		Consultation:    createTherapistDTO.Consultation,
		AppointmentRate: createTherapistDTO.AppointmentRate,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	apiErr = ctrl.TherapistService.AddTherapist(&therapist)
	if apiErr != nil {
		c.JSON(apiErr.HttpStatus, gin.H{"error": apiErr.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Therapist created successfully"})
}

func (ctrl *TherapistController) UpdateAvailability(c *gin.Context) {
	claims, exists := c.Get("claims")
	var availabilityDTO map[string]interface{}
	
	if !exists {
		c.JSON(http.StatusUnauthorized, apierror.NewApiErrorBuilder().
			WithStatus(http.StatusUnauthorized).
			WithMessage("Unauthorized access").
			Build())
		return
	}

	userClaims := claims.(*utilities.Claims)

	therapistUUID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, apierror.NewApiErrorBuilder().
			WithStatus(http.StatusBadRequest).
			WithMessage("Invalid user ID in token").
			Build())
		return
	}

	if err := c.ShouldBindJSON(&availabilityDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	if _, ok := availabilityDTO["is_available"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'is_available' is required"})
		return
	}
	
	var finalDTO struct {
		Date        string `json:"date"`
		IsAvailable bool   `json:"is_available"`
	}
	
	finalDTO.Date = availabilityDTO["date"].(string)
	finalDTO.IsAvailable = availabilityDTO["is_available"].(bool)
	date, err := time.Parse("2006-01-02",finalDTO.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	apiErr := ctrl.TherapistService.UpdateAvailabilityByDate(therapistUUID.String(), date, finalDTO.IsAvailable)
	if apiErr != nil {
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Availability updated successfully"})
}


func (ctrl *TherapistController) GetTherapistDetails(c *gin.Context) {
	therapistID := c.Param("id")

	therapist, err := ctrl.TherapistService.GetTherapistDetails(therapistID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Respond with both Therapist and User details
	c.JSON(http.StatusOK, gin.H{
		"therapist": gin.H{
			"name":        therapist.User.Name,
			"email":       therapist.User.Email,
			"phone_number": therapist.User.PhoneNumber,
			"image_url":   therapist.User.ImageURL,
			"role":        therapist.User.Role,
			"location":    therapist.Location,
			"specialization": therapist.Specialization,
			"consultation_type": therapist.Consultation,
			"appointment_rate": therapist.AppointmentRate,
		},
	})
}


// GetTherapistSchedule - Fetch available schedule
func (ctrl *TherapistController) GetTherapistSchedule(c *gin.Context) {
	therapistID := c.Param("id")
	date := c.Query("date") // Optional
	consultationType := c.Query("type") // online/offline

	schedules, err := ctrl.TherapistService.GetAvailableSchedules(therapistID, date, consultationType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Schedule not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"schedules": schedules})
}

func (ctrl *TherapistController) GetTherapistAppointments(c *gin.Context) {
	// Extract therapist details from claims (JWT token)
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userClaims := claims.(*utilities.Claims)

	// Fetch appointments for the therapist
	appointments, err := ctrl.AppointmentService.GetAppointmentsByTherapistID(userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments"})
		return
	}

	// Return the appointments
	c.JSON(http.StatusOK, appointments)
}

// AddPrescriptionAndMedication allows therapists to add a prescription and medication after the appointment
func (ctrl *TherapistController) AddPrescriptionAndMedication(c *gin.Context) {
	appointmentID, err := uuid.Parse(c.Param("appointmentID"))
	fmt.Println(appointmentID)
	if err != nil {
		apiErr := apierror.NewApiErrorBuilder().
			WithStatus(http.StatusBadRequest).
			WithMessage("Invalid appointment ID").
			Build()
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	var req struct {
		Conclusion  string `json:"conclusion"`
		Medications []struct {
			MedicationID uuid.UUID `json:"medication_id"`
			Dosage       string    `json:"dosage"`
			Quantity    string `json:"quantity"`
		} `json:"medications"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := apierror.NewApiErrorBuilder().
			WithStatus(http.StatusBadRequest).
			WithMessage(apierror.ErrInvalidInput).
			Build()
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	// Save the consultation history
	consultationHistory := models.ConsultationHistory{
		ID:         uuid.New(),
		AppointmentID:    appointmentID,
		Conclusion:       req.Conclusion,
		ConsultationDate: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	fmt.Println(consultationHistory);
	fmt.Println("Service", ctrl.ConsultationHistoryService)
	fmt.Println("tes", ctrl.ConsultationHistoryService.DB)

	if err := ctrl.ConsultationHistoryService.CreateConsultationHistory(&consultationHistory); err != nil {
		apiErr := apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage(apierror.ErrInternalServerError).
			Build()
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	// Save each medication with the associated consultation history
	for _, med := range req.Medications {
		prescription := models.Prescription{
			ConsultationHistoryID: consultationHistory.ID,
			MedicationID:          med.MedicationID,
			Dosage:                med.Dosage,
			Quantity: 				med.Quantity,
		}
		if err := ctrl.PrescriptionService.CreatePrescription(&prescription); err != nil {
			apiErr := apierror.NewApiErrorBuilder().
				WithStatus(http.StatusInternalServerError).
				WithMessage("Failed to save prescription").
				Build()
			c.JSON(apiErr.HttpStatus, apiErr)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prescription and medication saved successfully"})
}
