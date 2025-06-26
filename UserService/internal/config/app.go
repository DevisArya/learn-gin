package config

import (
	"github.com/DevisArya/BE-challenge-syn/UserService/client/notification"
	ctrl "github.com/DevisArya/BE-challenge-syn/UserService/internal/delivery/http"
	"github.com/DevisArya/BE-challenge-syn/UserService/internal/delivery/http/route"
	"github.com/DevisArya/BE-challenge-syn/UserService/internal/helper"
	"github.com/DevisArya/BE-challenge-syn/UserService/internal/repository"
	"github.com/DevisArya/BE-challenge-syn/UserService/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB                 *gorm.DB
	App                *gin.Engine
	Validate           *validator.Validate
	NotificationClient notification.NotificationClient
	Log                *logrus.Logger
}

func Bootstrap(cfg *BootstrapConfig) {
	//setup helper
	helperTx := helper.NewHelper(cfg.DB)

	// setup repositories
	userRepository := repository.NewUserRepository(cfg.Log, cfg.DB)

	// setup use cases
	userUseCase := usecase.NewUserUseCase(userRepository, cfg.Validate, cfg.NotificationClient, cfg.Log, helperTx)

	// setup controller
	userController := ctrl.NewUserController(userUseCase, cfg.Log)

	// // setup middleware
	// authMiddleware := middleware.NewAuth(userUseCase)

	routeConfig := route.RouteConfig{
		App:            cfg.App,
		UserController: userController,
		// AuthMiddleware: authMiddleware,
	}
	routeConfig.Setup()
}
