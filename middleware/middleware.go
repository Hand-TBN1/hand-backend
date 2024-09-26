package middleware

import (
	"net/http"
	"strings"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/utilities"
	"github.com/gin-gonic/gin"
)

// RoleMiddleware checks if the user has one of the required roles
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, apierror.NewApiErrorBuilder().
				WithStatus(http.StatusUnauthorized).
				WithMessage("Authorization header missing").
				Build())
			c.Abort()
			return
		}

		// Split "Bearer token"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, apierror.NewApiErrorBuilder().
				WithStatus(http.StatusUnauthorized).
				WithMessage("Invalid token format").
				Build())
			c.Abort()
			return
		}

		// Validate the token
		claims, err := utilities.ValidateJWT(tokenParts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, apierror.NewApiErrorBuilder().
				WithStatus(http.StatusUnauthorized).
				WithMessage("Invalid token").
				Build())
			c.Abort()
			return
		}

		// Check if the user's role is allowed
		userRole := claims.Role
		for _, role := range allowedRoles {
			if role == userRole {
				c.Next()
				return
			}
		}

		// If the role doesn't match any allowed roles, return forbidden
		c.JSON(http.StatusForbidden, apierror.NewApiErrorBuilder().
			WithStatus(http.StatusForbidden).
			WithMessage("Forbidden: You don't have access to this resource").
			Build())
		c.Abort()
	}
}


func Authenticate(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, apierror.NewApiErrorBuilder().
				WithStatus(http.StatusUnauthorized).
				WithMessage("Authorization header missing").
				Build())
			c.Abort()
			return
		}

		// Split "Bearer token"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, apierror.NewApiErrorBuilder().
				WithStatus(http.StatusUnauthorized).
				WithMessage("Invalid token format").
				Build())
			c.Abort()
			return
		}

		// Validate the token
		claims, err := utilities.ValidateJWT(tokenParts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, apierror.NewApiErrorBuilder().
				WithStatus(http.StatusUnauthorized).
				WithMessage("Invalid token").
				Build())
			c.Abort()
			return
		}

		// Check if the user's role is allowed
		c.Set("user.id", claims.UserID)
		c.Set("user.name", claims.Name)
		c.Set("user.role", claims.Role)
		c.Next()
	}
}


// Contoh Penggunaan
// router.Use(middleware.RoleMiddleware("admin", "therapist")) // Allow only admin and therapist
// 	{
// 		// Protected routes that only admin and therapist can access
// 		protected.GET("/some-endpoint", someHandler)
// 	}