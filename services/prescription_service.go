package services

import (
	"github.com/Hand-TBN1/hand-backend/models"
	"gorm.io/gorm"
)

type PrescriptionService struct {
	DB *gorm.DB
}

func (service *PrescriptionService) CreatePrescription(prescription *models.Prescription) error {
	return service.DB.Create(prescription).Error
}
