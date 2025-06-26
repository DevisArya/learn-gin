package dto

type SendUserNotificationRequest struct {
	UserID      uint32 `json:"user_id" validate:"required"`
	Action      string `json:"action" validate:"required"`
	Description string `json:"description" validate:"required"`
}
