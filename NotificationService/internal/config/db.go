package config

import (
	"fmt"
	"os"
	"time"

	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/entity"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/helper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	DB_Username string
	DB_Password string
	DB_Port     string
	DB_Host     string
	DB_Name     string
	DB_Ssl_Mode string
	DB_TimeZone string
}

func NewDB() *gorm.DB {

	config := &DBConfig{
		DB_Username: os.Getenv("DB_USER"),
		DB_Password: os.Getenv("DB_PASSWORD"),
		DB_Port:     os.Getenv("DB_PORT"),
		DB_Host:     os.Getenv("NOTIF_DB_HOST"),
		DB_Name:     os.Getenv("NOTIF_DB_NAME"),
		DB_Ssl_Mode: os.Getenv("DB_SSL_MODE"),
		DB_TimeZone: os.Getenv("DB_TIMEZONE"),
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.DB_Host,
		config.DB_Username,
		config.DB_Password,
		config.DB_Name,
		config.DB_Port,
		config.DB_Ssl_Mode,
		config.DB_TimeZone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	helper.PanicIfError(err)

	sqlDB, err := db.DB()
	helper.PanicIfError(err)

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)
	sqlDB.SetConnMaxIdleTime(60 * time.Minute)

	return db
}

func InitialMigration(db *gorm.DB) {

	err := db.Exec(`DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'notification_action') THEN
			CREATE TYPE notification_action AS ENUM ('CREATE', 'UPDATE', 'DELETE');
		END IF;
	END
	$$;`).Error
	helper.PanicIfError(err)

	err = db.AutoMigrate(
		entity.Notification{},
	)

	helper.PanicIfError(err)
}
