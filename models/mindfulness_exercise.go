package models

import (
	"time"

	"github.com/google/uuid"
)

type MindfulnessExercise struct {
    ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    ExerciseName     string
    ExerciseDescription string
    ExerciseType     string
    ImageURL         string
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
