package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/dto"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/helper"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type NotificationController interface {
	GetNotifications(c *gin.Context)
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
func (n *NotificationControllerImpl) GetNotifications(c *gin.Context) {

	var req dto.GetNotificationRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		n.Log.Warnf("Invalid query parameters: %v", err)
		helper.SendErrorResponse(c, http.StatusBadRequest, []string{"Invalid query parameters"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := n.NotificationUseCase.GetNotifications(ctx, &req)
	if err != nil {
		n.Log.Errorf("Failed to get notifications: %v", err)
		helper.SendCustomErrorResponse(c, err)
		return
	}

	n.Log.Infof("Successfully fetched notifications for UserID: %v", req.UserID)

	helper.SendSuccessResponseWithData(c, http.StatusOK, "success get notifications", res)
}
