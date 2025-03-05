package services

import (
	"time"

	"github.com/Hand-TBN1/hand-backend/models"
	"gorm.io/gorm"
)

type AvailabilityService struct {
	DB *gorm.DB
}

func (service *AvailabilityService) GetBlockedAvailabilityByTherapistID(therapistID string, currentTime time.Time) ([]models.Availability, error) {
	var blockedDates []models.Availability

	err := service.DB.
		Where("therapist_id = ? AND date >= ? AND is_available = false", therapistID, currentTime).
		Order("date asc").
		Find(&blockedDates).Error

	if err != nil {
		return nil, err
	}

	return blockedDates, nil
}
