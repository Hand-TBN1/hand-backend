package models

import (
	"time"

	"github.com/google/uuid"
)

type Medication struct {
    ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    ImageURL         string   `json:"image_url"`
    Stock            int
    Name             string
    Price            int64
    Description      string
    RequiresPrescription bool
    CreatedAt        time.Time
    UpdatedAt        time.Time

    // Associations
    Prescriptions    []Prescription    `gorm:"foreignKey:MedicationID"`
}
