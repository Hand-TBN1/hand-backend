package middleware

import (
	"log"
	"net/http"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/utilities"
	"github.com/gin-gonic/gin"
)

// RoleMiddleware checks if the user is authenticated and optionally verifies their role.
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract the JWT from the "auth_token" HTTP-only cookie
		log.Println(c)
        token, err := c.Request.Cookie("authToken")
		log.Println(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, apierror.NewApiErrorBuilder().
                WithStatus(http.StatusUnauthorized).
                WithMessage("Authorization token missing").
                Build())
            c.Abort()
            return
        }

        // Validate the token
        claims, err := utilities.ValidateJWT(token.Value)
        if err != nil {
            c.JSON(http.StatusUnauthorized, apierror.NewApiErrorBuilder().
                WithStatus(http.StatusUnauthorized).
                WithMessage("Invalid token").
                Build())
            c.Abort()
            return
        }

        // Store the claims in the context for future use
        c.Set("claims", claims)

        // Extract the user's role from the "user_role" cookie
        userRole, err := c.Request.Cookie("user_role")
        if err != nil {
            c.JSON(http.StatusForbidden, apierror.NewApiErrorBuilder().
                WithStatus(http.StatusForbidden).
                WithMessage("User role not found").
                Build())
            c.Abort()
            return
        }

        // If no roles are specified, just check if the user is authenticated
        if len(allowedRoles) == 0 {
            // Proceed since no role check is needed
            c.Next()
            return
        }

        // Check if the user's role matches any of the allowed roles
        for _, role := range allowedRoles {
            if role == userRole.Value {
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
