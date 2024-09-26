package services

import (
	"errors"
	"time"

	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CheckInService struct {
	DB *gorm.DB
}

func (service *CheckInService) CreateCheckIn(checkIn *models.CheckIn) error {
	if err := service.DB.Create(checkIn).Error; err != nil {
		return err
	}
	return nil
}

func (service *CheckInService) GetCheckIn(id string) (*models.CheckIn, error) {
	var checkIn models.CheckIn
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	if err := service.DB.First(&checkIn, "id = ?", parsedID).Error; err != nil {
		return nil, err
	}
	return &checkIn, nil
}

func (service *CheckInService) GetAllCheckIns() ([]models.CheckIn, error) {
	var checkIns []models.CheckIn
	if err := service.DB.Find(&checkIns).Error; err != nil {
		return nil, err
	}
	return checkIns, nil
}

func (service *CheckInService) UpdateCheckIn(id string, newCheckIn models.CheckIn) error {
	var checkIn models.CheckIn
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	if err := service.DB.First(&checkIn, "id = ?", parsedID).Error; err != nil {
		return errors.New("check-in not found")
	}
	checkIn.MoodScore = newCheckIn.MoodScore
	checkIn.Notes = newCheckIn.Notes
	checkIn.UpdatedAt = time.Now()

	if err := service.DB.Save(&checkIn).Error; err != nil {
		return err
	}
	return nil
}