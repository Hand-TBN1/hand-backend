package main

import (
	"log"
	"os"

	"github.com/Hand-TBN1/hand-backend/config"
	"github.com/Hand-TBN1/hand-backend/middleware"
	"github.com/Hand-TBN1/hand-backend/models"
	"github.com/Hand-TBN1/hand-backend/routes"
	"github.com/Hand-TBN1/hand-backend/services"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
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
        &models.CheckIn{},
        &models.ChatMessage{},
        &models.ChatRoom{},
        &models.PositiveAffirmation{},
        &models.EmergencyHistory{},
        &models.Media{},
        &models.Journal{},
        &models.Availability{},
        &models.PersonalHealthPlan{},
        &models.Appointment{},
        &models.ConsultationHistory{},
        &models.Medication{},
        &models.Prescription{},
        &models.MedicationHistoryTransaction{},
        &models.MedicationHistoryItem{},
    )

    redisClient := config.NewRedis()
    config.LoadR2Config()

    if redisClient != nil {
        log.Println("Connect Redis Successful")
    }

    if db != nil {
        log.Println("Connect Successful")
    } else {
        log.Println("Failed Connect")
    }

    config.SetupMidtrans()
    paymentService := &services.PaymentService{}
    checkInService := &services.CheckInService{DB: db}

    engine := config.NewGin()
    engine.Use(middleware.CORS())

    routes.SetupAuthRoutes(engine, db)
    routes.RegisterCheckInRoutes(engine, db)
    routes.RegisterMedicationRoutes(engine, db)
    routes.RegisterMediaRoutes(engine, db)
    routes.RegisterMedicationTransactionHistoryRoutes(engine, db, paymentService)
    routes.RegisterTherapistRoutes(engine, db)
    routes.SetupPaymentRoutes(engine, db)  
    routes.RegisterUserRoutes(engine, db)  
    routes.RegisterAppointmentRoutes(engine, db,paymentService)  
    routes.RegisterJournalRoutes(engine, db)
    routes.RegisterPrescriptionRoutes(engine, db)
    routes.RegisterChatRoutes(engine,db)
    routes.RegisterCloudflareRoutes(engine)

    c := cron.New()
    _, err = c.AddFunc("30 14 * * *", func() { 
        users, err := checkInService.CheckUserCheckIns()
        if err != nil {
            log.Println("Error fetching users:", err)
            return
        }

        for _, user := range users {
            notificationErr := checkInService.SendReminder(user.PhoneNumber)
            if notificationErr != nil {
                log.Println("Error sending reminder:", notificationErr)
            }
        }
    })
    if err != nil {
        log.Fatalf("Error scheduling the task: %v", err)
    }
    c.Start()
    

	go func() {
		select {}
	}()

    log.Printf("Running on port %s", config.Env.ApiPort) 
    if err := engine.Run(config.Env.ApiPort); err != nil {
        log.Fatalf("Failed to start server: %v\n", err)
    }
}
