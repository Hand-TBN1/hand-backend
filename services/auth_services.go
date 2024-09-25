package services

import (
	"net/http"
	"time"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/Hand-TBN1/hand-backend/utilities"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

// Register a new user
func (service *AuthService) Register(user *models.User) *apierror.ApiError {
	// Check if the user already exists
	var existingUser models.User
	if err := service.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusConflict).
			WithMessage(apierror.ErrUserAlreadyExists).
			Build()
	}

	// Hash the password before saving
	hashedPassword, err := utilities.HashPassword(user.Password)
	if err != nil {
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to hash password").
			Build()
	}
	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if result := service.DB.Create(user); result.Error != nil {
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage(result.Error.Error()).
			Build()
	}
	return nil
}

// Login function to validate user credentials
func (service *AuthService) Login(email, password string) (*models.User, string, *apierror.ApiError) {
	var user models.User
	if err := service.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", apierror.NewApiErrorBuilder().
				WithStatus(http.StatusNotFound).
				WithMessage(apierror.ErrUserNotFound).
				Build()
		}
		return nil, "", apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Database error").
			Build()
	}

	// Check the password
	if !utilities.CheckPasswordHash(password, user.Password) {
		return nil, "", apierror.NewApiErrorBuilder().
			WithStatus(http.StatusUnauthorized).
			WithMessage(apierror.ErrInvalidCredentials).
			Build()
	}

	// Generate JWT token
	token, err := utilities.GenerateJWT(user.ID.String(), string(user.Role), user.Name)
	if err != nil {
		return nil, "", apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to generate token").
			Build()
	}

	return &user, token, nil
}
