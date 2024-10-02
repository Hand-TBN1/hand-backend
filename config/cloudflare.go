package config

import (
	"log"
	"os"
)

type CloudflareR2Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string
	BucketName      string
}

var R2Config *CloudflareR2Config

func LoadR2Config() {
	R2Config = &CloudflareR2Config{
		AccessKeyID:     os.Getenv("CLOUDFLARE_R2_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("CLOUDFLARE_R2_SECRET_ACCESS_KEY"),
		Endpoint:        os.Getenv("CLOUDFLARE_R2_ENDPOINT"),
		BucketName:      os.Getenv("CLOUDFLARE_R2_BUCKET_NAME"),
	}

	// Check if all environment variables are set
	if R2Config.AccessKeyID == "" || R2Config.SecretAccessKey == "" || R2Config.Endpoint == "" || R2Config.BucketName == "" {
		log.Fatal("Cloudflare R2 environment variables not set")
	}
}
