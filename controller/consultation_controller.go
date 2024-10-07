package controller

import (
	"net/http"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
)

type ConsultationHistoryController struct {
	ConsultationHistoryService *services.ConsultationHistoryService
}
func (ctrl *ConsultationHistoryController) GetAllUserConsultationHistory(c *gin.Context) {
	userID := c.Param("user_id")

	history, err := ctrl.ConsultationHistoryService.GetConsultationHistoryByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch consultation history"})
		return
	}

	// Format the response to include doctor (therapist) name
	var response []gin.H
	for _, h := range history {
		response = append(response, gin.H{
			"id":                h.ID,
			"appointment_id":    h.AppointmentID,
			"conclusion":        h.Conclusion,
			"consultation_date": h.ConsultationDate,
			"doctor": gin.H{
				"name": h.Appointment.Therapist.Name, // Get the therapist's name from the appointment
			},
			"prescriptions": h.Prescription, // This includes the medication and dosage info
		})
	}

	if len(response) == 0 {
		response = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{
		"consultations": response,
	})
}
