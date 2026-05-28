package model

import (
	"encoding/json"
	"gochat/internal/notification/domain"
	"gochat/internal/shared/kernel"
)

type Notification struct {
	ID kernel.MessageID `gorm:"primaryKey;type:char(36)"`
	// 建立索引 recipient_id, state, 以及联合索引 (recipient_id, state)
	RecipientID kernel.UserID            `gorm:"type:char(36);not null;index:idx_notification_notifications_recipient_id;index:idx_notification_notifications_recipient_state,priority:1"`
	State       domain.NotificationState `gorm:"type:varchar(36);not null;index:idx_notification_notifications_state;index:idx_notification_notifications_recipient_state,priority:2"`
	RawPayload  json.RawMessage          `gorm:"type:longblob"`
}

func (*Notification) TableName() string {
	return "notification_notifications"
}
