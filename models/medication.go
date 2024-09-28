package models

import (
	"time"

	"github.com/google/uuid"
)

type Medication struct {
    ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
    ImageURL         string   `json:"image_url"`
    Stock            int      `json:"stock"`
    Name             string `json:"name"`
    Price            int64 `json:"price"`
    Description      string `json:"description"`
    RequiresPrescription bool `json:"requiresPrescription"`
    CreatedAt        time.Time `json:"createdAt"`
    UpdatedAt        time.Time `json:"updatedAt"`

    // Associations
    Prescriptions    []Prescription    `gorm:"foreignKey:MedicationID"`
}
