package services

import (
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

	query := service.DB.Model(&models.Therapist{})

	if consultationType != "" {
		query = query.Where("consultation = ?", consultationType)
	}
	if location != "" {
		query = query.Where("location = ?", location)
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
