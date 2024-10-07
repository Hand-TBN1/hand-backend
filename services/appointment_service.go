package services

import (
	"time"
	"github.com/Hand-TBN1/hand-backend/models"
	"gorm.io/gorm"
	"github.com/google/uuid"
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
    query = query.Order("created_at desc")

    if err := query.Find(&appointments).Error; err != nil {
        return nil, err
    }

    return appointments, nil
}

// UpdatePaymentAndAppointmentStatus updates the payment and appointment status
func (service *AppointmentService) UpdatePaymentAndAppointmentStatus(orderID string, status string) error {
	var appointment models.Appointment

	// Find the appointment by order ID
	if err := service.DB.Where("id = ?", orderID).First(&appointment).Error; err != nil {
		return err
	}

	// Update payment and appointment status based on Midtrans transaction status
	switch status {
	case "settlement":
		appointment.PaymentStatus = models.MidtransStatusSuccess
		appointment.Status = models.Success
	case "expire":
		appointment.PaymentStatus = models.MidtransStatusFailure
		appointment.Status = models.Canceled
	case "deny", "cancel", "failure":
		appointment.PaymentStatus = models.MidtransStatusFailure
		appointment.Status = models.Canceled
	default:
		appointment.PaymentStatus = models.MidtransStatusPending
	}

	// Save the updated appointment
	return service.DB.Save(&appointment).Error
}

// GetAppointmentsByTherapistID fetches all appointments associated with a therapist
func (service *AppointmentService) GetAppointmentsByTherapistID(therapistID string) ([]models.Appointment, error) {
	var appointments []models.Appointment

	err := service.DB.Preload("User").Where("therapist_id = ?", therapistID).Find(&appointments).Error
	if err != nil {
		return nil, err
	}

	return appointments, nil
}

func (service *AppointmentService) GetAppointmentWithUserByID(appointmentID uuid.UUID, appointment *models.Appointment) error {
    return service.DB.Preload("User").First(appointment, "id = ?", appointmentID).Error
}


func (service *AppointmentService) GetUpcomingAppointmentsByTherapistID(therapistID string, currentTime time.Time) ([]models.Appointment, error) {
    var appointments []models.Appointment

    err := service.DB.Preload("User").
        Joins("LEFT JOIN consultation_histories ON consultation_histories.appointment_id = appointments.id").
        Where("therapist_id = ?", therapistID).
        Where("consultation_histories.conclusion IS NULL OR consultation_histories.conclusion = ''").
		Where("status = ?", models.Success). 
		Where("payment_status = ?", models.MidtransStatusSuccess).
        Order("appointment_date asc").
        Find(&appointments).Error

    if err != nil {
        return nil, err
    }

    return appointments, nil
}
