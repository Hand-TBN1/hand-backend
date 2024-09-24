package main

import (
	"log"
	"os"

	"github.com/Hand-TBN1/hand-backend/config"
	"github.com/Hand-TBN1/hand-backend/schema"
	"github.com/joho/godotenv"
)

func main(){
	err := godotenv.Load()
	apiEnv := os.Getenv("ENV")
	if err != nil && apiEnv == "" {
		log.Println("fail to load env", err)
	}
	config.LoadEnv()

	db := config.NewPostgresql(
		&schema.User{},
	)

	if db != nil {
		log.Println("Connect Succesful")
	}else{
		log.Println("Failed Connect")
	}

	engine := config.NewGin()
	log.Printf("Running on port %s", config.Env.ApiPort) 
	if err := engine.Run(":" + config.Env.ApiPort); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}