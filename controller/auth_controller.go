package controller

import (
	"net/http"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/Hand-TBN1/hand-backend/utilities"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService *services.AuthService
}

// Register a new user
func (ctrl *AuthController) Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		apiErr := apierror.NewApiErrorBuilder().
			WithStatus(http.StatusBadRequest).
			WithMessage(apierror.ErrInvalidInput).
			Build()
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	if apiErr := ctrl.AuthService.Register(&user); apiErr != nil {
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login a user and return JWT token in an HTTP-only cookie
	func (ctrl *AuthController) Login(c *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			apiErr := apierror.NewApiErrorBuilder().
				WithStatus(http.StatusBadRequest).
				WithMessage(apierror.ErrInvalidInput).
				Build()
			c.JSON(apiErr.HttpStatus, apiErr)
			return
		}

		user, token, apiErr := ctrl.AuthService.Login(req.Email, req.Password)
		if apiErr != nil {
			c.JSON(apiErr.HttpStatus, apiErr)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token" : token,
			"user" : user,
		})
	}

// Handler for sending OTP
func (ctrl *AuthController) SendOTP(c *gin.Context) {
    claims, exists := c.Get("claims")
    if !exists {
        apiErr := apierror.NewApiErrorBuilder().
            WithStatus(http.StatusUnauthorized).
            WithMessage("Unauthorized access").
            Build()
        c.JSON(apiErr.HttpStatus, apiErr)
        return
    }
    userClaims := claims.(*utilities.Claims)

    user, apiErr := ctrl.AuthService.GetUserByID(userClaims.UserID)
    if apiErr != nil {
        c.JSON(apiErr.HttpStatus, apiErr)
        return
    }

    if apiErr := ctrl.AuthService.SendOTP(user.PhoneNumber); apiErr != nil {
        c.JSON(apiErr.HttpStatus, apiErr)
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// Handler for verifying OTP
func (ctrl *AuthController) VerifyOTP(c *gin.Context) {
    var req struct {
        OTP string `json:"otp"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        apiErr := apierror.NewApiErrorBuilder().
            WithStatus(http.StatusBadRequest).
            WithMessage(apierror.ErrInvalidInput).
            Build()
        c.JSON(apiErr.HttpStatus, apiErr)
        return
    }

    claims, exists := c.Get("claims")
    if !exists {
        apiErr := apierror.NewApiErrorBuilder().
            WithStatus(http.StatusUnauthorized).
            WithMessage("Unauthorized access").
            Build()
        c.JSON(apiErr.HttpStatus, apiErr)
        return
    }
    userClaims := claims.(*utilities.Claims)

    user, apiErr := ctrl.AuthService.GetUserByID(userClaims.UserID)
    if apiErr != nil {
        c.JSON(apiErr.HttpStatus, apiErr)
        return
    }

    if apiErr := ctrl.AuthService.VerifyOTP(user.PhoneNumber, req.OTP); apiErr != nil {
        c.JSON(apiErr.HttpStatus, apiErr)
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}
