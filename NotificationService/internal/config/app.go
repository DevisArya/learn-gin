package config

import (
	"net"
	"os"

	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/delivery/grpcdelivery"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/delivery/middleware"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/delivery/rest/controller"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/delivery/rest/router"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/pb"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/repository"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type BootstrapConfigGrpc struct {
	DB       *gorm.DB
	Validate *validator.Validate
	Log      *logrus.Logger
}
type BootstrapConfigRest struct {
	DB       *gorm.DB
	App      *fiber.App
	Validate *validator.Validate
	Log      *logrus.Logger
}

type BootstrapResult struct {
	GRPCServer *grpc.Server
	Listener   net.Listener
}

func BootstrapGrpc(cfg *BootstrapConfigGrpc) (*BootstrapResult, error) {

	// Setup TCP listener
	network := os.Getenv("NETWORK")
	portGrpc := os.Getenv("NOTIF_PORT_GRPC")
	lis, err := net.Listen(network, portGrpc)
	if err != nil {
		return nil, err
	}

	//Init depedencies
	notifRepo := repository.NewNotificationRepository(cfg.Log)
	notifUc := usecase.NewNotificationUseCase(notifRepo, cfg.DB, cfg.Validate, cfg.Log)
	notifCtrl := grpcdelivery.NewNotificationGrpc(notifUc, cfg.Log)

	//init grpc server & register service

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, notifCtrl)

	return &BootstrapResult{
		GRPCServer: grpcServer,
		Listener:   lis,
	}, nil
}

func BootstrapRest(cfg *BootstrapConfigRest) {
	// setup repositories
	notificationRepository := repository.NewNotificationRepository(cfg.Log)

	// setup use cases
	notificationUseCase := usecase.NewNotificationUseCase(notificationRepository, cfg.DB, cfg.Validate, cfg.Log)

	// setup controller
	notificationController := controller.NewNotificationController(notificationUseCase, cfg.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth(notificationUseCase)

	routeConfig := router.RouteConfig{
		App:                    cfg.App,
		NotificationController: notificationController,
		AuthMiddleware:         authMiddleware,
	}
	routeConfig.Setup()
}
