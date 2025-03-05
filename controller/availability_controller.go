package controller

import (
	"net/http"
	"time"

	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
)

type AvailabilityController struct {
	AvailabilityService *services.AvailabilityService
}


func (ctrl *AvailabilityController) GetBlockedAvailability(c *gin.Context) {
	therapistID := c.Param("id")

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid timezone"})
		return
	}
	currentTime := time.Now().In(location)

	blockedDates, err := ctrl.AvailabilityService.GetBlockedAvailabilityByTherapistID(therapistID, currentTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blocked availability"})
		return
	}


	response := make([]gin.H, 0)
	for _, date := range blockedDates {
		response = append(response, gin.H{
			"date": date.Date.Format("2006-01-02"),
		})
	}

	c.JSON(http.StatusOK, gin.H{"blocked_dates": response})
}
