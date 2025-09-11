package service

import (
	"github.com/Kevinmajesta/depublic-backend/internal/entity"
	"github.com/Kevinmajesta/depublic-backend/internal/repository"
	"github.com/Kevinmajesta/depublic-backend/pkg/token"
	"github.com/google/uuid"
)

type NotificationService interface {
	GetUserNotifications(userID uuid.UUID) ([]*entity.Notification, error)
	GetUserNotificationsNoRead(userID uuid.UUID) ([]*entity.Notification, error)
	CreateNotification(notification *entity.Notification) error
	MarkNotificationAsRead(notificationID uuid.UUID) error
	CheckUserExists(id uuid.UUID) (bool, error)
}

type notificationService struct {
	notificationRepo repository.NotificationRepository
	tokenUseCase     token.TokenUseCase
	userRepository   repository.UserRepository
}

func NewNotificationService(notificationRepo repository.NotificationRepository, tokenUseCase token.TokenUseCase,
	userRepository repository.UserRepository) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		tokenUseCase:     tokenUseCase,
		userRepository:   userRepository,
	}
}

func (s *notificationService) CheckUserExists(id uuid.UUID) (bool, error) {
	return s.notificationRepo.CheckUserExists(id)
}

func (s *notificationService) GetUserNotificationsNoRead(userID uuid.UUID) ([]*entity.Notification, error) {
	return s.notificationRepo.GetUserNotificationsNoRead(userID)
}

func (s *notificationService) GetUserNotifications(userID uuid.UUID) ([]*entity.Notification, error) {
	return s.notificationRepo.GetUserNotifications(userID)
}

func (s *notificationService) CreateNotification(notification *entity.Notification) error {
	userIds, err := s.userRepository.GetAllUserIds()
	if err != nil {
		return err
	}

	for _, userID := range userIds {
		// Buat salinan notifikasi untuk setiap pengguna
		userNotification := *notification
		userNotification.UserID = userID
		if err := s.notificationRepo.CreateNotification(&userNotification); err != nil {
			return err
		}
	}

	return nil
}

func (s *notificationService) MarkNotificationAsRead(notificationID uuid.UUID) error {
	return s.notificationRepo.MarkNotificationAsRead(notificationID)
}
