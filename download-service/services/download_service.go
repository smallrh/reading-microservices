package services

import (
	"reading-microservices/download-service/models"
	"reading-microservices/download-service/repositories"
)

type DownloadService struct {
	downloadRepo *repositories.DownloadRepository
}

func NewDownloadService(downloadRepo *repositories.DownloadRepository) *DownloadService {
	return &DownloadService{
		downloadRepo: downloadRepo,
	}
}

func (s *DownloadService) CreateDownload(download *models.DownloadTask) error {
	return s.downloadRepo.Create(download)
}

func (s *DownloadService) GetDownloadByID(id string) (*models.DownloadTask, error) {
	return s.downloadRepo.FindByID(id)
}

func (s *DownloadService) GetDownloadsByUserID(userID string, limit, offset int) ([]*models.DownloadTask, error) {
	return s.downloadRepo.FindByUserID(userID, limit, offset)
}

func (s *DownloadService) UpdateDownload(download *models.DownloadTask) error {
	return s.downloadRepo.Update(download)
}