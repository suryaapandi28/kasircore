package entity

import (
	"github.com/google/uuid"
)

type Notification struct {
	Notification_ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"notification_id"`
	UserID          uuid.UUID `json:"user_id"`
	Type            string    `json:"type"`
	Message         string    `json:"message"`
	IsRead          bool      `json:"is_read"`
	Auditable
}

func NewNotification(tipe, message string, isRead bool) *Notification {
	return &Notification{
		Notification_ID: uuid.New(),
		Type:            tipe,
		Message:         message,
		IsRead:          isRead,
		Auditable:       NewAuditable(),
	}
}

func MarkNotificationAsRead(notificationID uuid.UUID) *Notification {
	return &Notification{
		Notification_ID: notificationID,
		IsRead:          true,
		Auditable:       UpdateAuditable(),
	}
}

func UserRequest(notificationID uuid.UUID) *Notification {
	return &Notification{
		Notification_ID: notificationID,
		IsRead:          true,
		Auditable:       UpdateAuditable(),
	}
}
