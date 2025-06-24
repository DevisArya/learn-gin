package main

import (
	"os"

	"github.com/DevisArya/BE-challenge-syn/UserService/client/notification"
	"github.com/DevisArya/BE-challenge-syn/UserService/internal/config"
	"github.com/DevisArya/BE-challenge-syn/UserService/internal/helper"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {

	errLoad := godotenv.Load("../.env")
	helper.PanicIfError(errLoad)

	db := config.NewDB()
	config.InitialMigration(db)

	log := config.NewLogger()
	validate := validator.New()
	app := gin.New()
	app.Use(gin.Recovery())
	// app.Use(logger.New(logger.Config{
	// 	Format: "${time} ${status} - ${method} ${path} ${latency} \n",
	// }))

	notifClient, err := notification.NewNotificationClientGRPC(os.Getenv("ADDR_NOTIF"))
	if err != nil {
		log.Fatal("Failed to connect to Notification service:", err)
	}
	defer notifClient.Close()

	config.Bootstrap(&config.BootstrapConfig{
		DB:                 db,
		App:                app,
		Validate:           validate,
		NotificationClient: notifClient,
		Log:                log,
	})

	port := os.Getenv("USER_PORT_REST")
	if port == "" {
		port = ":8080"
	}

	if err := app.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
