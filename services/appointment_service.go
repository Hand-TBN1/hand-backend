package services

import (
	"github.com/Hand-TBN1/hand-backend/models"
	"gorm.io/gorm"
)

type AppointmentService struct {
	DB *gorm.DB
}

// CreateAppointment creates a new appointment for a user
func (service *AppointmentService) CreateAppointment(appointment *models.Appointment) error {	
	return service.DB.Create(appointment).Error
}
