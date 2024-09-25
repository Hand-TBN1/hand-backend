package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
    Admin    Role = "admin"
    Patient  Role = "patient"
    RoleTherapist Role = "therapist"
)

type User struct {
    ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    Name             string
    Email            string
    PhoneNumber      string
    ImageURL         string
    Password         string
    Role             Role `gorm:"type:enum('admin', 'patient', 'therapist')"`
    IsMobileVerified bool
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
