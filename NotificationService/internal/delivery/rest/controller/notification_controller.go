package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/dto"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/helper"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type NotificationController interface {
	GetNotifications(c *fiber.Ctx) error
}
type NotificationControllerImpl struct {
	NotificationUseCase usecase.NotificationUseCase
	Log                 *logrus.Logger
}

func NewNotificationController(NotificationUseCase usecase.NotificationUseCase, log *logrus.Logger) NotificationController {
	return &NotificationControllerImpl{
		NotificationUseCase: NotificationUseCase,
		Log:                 log,
	}
}

// GetNotifications implements NotificationController.
func (n *NotificationControllerImpl) GetNotifications(c *fiber.Ctx) error {

	var req dto.GetNotificationRequest

	if err := c.QueryParser(&req); err != nil {
		n.Log.Warnf("Invalid query parameters: %v", err)
		return helper.SendErrorResponse(c, http.StatusBadRequest, []string{"Invalid query parameters"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := n.NotificationUseCase.GetNotifications(ctx, &req)
	if err != nil {
		n.Log.Errorf("Failed to get notifications: %v", err)
		return helper.SendCustomErrorResponse(c, err)
	}

	n.Log.Infof("Successfully fetched notifications for UserID: %v", req.UserID)

	return helper.SendSuccessResponseWithData(c, http.StatusOK, "success get notifications", res)
}
