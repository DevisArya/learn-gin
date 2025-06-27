package main

import (
	"fmt"
	"os"

	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/config"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/helper"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	errLoad := godotenv.Load(".env")
	helper.PanicIfError(errLoad)

	validate := validator.New()
	db := config.NewDB()
	log := config.NewLogger()

	go func() {
		bootstrapResult, err := config.BootstrapGrpc(&config.BootstrapConfigGrpc{
			DB:       db,
			Validate: validate,
			Log:      log,
		})
		if err != nil {
			log.Fatalf("failed to bootstrap gRPC: %v", err)
		}
		fmt.Println("Notification gRPC server running on port 5000")
		if err := bootstrapResult.GRPCServer.Serve(bootstrapResult.Listener); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	app := gin.New()
	config.InitialMigration(db)
	// app.Use(logger.New(logger.Config{
	// 	Format: "${time} ${status} - ${method} ${path} ${latency} \n",
	// }))
	config.BootstrapRest(&config.BootstrapConfigRest{
		DB:       db,
		App:      app,
		Validate: validate,
		Log:      log,
	})

	port := os.Getenv("NOTIF_PORT_REST")
	if port == "" {
		port = ":8080"
	}

	err := app.Run(port)
	if err != nil {
		log.Fatalf("Failed to start REST server: %v", err)
	}
}
