package services

import (
	"reading-microservices/notification-service/models"
	"reading-microservices/notification-service/repositories"
)

type NotificationService struct {
	notificationRepo *repositories.NotificationRepository
}

func NewNotificationService(notificationRepo *repositories.NotificationRepository) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
	}
}

func (s *NotificationService) CreateNotification(notification *models.Notification) error {
	return s.notificationRepo.Create(notification)
}

func (s *NotificationService) GetNotificationsByUserID(userID uint, limit, offset int) ([]*models.Notification, error) {
	return s.notificationRepo.FindByUserID(userID, limit, offset)
}

func (s *NotificationService) MarkAsRead(id uint) error {
	notification, err := s.notificationRepo.FindByID(id)
	if err != nil {
		return err
	}

	notification.IsRead = true
	return s.notificationRepo.Update(notification)
}

func (s *NotificationService) DeleteNotification(id uint) error {
	return s.notificationRepo.Delete(id)
}