package services

import (
	"time"

	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/google/uuid"
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
		Where("appointment_date >= ?", currentTime).
        Order("appointment_date asc").
        Find(&appointments).Error

    if err != nil {
        return nil, err
    }

    return appointments, nil
}

func (service *AppointmentService) GetAppointmentSummaryByTherapistID(therapistID string) (map[string]int, error) {
    var totalAppointments int64 = 0
    var completedAppointments int64 = 0
    var upcomingAppointments int64 = 0

    if err := service.DB.Model(&models.Appointment{}).
        Where("therapist_id = ? AND status = ? AND payment_status = ?", 
            therapistID, models.Success, models.MidtransStatusSuccess).
        Count(&totalAppointments).Error; err != nil {
        return nil, err
    }

    if err := service.DB.Model(&models.Appointment{}).
        Where("therapist_id = ? AND status = ? AND payment_status = ? AND appointment_date < ?", 
            therapistID, models.Success, models.MidtransStatusSuccess, time.Now()).
        Count(&completedAppointments).Error; err != nil {
        return nil, err
    }

    if err := service.DB.Model(&models.Appointment{}).
        Where("therapist_id = ? AND status = ? AND payment_status = ? AND appointment_date >= ?", 
            therapistID, models.Success, models.MidtransStatusSuccess, time.Now()).
        Count(&upcomingAppointments).Error; err != nil {
        return nil, err
    }

    summary := map[string]int{
        "total_appointments":     int(totalAppointments),
        "completed_appointments": int(completedAppointments),
        "upcoming_appointments":  int(upcomingAppointments),
    }

    return summary, nil
}

type AppointmentHistoryResponse struct {
	AppointmentID string `json:"appointment_id"`
	Conclusion    string `json:"conclusion,omitempty"`
	Date          string `json:"date"`
	Medications   []struct {
		Name    string `json:"name"`
		Dosage  string `json:"dosage,omitempty"`
		Quantity string `json:"quantity,omitempty"`
	} `json:"medications,omitempty"`
}

func (s *AppointmentService) GetAppointmentHistoryByUserAndTherapist(userID, therapistID uuid.UUID) ([]AppointmentHistoryResponse, error) {
    var consultationHistories []models.ConsultationHistory

    err := s.DB.Preload("Appointment").
        Preload("Prescription.Medication").
        Joins("JOIN appointments ON appointments.id = consultation_histories.appointment_id").
        Where("appointments.user_id = ? AND appointments.therapist_id = ?", userID, therapistID).
        Find(&consultationHistories).Error

    if err != nil {
        return []AppointmentHistoryResponse{}, err 
    }

    response := make([]AppointmentHistoryResponse, 0) 

    for _, history := range consultationHistories {
        medications := make([]struct {
            Name     string `json:"name"`
            Dosage   string `json:"dosage,omitempty"`
            Quantity string `json:"quantity,omitempty"`
        }, 0) 

        for _, prescription := range history.Prescription {
            medications = append(medications, struct {
                Name     string `json:"name"`
                Dosage   string `json:"dosage,omitempty"`
                Quantity string `json:"quantity,omitempty"`
            }{
                Name:     prescription.Medication.Name,
                Dosage:   prescription.Dosage,
                Quantity: prescription.Quantity,
            })
        }

        response = append(response, AppointmentHistoryResponse{
            AppointmentID: history.AppointmentID.String(),
            Conclusion:    history.Conclusion,
            Date:          history.ConsultationDate.In(time.FixedZone("Asia/Jakarta", 7*60*60)).Format("2006-01-02"),
            Medications:   medications,
        })
    }

    return response, nil
}
