package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
    Admin        Role = "admin"
    Patient      Role = "patient"
    RoleTherapist Role = "therapist"
)

type User struct {
    ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    Name             string    `json:"name"`
    Email            string    `json:"email" gorm:"unique;not null"`
    PhoneNumber      string    `json:"phone_number" gorm:"unique;not null"`
    ImageURL         string    `json:"image_url"`
    Password         string    `json:"password"`
    Role             Role      `json:"role" gorm:"type:role_enum"`
    IsMobileVerified bool      `json:"is_mobile_verified"`
    CreatedAt        time.Time
    UpdatedAt        time.Time

    // Associations
    Therapist         *Therapist         `gorm:"foreignKey:UserID"`
    Appointment    []Appointment   `gorm:"foreignKey:UserID"`
    PositiveAffirmations []PositiveAffirmation `gorm:"foreignKey:PatientID"`
    PersonalHealthPlans  []PersonalHealthPlan  `gorm:"foreignKey:PatientID"`
    CheckIns            []CheckIn          `gorm:"foreignKey:UserID"`
    ChatMessages       []ChatMessage      `gorm:"foreignKey:SenderID"`
    EmergencyHistories []EmergencyHistory `gorm:"foreignKey:UserID"`
}