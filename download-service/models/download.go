package models

import (
	"time"
	"gorm.io/gorm"
)

type DownloadTask struct {
	ID                 string     `gorm:"type:varchar(36);primarykey" json:"id"`
	UserID             string     `gorm:"type:varchar(36);not null;index" json:"user_id"`
	NovelID            string     `gorm:"type:varchar(36);not null" json:"novel_id"`
	Status             string     `gorm:"type:enum('pending','downloading','paused','completed','failed');default:'pending';index" json:"status"`
	Progress           float64    `gorm:"type:decimal(5,2);default:0.00" json:"progress"`
	DownloadedChapters int        `gorm:"default:0" json:"downloaded_chapters"`
	TotalChapters      int        `gorm:"not null" json:"total_chapters"`
	FileSize           int64      `gorm:"default:0" json:"file_size"`
	DownloadSpeed      *string    `gorm:"type:varchar(20)" json:"download_speed"`
	ErrorMessage       *string    `gorm:"type:text" json:"error_message"`
	StartedAt          *time.Time `json:"started_at"`
	CompletedAt        *time.Time `json:"completed_at"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	// 关联
	DownloadChapters []DownloadChapter `gorm:"foreignKey:DownloadTaskID" json:"download_chapters,omitempty"`
}

type DownloadChapter struct {
	ID             string     `gorm:"type:varchar(36);primarykey" json:"id"`
	DownloadTaskID string     `gorm:"type:varchar(36);not null;index" json:"download_task_id"`
	ChapterID      string     `gorm:"type:varchar(36);not null" json:"chapter_id"`
	ChapterNumber  int        `gorm:"not null" json:"chapter_number"`
	FilePath       *string    `gorm:"type:varchar(500)" json:"file_path"`
	FileSize       int64      `gorm:"default:0" json:"file_size"`
	IsDownloaded   bool       `gorm:"default:false" json:"is_downloaded"`
	DownloadedAt   *time.Time `json:"downloaded_at"`
	CreatedAt      time.Time  `json:"created_at"`

	// 关联
	DownloadTask DownloadTask `gorm:"foreignKey:DownloadTaskID" json:"download_task,omitempty"`
}

// Hook functions
func (dt *DownloadTask) BeforeCreate(tx *gorm.DB) error {
	if dt.ID == "" {
		dt.ID = generateUUID()
	}
	return nil
}

func (dc *DownloadChapter) BeforeCreate(tx *gorm.DB) error {
	if dc.ID == "" {
		dc.ID = generateUUID()
	}
	return nil
}