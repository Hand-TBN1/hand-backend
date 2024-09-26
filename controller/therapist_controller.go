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

type TherapistController struct {
	TherapistService *services.TherapistService
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
	IsAvailableToday bool                 `json:"is_available_today"`
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
