package usecase

import (
	"context"
	"net/http"

	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/dto"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/entity"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/helper"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type NotificationUseCase interface {
	SaveNotification(ctx context.Context, dataNotif *dto.CreateNotificationRequest) error
	GetNotifications(ctx context.Context, request *dto.GetNotificationRequest) ([]*entity.Notification, error)
}

type NotificationUseCaseImpl struct {
	NotificationRepository repository.NotificationRepository
	DB                     *gorm.DB
	validate               *validator.Validate
	Log                    *logrus.Logger
}

func NewNotificationUseCase(notificationRepository repository.NotificationRepository, DB *gorm.DB, validate *validator.Validate, log *logrus.Logger) NotificationUseCase {
	return &NotificationUseCaseImpl{
		NotificationRepository: notificationRepository,
		DB:                     DB,
		validate:               validate,
		Log:                    log,
	}
}

// SaveNotification implements NotificationUseCase.
func (uc *NotificationUseCaseImpl) SaveNotification(ctx context.Context, request *dto.CreateNotificationRequest) error {
	uc.Log.Infof("Start saving notification for user_id=%v, action=%v", request.UserID, request.Action)

	if validationMessages := helper.ValidateStruct(uc.validate, request); len(validationMessages) > 0 {
		uc.Log.Warnf("Validation failed: %v", validationMessages)
		return helper.NewCustomError(http.StatusBadRequest, validationMessages)
	}

	tx := uc.DB.Begin()
	defer helper.CommitOrRollback(tx)

	dataNotif := entity.Notification{
		UserID:      request.UserID,
		Action:      entity.ActionType(request.Action),
		Description: request.Description,
	}

	err := uc.NotificationRepository.SaveNotification(ctx, tx, &dataNotif)
	if err != nil {
		uc.Log.Errorf("Failed to save notification: %v", err)
		return err
	}

	uc.Log.Infof("Notification successfully saved for user_id=%v", request.UserID)
	return nil
}

// GetNotifications implements NotificationUseCase.
func (uc *NotificationUseCaseImpl) GetNotifications(ctx context.Context, request *dto.GetNotificationRequest) ([]*entity.Notification, error) {

	uc.Log.Infof("Fetching notifications for user_id=%v", request.UserID)

	if validationMessages := helper.ValidateStruct(uc.validate, request); len(validationMessages) > 0 {
		uc.Log.Warnf("Validation failed: %v", validationMessages)
		return nil, helper.NewCustomError(http.StatusBadRequest, validationMessages)
	}

	tx := uc.DB.Begin()
	defer helper.CommitOrRollback(tx)

	res, err := uc.NotificationRepository.GetNotifications(ctx, tx, request)
	if err != nil {
		uc.Log.Errorf("Failed to fetch notifications: %v", err)
		return nil, err
	}

	uc.Log.Infof("Successfully fetched %d notifications for user_id=%v", len(res), request.UserID)
	return res, nil
}
