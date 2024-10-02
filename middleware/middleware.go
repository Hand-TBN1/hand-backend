package middleware

import (
    "log"
	"net/http"
	"strings"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/utilities"
	"github.com/gin-gonic/gin"
)

// RoleMiddleware checks if the user is authenticated and optionally verifies their role.
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, apierror.NewApiErrorBuilder().
				WithStatus(http.StatusUnauthorized).
				WithMessage("Authorization header missing").
				Build())
			c.Abort()
			return
		}
        log.Println(authHeader);

		// Check if the token is a Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, apierror.NewApiErrorBuilder().
				WithStatus(http.StatusUnauthorized).
				WithMessage("Invalid Authorization header format").
				Build())
			c.Abort()
			return
		}

		// Extract the token part
		token := tokenParts[1]

		// Validate the token
		claims, err := utilities.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, apierror.NewApiErrorBuilder().
				WithStatus(http.StatusUnauthorized).
				WithMessage("Invalid token").
				Build())
			c.Abort()
			return
		}

        c.Set("claims", claims)
		// Access the role from the claims struct
		userRole := claims.Role

		// If no roles are specified, just check if the user is authenticated
		if len(allowedRoles) == 0 {
			c.Next()
			return
		}

		// Check if the user's role matches any of the allowed roles
		for _, role := range allowedRoles {
			if role == userRole {
				c.Next()
				return
			}
		}

		// If the user's role is not allowed, return a "Forbidden" response
		c.JSON(http.StatusForbidden, apierror.NewApiErrorBuilder().
			WithStatus(http.StatusForbidden).
			WithMessage("Forbidden: You don't have access to this resource").
			Build())
		c.Abort()
	}
}
