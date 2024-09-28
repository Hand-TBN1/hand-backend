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

func (service *AppointmentService) GetAppointmentsByUserID(userID string, status string) ([]models.Appointment, error) {
    var appointments []models.Appointment

	query := service.DB.Preload("Therapist").Preload("Therapist.Therapist").Where("user_id = ?", userID)

    if status != "" {
        query = query.Where("status = ?", status)
    }

    if err := query.Find(&appointments).Error; err != nil {
        return nil, err
    }

    return appointments, nil
}