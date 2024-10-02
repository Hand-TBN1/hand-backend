package services

import (
	"net/http"
	"time"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/config"
	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/Hand-TBN1/hand-backend/utilities"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
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

	if user.Role != models.Patient {
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusBadRequest).
			WithMessage("Only users with role patient can register").
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

	user.IsMobileVerified = false
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

func (service *AuthService) SendOTP(phoneNumber string) *apierror.ApiError {
    client := twilio.NewRestClientWithParams(twilio.ClientParams{
        Username: config.Env.TwilioAccountSID,
        Password: config.Env.TwilioAuthToken,
    })

    params := &verify.CreateVerificationParams{}
    params.SetTo(phoneNumber)
    params.SetChannel("sms") 

    _, err := client.VerifyV2.CreateVerification(config.Env.TwilioVerifyServiceSID, params)
    if err != nil {
        return apierror.NewApiErrorBuilder().
            WithStatus(http.StatusInternalServerError).
            WithMessage("Failed to send OTP").
            Build()
    }

    return nil
}

func (service *AuthService) VerifyOTP(phoneNumber string, otp string) *apierror.ApiError {
    client := twilio.NewRestClientWithParams(twilio.ClientParams{
        Username: config.Env.TwilioAccountSID,
        Password: config.Env.TwilioAuthToken,
    })

    params := &verify.CreateVerificationCheckParams{}
    params.SetTo(phoneNumber)
    params.SetCode(otp)

    // Correct method to call CreateVerificationCheck
    resp, err := client.VerifyV2.CreateVerificationCheck(config.Env.TwilioVerifyServiceSID, params)
    if err != nil || *resp.Status != "approved" {
        return apierror.NewApiErrorBuilder().
            WithStatus(http.StatusUnauthorized).
            WithMessage("Invalid OTP").
            Build()
    }

	service.DB.Model(&models.User{}).Where("phone_number = ?", phoneNumber).Update("is_mobile_verified", true)

    return nil
}

