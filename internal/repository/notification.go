package repository

import (
	"github.com/google/uuid"
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/pkg/cache"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	GetUserNotifications(userID uuid.UUID) ([]*entity.Notification, error)
	GetUserNotificationsNoRead(userID uuid.UUID) ([]*entity.Notification, error)
	CreateNotification(notification *entity.Notification) error
	MarkNotificationAsRead(notificationID uuid.UUID) error
	CheckUserExists(id uuid.UUID) (bool, error)
}

type notificationRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewNotificationRepository(db *gorm.DB, cacheable cache.Cacheable) *notificationRepository {
	return &notificationRepository{db: db, cacheable: cacheable}
}

//--------------PAKE CACHE TAPI ERROR GA KELUAR----------------//
// func (r *notificationRepository) GetUserNotifications(userID uuid.UUID) ([]*entity.Notification, error) {
// 	notifications := make([]*entity.Notification, 0)
// 	key := "user_notifications"

// 	data, _ := r.cacheable.Get(key)
// 	if data == "" {
// 		result := r.db.Where("user_id = ?", userID).Find(&notifications)
// 		if result.Error != nil {
// 			return nil, result.Error
// 		}
// 		marshalledNotifications, _ := json.Marshal(notifications)
// 		err := r.cacheable.Set(key, marshalledNotifications, 5*time.Minute)
// 		if err != nil {
// 			return nil, err
// 		}
// 	} else {
// 		err := json.Unmarshal([]byte(data), &notifications)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	for _, notification := range notifications {
// 		err := r.MarkNotificationAsRead(notification.Notification_ID)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return notifications, nil
// }

// --------------GAK PAKE CACHE----------------//
func (r *notificationRepository) GetUserNotifications(userID uuid.UUID) ([]*entity.Notification, error) {
	notifications := make([]*entity.Notification, 0)
	result := r.db.Where("user_id = ?", userID).Find(&notifications)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, notification := range notifications {
		err := r.MarkNotificationAsRead(notification.Notification_ID)
		if err != nil {
			return nil, err
		}
	}

	if len(notifications) == 0 {
		return notifications, nil
	}
	return notifications, nil
}

func (r *notificationRepository) CheckUserExists(id uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Notification{}).Where("user_id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *notificationRepository) GetUserNotificationsNoRead(userID uuid.UUID) ([]*entity.Notification, error) {
	var notifications []*entity.Notification
	result := r.db.Where("user_id = ? AND is_read = ?", userID, false).Find(&notifications)
	if result.Error != nil {
		return nil, result.Error
	}

	if len(notifications) == 0 {
		return notifications, nil
	}

	err := r.db.Model(&entity.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *notificationRepository) CreateNotification(notification *entity.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) MarkNotificationAsRead(notificationID uuid.UUID) error {
	return r.db.Model(&entity.Notification{}).Where("notification_id = ?", notificationID).Update("is_read", true).Error
}
