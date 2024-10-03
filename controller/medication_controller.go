package controller

import (
	"net/http"
	"strconv"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/gin-gonic/gin"
)

type MedicationController struct {
	MedicationService *services.MedicationService
}

func (ctrl *MedicationController) GetMedications(c *gin.Context) {
	name := c.Query("name")

	medications, apiErr := ctrl.MedicationService.GetMedications(name)
	if apiErr != nil {
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	c.JSON(http.StatusOK, medications)
}

func (ctrl *MedicationController) AddMedication(c *gin.Context) {
	var medication models.Medication
	
	stock, err := strconv.Atoi(c.PostForm("stock"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock value"})
		return
	}

	price, err := strconv.ParseFloat(c.PostForm("price"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price value"})
		return
	}

	medication.Name = c.PostForm("name")
	medication.Stock = stock
	medication.Price = int64(price)
	medication.Description = c.PostForm("description")
	medication.RequiresPrescription, _ = strconv.ParseBool(c.PostForm("requiresPrescription"))
	medication.ImageURL = c.PostForm("image")

	apiErr := ctrl.MedicationService.AddMedication(&medication)
	if apiErr != nil {
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Medication added successfully"})
}

func (ctrl *MedicationController) UpdateMedication(c *gin.Context) {
	id := c.Param("id")
	var updatedMedication models.Medication

	if err := c.ShouldBindJSON(&updatedMedication); err != nil {
		apiErr := apierror.NewApiErrorBuilder().
			WithStatus(http.StatusBadRequest).
			WithMessage("Invalid input").
			Build()
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	apiErr := ctrl.MedicationService.UpdateMedication(id, &updatedMedication)
	if apiErr != nil {
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medication updated successfully"})
}

func (ctrl *MedicationController) DeleteMedication(c *gin.Context) {
	id := c.Param("id")

	apiErr := ctrl.MedicationService.DeleteMedication(id)
	if apiErr != nil {
		c.JSON(apiErr.HttpStatus, apiErr)
		return
	}

	c.Status(http.StatusNoContent)
}
