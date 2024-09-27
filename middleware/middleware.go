package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/utilities"
	"github.com/gin-gonic/gin"
)

// RoleMiddleware checks if the user is authenticated and optionally verifies their role.
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

		// Store the claims in the context
		c.Set("claims", claims)

		// If allowedRoles is empty, just check that the user is authenticated
		if len(allowedRoles) == 0 {
			// No role check needed, proceed
			c.Next()
			return
		}

		// Check if the user's role is allowed
		for _, role := range allowedRoles {
			fmt.Printf("User role from token: %s, Allowed role: %s\n", claims.Role, role)
			if role == claims.Role {
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
