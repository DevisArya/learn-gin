package dto

import "github.com/DevisArya/BE-challenge-syn/NotificationService/internal/entity"

type GetNotificationRequest struct {
	UserID string `query:"user_id"`
	Action string `query:"action,omitempty"`
	Seen   *bool  `query:"seen,omitempty"`
}

type SaveNotificationResponse struct {
	Id uint `json:"id"`
}

type CreateNotificationRequest struct {
	UserID      uint32            `json:"user_id" validate:"required"`
	Action      entity.ActionType `json:"action" validate:"required"`
	Description string            `json:"description"`
}
