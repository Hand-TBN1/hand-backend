package services

import (
	"log"
	"net/http"
	"time"

	"github.com/Hand-TBN1/hand-backend/apierror"
	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JournalService struct {
	DB *gorm.DB
}

func (service *JournalService) GetUserJournals(userID uuid.UUID, date *time.Time) ([]models.Journal, *apierror.ApiError) {
    var journals []models.Journal
    query := service.DB.Where("user_id = ?", userID)

    if date != nil {
        startOfDay := date.Truncate(24 * time.Hour)
        endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)

        log.Printf("Start of Day UTC: %v, End of Day UTC: %v", startOfDay, endOfDay)

        query = query.Where("created_at >= ? AND created_at <= ?", startOfDay, endOfDay)
    }

    if err := query.Find(&journals).Error; err != nil {
        return nil, apierror.NewApiErrorBuilder().
            WithStatus(http.StatusInternalServerError).
            WithMessage("Failed to retrieve journals").
            Build()
    }

    return journals, nil
}

func (service *JournalService) CreateJournal(journal *models.Journal) *apierror.ApiError {
	if err := service.DB.Create(journal).Error; err != nil {
		return apierror.NewApiErrorBuilder().
			WithStatus(http.StatusInternalServerError).
			WithMessage("Failed to create journal").
			Build()
	}
	return nil
}
