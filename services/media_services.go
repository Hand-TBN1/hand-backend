package services

import (
	"net/http"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/models"
	"gorm.io/gorm"
)

type MediaService struct {
	DB *gorm.DB
}

func (service *MediaService) AddMedia(media *models.Media) *apierror.ApiError {
	if err := service.DB.Create(media).Error; err != nil {
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to create media").
			Build()
	}
	return nil
}

func (service *MediaService) GetAllMedia() ([]models.Media, *apierror.ApiError) {
	var mediaList []models.Media
	if err := service.DB.Find(&mediaList).Error; err != nil {
		return nil, apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to retrieve media").
			Build()
	}
	return mediaList, nil
}

func (service *MediaService) GetMedia(id string) (*models.Media, *apierror.ApiError) {
	var media models.Media
	if err := service.DB.Where("id = ?", id).First(&media).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierror.NewApiErrorBuilder().
				WithStatus(http.StatusNotFound).
				WithMessage("Media not found").
				Build()
		}
		return nil, apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to retrieve media").
			Build()
	}
	return &media, nil
}

// UpdateMedia updates an existing media entry
func (service *MediaService) UpdateMedia(media *models.Media) *apierror.ApiError {
	if err := service.DB.Save(media).Error; err != nil {
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to update media").
			Build()
	}
	return nil
}

// DeleteMedia deletes a media entry by ID
func (service *MediaService) DeleteMedia(id string) *apierror.ApiError {
	if err := service.DB.Where("id = ?", id).Delete(&models.Media{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return apierror.NewApiErrorBuilder().
				WithStatus(http.StatusNotFound).
				WithMessage("Media not found").
				Build()
		}
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to delete media").
			Build()
	}
	return nil
}
