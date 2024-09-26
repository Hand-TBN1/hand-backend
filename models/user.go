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
    Name             string
    Email            string `gorm:"unique;not null"`
    PhoneNumber      string `gorm:"unique;not null"`
    ImageURL         string
    Password         string
    Role             Role `gorm:"type:role_enum"`
    IsMobileVerified bool
    CreatedAt        time.Time
    UpdatedAt        time.Time

    // Associations
    Therapist         Therapist         `gorm:"foreignKey:UserID"`
    BookedSchedules    []BookedSchedule   `gorm:"foreignKey:UserID"`
    PositiveAffirmations []PositiveAffirmation `gorm:"foreignKey:PatientID"`
    PersonalHealthPlans  []PersonalHealthPlan  `gorm:"foreignKey:PatientID"`
    CheckIns            []CheckIn          `gorm:"foreignKey:UserID"`
    ChatMessages       []ChatMessage      `gorm:"foreignKey:SenderID"`
    EmergencyHistories []EmergencyHistory `gorm:"foreignKey:UserID"`
}
