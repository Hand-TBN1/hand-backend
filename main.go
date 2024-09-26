package main

import (
	"log"
	"os"

	"github.com/Hand-TBN1/hand-backend/config"
	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/Hand-TBN1/hand-backend/routes"
	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    apiEnv := os.Getenv("ENV")
    if err != nil && apiEnv == "" {
        log.Println("fail to load env", err)
    }
    config.LoadEnv()

    // Pass all models to the migration function
    db := config.NewPostgresql(
        &models.User{},
        &models.Therapist{},
        &models.BookedSchedule{},
        &models.CheckIn{},
        &models.ChatMessage{},
        &models.ChatRoom{},
        &models.PositiveAffirmation{},
        &models.EmergencyHistory{},
        &models.MindfulnessExercise{},
        &models.PersonalHealthPlan{},
        &models.Appointment{},
        &models.ConsultationHistory{},
        &models.Medication{},
        &models.Prescription{},
    )

    redisClient := config.NewRedis()

    if redisClient != nil {
        log.Println("Connect Redis Successful")
    }

    if db != nil {
        log.Println("Connect Successful")
    } else {
        log.Println("Failed Connect")
    }

    engine := config.NewGin()

    routes.SetupRoutes(engine, db)

    log.Printf("Running on port %s", config.Env.ApiPort) 
    if err := engine.Run(":" + config.Env.ApiPort); err != nil {
        log.Fatalf("Failed to start server: %v\n", err)
    }
}
