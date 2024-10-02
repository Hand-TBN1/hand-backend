package services

import (
	"errors"
	"time"

	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CheckInService struct {
	DB *gorm.DB
}

func (service *CheckInService) CreateCheckIn(checkIn *models.CheckIn) error {
	if err := service.DB.Create(checkIn).Error; err != nil {
		return err
	}
	return nil
}

func (service *CheckInService) GetCheckIn(id string) (*models.CheckIn, error) {
	var checkIn models.CheckIn
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	if err := service.DB.First(&checkIn, "id = ?", parsedID).Error; err != nil {
		return nil, err
	}
	return &checkIn, nil
}

func (service *CheckInService) GetAllCheckIns() ([]models.CheckIn, error) {
	var checkIns []models.CheckIn
	if err := service.DB.Find(&checkIns).Error; err != nil {
		return nil, err
	}
	return checkIns, nil
}

func (service *CheckInService) UpdateCheckIn(id string, newCheckIn models.CheckIn) error {
	var checkIn models.CheckIn
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	if err := service.DB.First(&checkIn, "id = ?", parsedID).Error; err != nil {
		return errors.New("check-in not found")
	}
	checkIn.MoodScore = newCheckIn.MoodScore
	checkIn.Notes = newCheckIn.Notes
	checkIn.UpdatedAt = time.Now()

	if err := service.DB.Save(&checkIn).Error; err != nil {
		return err
	}
	return nil
}
	
func (service *CheckInService) FindCheckInByUserIDAndDate(userID uuid.UUID, date string) (*models.CheckIn, error) {
	var checkIn models.CheckIn

	// Query the database for a check-in by userID and the date in UTC
	err := service.DB.Where("user_id = ? AND DATE(check_in_date) = ?", userID, date).First(&checkIn).Error
	if err != nil {
		return nil, err
	}

	return &checkIn, nil
}

	
func (service *CheckInService) UpdateCheckInByUserIDAndDate(userID uuid.UUID, date string, updatedCheckIn models.CheckIn) error {
	return service.DB.Model(&models.CheckIn{}).
		Where("user_id = ? AND DATE(check_in_date) = ?", userID, date).
		Updates(map[string]interface{}{
			"mood_score": updatedCheckIn.MoodScore,
			"notes":      updatedCheckIn.Notes,
			"feelings" : updatedCheckIn.Feelings,
			"updated_at": updatedCheckIn.UpdatedAt,
		}).Error
}

func (service *CheckInService) FindCheckInByUserIDAndDateRange(userID uuid.UUID, start time.Time, end time.Time) (*models.CheckIn, error) {
    var checkIn models.CheckIn

    err := service.DB.Where("user_id = ? AND check_in_date BETWEEN ? AND ?", userID, start, end).First(&checkIn).Error
    if err != nil {
        return nil, err
    }

    return &checkIn, nil
}


func (service *CheckInService) CheckTodayCheckIn(userID uuid.UUID) (*models.CheckIn, error) {
    location, _ := time.LoadLocation("Asia/Jakarta") 

    localNow := time.Now().In(location)

    startOfDay := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, location).UTC()
    endOfDay := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 23, 59, 59, 999999999, location).UTC()

    return service.FindCheckInByUserIDAndDateRange(userID, startOfDay, endOfDay)
}


func (s *CheckInService) GetAllUserCheckIns(userID uuid.UUID) ([]models.CheckIn, error) {
    var checkIns []models.CheckIn

    err := s.DB.Where("user_id = ?", userID).Order("check_in_date desc").Find(&checkIns).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, gorm.ErrRecordNotFound
        }
        return nil, err
    }

    return checkIns, nil
}