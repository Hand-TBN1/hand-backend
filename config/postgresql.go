package config

import (
	"fmt"
	"log"

	"gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresql(migrations ...any) *gorm.DB {
	gormLogger := logger.Default
	if Env.ENV != "production" {
		gormLogger = gormLogger.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
			Env.PostgresHost,
			Env.PostgresUser,
			Env.PostgresPassword,
			Env.PostgresDbName,
			Env.PostgresPort,
		),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := migratePostgresqlTables(db, migrations...); err != nil {
		log.Fatalln(err)
	}

	return db
}

func migratePostgresqlTables(db *gorm.DB, migrations ...any) error {

	if err := db.AutoMigrate(
		migrations..., // BREAKING: entities should be passed from cmd/api/main.go due to circular dependency issue
	); err != nil {
		return err
	}

	return nil
}
