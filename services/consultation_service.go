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

func (service *ConsultationHistoryService) GetConsultationHistoryByUserID(userID string) ([]models.ConsultationHistory, error) {
	var consultationHistories []models.ConsultationHistory

	// Preload the related data: Appointment -> Therapist (which is a User) -> Prescriptions -> Medication
	err := service.DB.Preload("Appointment").Preload("Appointment.Therapist").Preload("Prescription").Preload("Prescription.Medication").
		Joins("JOIN appointments ON consultation_histories.appointment_id = appointments.id").
		Where("appointments.user_id = ?", userID).
		Find(&consultationHistories).Error

	if err != nil {
		return nil, err
	}

	return consultationHistories, nil
}