package services

import (
	"errors"
	"time"

	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/Hand-TBN1/hand-backend/utilities"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

// Register a new user
func (service *AuthService) Register(user *models.User) error {
	// Hash the password before saving
	hashedPassword, err := utilities.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if result := service.DB.Create(user); result.Error != nil {
		return result.Error
	}
	return nil
}

// Login function to validate user credentials
func (service *AuthService) Login(email, password string) (*models.User, string, error) {
	var user models.User
	if err := service.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("user not found")
		}
		return nil, "", err
	}

	// Check the password
	if !utilities.CheckPasswordHash(password, user.Password) {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utilities.GenerateJWT(user.ID.String(), string(user.Role), user.Name)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}
