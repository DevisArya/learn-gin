package entity

import (
	"time"
)

type ActionType string

const (
	ActionCreate ActionType = "CREATE"
	ActionUpdate ActionType = "UPDATE"
	ActionDelete ActionType = "DELETE"
)

type Notification struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint32     `gorm:"not null" json:"user_id"`
	Action      ActionType `gorm:"type:notification_action;not null" json:"action"`
	Description string     `gorm:"type:text" json:"description"`
	Seen        bool       `gorm:"default:false" json:"seen"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Notification) TableName() string {
	return "notifications"
}
