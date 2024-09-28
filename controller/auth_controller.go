package controller

import (
	"net/http"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/Hand-TBN1/hand-backend/services"
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

		// Set JWT as HTTP-only cookie
		c.SetCookie(
			"auth_token",  
			token,         
			60*60*24*365*0.5, // Expiry time in seconds (0.5 year)
			"/",           // Cookie path
			"",            // Domain (empty means default, based on request domain)
			true,          // Secure (set to true for HTTPS)
			true,          // HTTPOnly (prevents access from JavaScript)
		)

		c.SetCookie(
			"user_id",  
			user.ID.String(),         
			60*60*24*365*0.5, // Expiry time in seconds (0.5 year)
			"/",           
			"",            
			false,         // Not HTTP-only 
			false,         
		)
		c.SetCookie(
			"user_name",  
			user.Name,         
			60*60*24*365*0.5, // Expiry time in seconds (0.5 year)
			"/",           
			"",            
			false,         // Not HTTP-only 
			false,         
		)
		
		c.SetCookie(	
			"user_role",  
			string(user.Role),         
			60*60*24*365*0.5, // Expiry time in seconds (0.5 year)
			"/",           
			"",            
			false,         // Not HTTP-only
			false,         
		)
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token" : token,
			"user" : user,
		})
	}

