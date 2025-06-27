package repository

import (
	"context"

	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/dto"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	SaveNotification(ctx context.Context, tx *gorm.DB, dataNotif *entity.Notification) error
	GetNotifications(ctx context.Context, tx *gorm.DB, request *dto.GetNotificationRequest) ([]*entity.Notification, error)
}

type NotificationRepositoryImpl struct {
	Log *logrus.Logger
}

func NewNotificationRepository(log *logrus.Logger) NotificationRepository {
	return &NotificationRepositoryImpl{
		Log: log,
	}
}

func (r *NotificationRepositoryImpl) SaveNotification(ctx context.Context, tx *gorm.DB, dataNotif *entity.Notification) error {

	if err := tx.WithContext(ctx).
		Create(dataNotif).
		Error; err != nil {
		r.Log.Errorf("Failed to save notification: %v", err)
		return err
	}

	r.Log.Infof("Notification saved successfully for user_id=%v", dataNotif.UserID)
	return nil
}
func (r *NotificationRepositoryImpl) GetNotifications(ctx context.Context, tx *gorm.DB, request *dto.GetNotificationRequest) ([]*entity.Notification, error) {

	var notifications []*entity.Notification
	query := tx.WithContext(ctx).Where("user_id = ?", request.UserID)

	if request.Action != "" {
		query = query.Where("action = ?", request.Action)
	}

	if request.Seen != nil {
		query = query.Where("seen = ?", request.Seen)
	}

	err := query.Find(&notifications).Error
	if err != nil {
		r.Log.Errorf("Failed to fetch notifications: %v", err)
		return nil, err
	}

	r.Log.Infof("Fetched %d notifications for user_id=%v", len(notifications), request.UserID)
	return notifications, err
}
