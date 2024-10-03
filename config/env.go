package config

import (
	"log"
	"os"
	"strconv"

	"github.com/midtrans/midtrans-go"
)

type environmentVariables struct {
	ENV              string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDbName   string

	RedisHost        string
	RedisPort        string
	RedisPassword    string
	RedisDatabase    int

	ApiPort          string
	MidtransClientKey string
	MidtransServerKey string
	MidtransEnvironment midtrans.EnvironmentType
	

	TwilioAccountSID        string 
	TwilioAuthToken       string 
	TwilioVerifyServiceSID string

	FonnteAPIKey string
}

var Env *environmentVariables
func LoadEnv() {
	env := &environmentVariables{}
	var err error
	env.ENV = os.Getenv("ENV")
	if env.ENV == "" {
		log.Fatal("ENV is not set")
	}

	env.ApiPort = os.Getenv("API_PORT")

	env.PostgresHost = os.Getenv("POSTGRES_HOST")
	env.PostgresPort = os.Getenv("POSTGRES_PORT")
	env.PostgresUser = os.Getenv("POSTGRES_USER")
	env.PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	env.PostgresDbName = os.Getenv("POSTGRES_DB")

	env.RedisHost = os.Getenv("REDIS_HOST")
	env.RedisPort = os.Getenv("REDIS_PORT")
	env.RedisPassword = os.Getenv("REDIS_PASSWORD")
	env.RedisDatabase, err = strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil && env.ENV != "test" {
		log.Fatal("Fail to parse REDIS_DATABASE")
	}

	env.MidtransClientKey = os.Getenv("MIDTRANS_CLIENT_KEY")
	env.MidtransServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	env.MidtransEnvironment = midtrans.Sandbox
	if env.MidtransClientKey == "" || env.MidtransServerKey == "" {
		log.Fatal("Midtrans keys are not set")
	}

	env.TwilioAccountSID = os.Getenv("TWILIO_ACCOUNT_SID")
	env.TwilioAuthToken = os.Getenv("TWILIO_AUTH_TOKEN")
	env.TwilioVerifyServiceSID = os.Getenv("TWILIO_VERIFY_SERVICE_SID")
	if env.TwilioAccountSID == "" || env.TwilioAuthToken == "" || env.TwilioVerifyServiceSID == "" {
		log.Fatal("Twilio credentials are not set")
	}
	
	env.FonnteAPIKey = os.Getenv("FONNTE_API_KEY")
    if env.FonnteAPIKey == "" {
        log.Fatal("Fonnte API Key is not set")
    }

	Env = env
}
