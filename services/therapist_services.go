package services

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TherapistService struct {
	DB *gorm.DB
}

func (service *TherapistService) GetTherapistsFiltered(consultationType, location string, date time.Time) ([]models.Therapist, *apierror.ApiError) {
	var therapists []models.Therapist

	query := service.DB.Model(&models.Therapist{}).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, image_url")
	})

	if consultationType != "" {
		query = query.Where("consultation = ?", consultationType)
	}
	if location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}

	if err := query.Find(&therapists).Error; err != nil {
		return nil, apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to retrieve therapists").
			Build()
	}

	if date.IsZero() {
		return therapists, nil
	}

	var availableTherapists []models.Therapist
	for _, therapist := range therapists {
		if isTherapistAvailableForDate(service.DB, therapist.UserID, date) && hasAvailableTimeSlots(service.DB, therapist.UserID, date) {
			availableTherapists = append(availableTherapists, therapist)
		}
	}

	return availableTherapists, nil
}


func isTherapistAvailableForDate(db *gorm.DB, therapistID uuid.UUID, date time.Time) bool {
	var availability models.Availability

	err := db.Where("therapist_id = ? AND date = ?", therapistID, date).First(&availability).Error

	if err == gorm.ErrRecordNotFound {
		fmt.Printf("No availability record found for therapist %s on %s. Defaulting to available.\n", therapistID, date.Format("2006-01-02"))
		return true
	} else if err != nil {
		fmt.Printf("Error querying availability for therapist %s on %s: %v\n", therapistID, date.Format("2006-01-02"), err)
		return false
	}

	fmt.Printf("Availability record found for therapist %s on %s: IsAvailable = %v\n", therapistID, date.Format("2006-01-02"), availability.IsAvailable)

	return availability.IsAvailable
}



func hasAvailableTimeSlots(db *gorm.DB, therapistID uuid.UUID, date time.Time) bool {
	// Define time slots (8 AM - 12 PM, 1 PM - 5 PM)
	timeSlots := []time.Time{}
	for hour := 8; hour < 12; hour++ {
		timeSlots = append(timeSlots, time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, time.UTC))
	}
	for hour := 13; hour < 17; hour++ {
		timeSlots = append(timeSlots, time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, time.UTC))
	}

	for _, timeSlot := range timeSlots {
		var count int64
		db.Model(&models.Appointment{}).
			Where("therapist_id = ?", therapistID).
			Where("appointment_date = ?", timeSlot).
			Where("status = ?", models.Success).
			Count(&count)
		
		if count == 0 {
			return true 
		}
	}

	return false
}

func (service *TherapistService) UpdateAvailabilityByDate(therapistID string, date time.Time, isAvailable bool) *apierror.ApiError {
	var availability models.Availability
	err := service.DB.Where("therapist_id = ? AND date = ?", therapistID, date).First(&availability).Error

	if err == gorm.ErrRecordNotFound {
		availability = models.Availability{
			ID:          uuid.New(),
			TherapistID: uuid.MustParse(therapistID),
			Date:        date,
			IsAvailable: isAvailable,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := service.DB.Create(&availability).Error; err != nil {
			return apierror.NewApiErrorBuilder().
				WithStatus(http.StatusInternalServerError).
				WithMessage("Failed to create availability").
				Build()
		}
	} else if err != nil {
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to check availability").
			Build()
	} else {
		availability.IsAvailable = isAvailable
		availability.UpdatedAt = time.Now()

		if err := service.DB.Save(&availability).Error; err != nil {
			return apierror.NewApiErrorBuilder().
				WithStatus(http.StatusInternalServerError).
				WithMessage("Failed to update availability").
				Build()
		}
	}

	return nil
}

func (service *TherapistService) AddTherapist(therapist *models.Therapist) *apierror.ApiError {
	if err := service.DB.Debug().Create(therapist).Error; err != nil { 
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage(err.Error()). 
			Build()
	}
	return nil
}

func (service *TherapistService) AddUser(user *models.User) *apierror.ApiError {
	if err := service.DB.Create(user).Error; err != nil {
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to create user").
			Build()
	}
	return nil
}


func (service *TherapistService) GetTherapistDetails(therapistID string) (*models.Therapist, error) {
	var therapist models.Therapist
	// Preload the associated User to retrieve details like role, name, etc.
	if err := service.DB.Preload("User").Where("user_id = ?", therapistID).First(&therapist).Error; err != nil {
		return nil, errors.New("therapist not found")
	}
	return &therapist, nil
}


// GetAvailableSchedules fetches available schedules for a therapist
func (service *TherapistService) GetAvailableSchedules(therapistID, date, consultationType string) ([]string, error) {
	var appointments []models.Appointment

	// Parse the provided date or default to today's date
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		parsedDate = time.Now() // Default to today if parsing fails
	}

	// Query successful appointments for the therapist on the given date and consultation type
	if err := service.DB.Where("therapist_id = ? AND type = ? AND status = ? AND DATE(appointment_date) = ?", therapistID, models.ConsultationType(consultationType), models.Success, parsedDate).Find(&appointments).Error; err != nil {
		return nil, errors.New("error fetching appointments")
	}

	// Define time slots (e.g., 8 AM to 11 AM, 1 PM to 5 PM)
	timeSlots := generateTimeSlots(parsedDate)

	// Loop through time slots and check if there's an appointment that blocks the slot
	availableSlots := filterAvailableSlots(timeSlots, appointments)

	return availableSlots, nil
}

// Helper function to generate time slots between 8 AM - 11 AM and 1 PM - 5 PM
func generateTimeSlots(date time.Time) []time.Time {
	var slots []time.Time

	// 8 AM to 11 AM slots
	for hour := 8; hour <= 10; hour++ {
		slots = append(slots, time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, time.Local))
	}

	// 1 PM to 5 PM slots
	for hour := 13; hour <= 16; hour++ {
		slots = append(slots, time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, time.Local))
	}

	return slots
}

// Helper function to filter available time slots based on existing appointments
func filterAvailableSlots(timeSlots []time.Time, appointments []models.Appointment) []string {
	var availableSlots []string

	for _, slot := range timeSlots {
		isAvailable := true

		for _, appointment := range appointments {
			appointmentStart := appointment.AppointmentDate

			// Compare the time slot with the appointment's start time
			if appointmentStart.Equal(slot) {
				isAvailable = false
				break
			}
		}

		if isAvailable {
			availableSlots = append(availableSlots, slot.Format("15:04")) // Format as "HH:MM"
		}
	}

	return availableSlots
}