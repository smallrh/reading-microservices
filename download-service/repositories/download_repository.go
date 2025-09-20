package repositories

import (
	"gorm.io/gorm"
	"reading-microservices/download-service/models"
)

type DownloadRepository struct {
	db *gorm.DB
}

func NewDownloadRepository(db *gorm.DB) *DownloadRepository {
	return &DownloadRepository{db: db}
}

func (r *DownloadRepository) Create(download *models.DownloadTask) error {
	return r.db.Create(download).Error
}

func (r *DownloadRepository) FindByID(id string) (*models.DownloadTask, error) {
	var download models.DownloadTask
	err := r.db.First(&download, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &download, nil
}

func (r *DownloadRepository) FindByUserID(userID string, limit, offset int) ([]*models.DownloadTask, error) {
	var downloads []*models.DownloadTask
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&downloads).Error
	return downloads, err
}

func (r *DownloadRepository) Update(download *models.DownloadTask) error {
	return r.db.Save(download).Error
}