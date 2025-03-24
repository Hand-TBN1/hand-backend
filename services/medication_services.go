package services

import (
	"net/http"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/models"
	"gorm.io/gorm"
)

type MedicationService struct {
	DB *gorm.DB
}

func (service *MedicationService) GetMedications(name string) ([]models.Medication, *apierror.ApiError) {
	var medications []models.Medication

	query := service.DB
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	if err := query.Find(&medications).Error; err != nil {
		return nil, apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Database error").
			Build()
	}

	return medications, nil
}

func (service *MedicationService) GetMedicationByID(id string) (*models.Medication, *apierror.ApiError) {
	var medication models.Medication
	if err := service.DB.Where("id = ?", id).First(&medication).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierror.NewApiErrorBuilder().
				WithStatus(http.StatusNotFound).
				WithMessage("Medication not found").
				Build()
		}
		return nil, apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Database error").
			Build()
	}
	return &medication, nil
}

type UpdateMedicationDTO struct {
    Name             string
    Stock            int
    Price            int64
    Description      string
    RequiresPrescription bool
}

// Function to update medication by ID
func (service *MedicationService) UpdateMedicationByID(id string, updateDTO *UpdateMedicationDTO) *apierror.ApiError {
    var medication models.Medication

    // Fetch medication by ID
    if err := service.DB.Where("id = ?", id).First(&medication).Error; err != nil {
        return apierror.NewApiErrorBuilder().
            WithStatus(404).
            WithMessage("Medication not found").
            Build()
    }

    // Update medication fields
    medication.Name = updateDTO.Name
    medication.Stock = updateDTO.Stock
    medication.Price = updateDTO.Price
    medication.Description = updateDTO.Description
    medication.RequiresPrescription = updateDTO.RequiresPrescription

    // Save updated medication back to the DB
    if err := service.DB.Save(&medication).Error; err != nil {
        return apierror.NewApiErrorBuilder().
            WithStatus(500).
            WithMessage("Failed to save updated medication").
            Build()
    }

    return nil
}

func (service *MedicationService) AddMedication(medication *models.Medication) *apierror.ApiError {
	if err := service.DB.Create(medication).Error; err != nil {
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to add medication").
			Build()
	}
	return nil
}

func (service *MedicationService) UpdateMedication(id string, updatedMedication *models.Medication) *apierror.ApiError {
	var medication models.Medication
	if err := service.DB.Where("id = ?", id).First(&medication).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return apierror.NewApiErrorBuilder().
				WithStatus(http.StatusNotFound).
				WithMessage("Medication not found").
				Build()
		}
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Database error").
			Build()
	}

	// Update medication fields
	medication.Name = updatedMedication.Name
	medication.Price = updatedMedication.Price
	medication.Stock = updatedMedication.Stock
	medication.Description = updatedMedication.Description
	medication.RequiresPrescription = updatedMedication.RequiresPrescription
	medication.ImageURL = updatedMedication.ImageURL
	medication.UpdatedAt = updatedMedication.UpdatedAt

	if err := service.DB.Save(&medication).Error; err != nil {
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to update medication").
			Build()
	}

	return nil
}

// Delete a medication
func (service *MedicationService) DeleteMedication(id string) *apierror.ApiError {
	if err := service.DB.Where("id = ?", id).Delete(&models.Medication{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return apierror.NewApiErrorBuilder().
				WithStatus(http.StatusNotFound).
				WithMessage("Medication not found").
				Build()
		}
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to delete medication").
			Build()
	}

	return nil
}
