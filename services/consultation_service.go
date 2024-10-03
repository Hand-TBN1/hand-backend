package services

import (
	"github.com/Hand-TBN1/hand-backend/models"
	"gorm.io/gorm"
)

type ConsultationHistoryService struct {
	DB *gorm.DB
}

func (service *ConsultationHistoryService) CreateConsultationHistory(history *models.ConsultationHistory) error {
	if err := service.DB.Create(history).Error; err != nil {
		return err
	}
	return nil
}
