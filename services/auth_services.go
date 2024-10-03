package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
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


func (service *AuthService) GetUserByID(userID string) (*models.User, *apierror.ApiError) {
    var user models.User
    if err := service.DB.First(&user, "id = ?", userID).Error; err != nil {
        return nil, apierror.NewApiErrorBuilder().
            WithStatus(http.StatusNotFound).
            WithMessage(apierror.ErrUserNotFound).
            Build()
    }
    return &user, nil
}

func (service *AuthService) SendOTP(phoneNumber string) *apierror.ApiError {
	apiKey := os.Getenv("FONNTE_API_KEY")
	url := "https://api.fonnte.com/send"
	otp := GenerateOTP()

	formattedPhoneNumber := phoneNumber
    if len(phoneNumber) > 0 && phoneNumber[0] == '0' {
		formattedPhoneNumber = "+62" + phoneNumber[1:] 
    }
	
	payload := map[string]string{
		"target":     formattedPhoneNumber,
		"message":    fmt.Sprintf("Your OTP is: %s", otp),
		"countryCode": "62",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("SendOTP Service: Failed to encode payload: %v", err)
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to prepare OTP payload").
			Build()
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("SendOTP Service: Failed to create request: %v", err)
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to create OTP request").
			Build()
	}

	req.Header.Set("Authorization", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Printf("SendOTP Service: Failed to send OTP: %v", err)
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to send OTP").
			Build()
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		log.Printf("SendOTP Service: Fonnte API Error: %s", string(body))
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage(fmt.Sprintf("Error sending OTP: %s", string(body))).
			Build()
	}

	// Save OTP to user
	if err := service.SaveOTPToUser(phoneNumber, otp); err != nil {
		log.Printf("SendOTP Service: Failed to save OTP: %v", err)
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to save OTP").
			Build()
	}

	return nil
}

func GenerateOTP() string {
	randGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", randGenerator.Intn(1000000))
}



func (service *AuthService) VerifyOTP(phone_number, inputOTP string) *apierror.ApiError {
    var user models.User
    if err := service.DB.First(&user, "phone_number = ?", phone_number).Error; err != nil {
        return apierror.NewApiErrorBuilder().
            WithStatus(http.StatusNotFound).
            WithMessage("User not found").
            Build()
    }

    if user.OTP != inputOTP {
        return apierror.NewApiErrorBuilder().
            WithStatus(http.StatusUnauthorized).
            WithMessage("Invalid OTP").
            Build()
    }

    if time.Now().After(user.OTPExpiresAt) {
        return apierror.NewApiErrorBuilder().
            WithStatus(http.StatusUnauthorized).
            WithMessage("OTP has expired").
            Build()
    }
	if err := service.DB.Model(&user).Update("is_mobile_verified", true).Error; err != nil {
        return apierror.NewApiErrorBuilder().
            WithStatus(http.StatusInternalServerError).
            WithMessage("Failed to update mobile verification status").
            Build()
    }

    log.Println("OTP verified successfully")
    return nil
}

func (service *AuthService) SaveOTPToUser(phoneNumber string, otp string) *apierror.ApiError {
    otpExpiry := time.Now().Add(5 * time.Minute) 

    if err := service.DB.Model(&models.User{}).Where("phone_number = ?", phoneNumber).
        Updates(map[string]interface{}{
            "otp":            otp,
            "otp_expires_at": otpExpiry,
        }).Error; err != nil {
        return apierror.NewApiErrorBuilder().
            WithStatus(http.StatusInternalServerError).
            WithMessage(fmt.Sprintf("Failed to save OTP: %v", err)).
            Build()
    }

    return nil
}
