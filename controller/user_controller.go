package controller

import (
	"net/http"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/Hand-TBN1/hand-backend/utilities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	UserService *services.UserService
}

// GetProfile handles the retrieval of the user's profile
func (ctrl *UserController) GetProfile(c *gin.Context) {
	// Extract the claims from middleware
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, apierror.NewApiErrorBuilder().
			WithStatus(http.StatusUnauthorized).
			WithMessage("Unauthorized access").
			Build())
		return
	}
	userClaims := claims.(*utilities.Claims)

	// Fetch the profile from the service
	user, apiErr := ctrl.UserService.GetProfile(userClaims.UserID)
	if apiErr != nil {
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	// Create a response with common user fields
	response := gin.H{
		"id":                 user.ID,
		"name":               user.Name,
		"email":              user.Email,
		"phone_number":       user.PhoneNumber,
		"image_url":          user.ImageURL,
		"role":               user.Role,
		"is_mobile_verified": user.IsMobileVerified,
	}

	// If the user is a therapist, include additional therapist-specific fields
	if user.Role == models.RoleTherapist && user.Therapist.ID != uuid.Nil {
		response["therapist"] = gin.H{
			"location":         user.Therapist.Location,
			"specialization":   user.Therapist.Specialization,
			"consultation":     user.Therapist.Consultation,
			"appointment_rate": user.Therapist.AppointmentRate,
		}
	} else {
		response["therapist"] = nil
	}

	c.JSON(http.StatusOK, gin.H{"profile": response})
}



// EditProfile handles the editing of the user's profile 
func (ctrl *UserController) EditProfile(c *gin.Context) {
	// Extract the claims from middleware
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, apierror.NewApiErrorBuilder().
			WithStatus(http.StatusUnauthorized).
			WithMessage("Unauthorized access").
			Build())
		return
	}
	userClaims := claims.(*utilities.Claims)

	var updatedUserRequest struct {
		Name     string `json:"name"`
		ImageURL string `json:"image_url"`
	}

	if err := c.ShouldBindJSON(&updatedUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, apierror.NewApiErrorBuilder().
			WithStatus(http.StatusBadRequest).
			WithMessage(apierror.ErrInvalidInput).
			Build())
		return
	}

	// Call the service to update only name and image URL
	apiErr := ctrl.UserService.EditProfile(userClaims.UserID, updatedUserRequest.Name, updatedUserRequest.ImageURL)
	if apiErr != nil {
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
