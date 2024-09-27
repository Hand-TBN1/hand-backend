package services

import (
	"errors"
	"net/http"
	"time"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/models"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

// GetProfile retrieves the user's profile from the database
func (service *UserService) GetProfile(userID string) (*models.User, *apierror.ApiError) {
	var user models.User
	err := service.DB.Preload("Therapist").Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.NewApiErrorBuilder().
				WithStatus(http.StatusNotFound).
				WithMessage(apierror.ErrUserNotFound).
				Build()
		}
		return nil, apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage(apierror.ErrInternalServerError).
			Build()
	}

	return &user, nil
}


// EditProfile updates the user's profile in the database
func (service *UserService) EditProfile(userID string, name, phoneNumber, imageURL string) *apierror.ApiError {
	// Find the user by ID
	var user models.User
	if err := service.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierror.NewApiErrorBuilder().
				WithStatus(http.StatusNotFound).
				WithMessage(apierror.ErrUserNotFound).
				Build()
		}
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage(err.Error()).
			Build()
	}

	// Update fields
	user.Name = name
	user.PhoneNumber = phoneNumber
	user.ImageURL = imageURL
	user.UpdatedAt = time.Now()

	// Save the changes
	if err := service.DB.Save(&user).Error; err != nil {
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage(err.Error()).
			Build()
	}

	return nil
}
