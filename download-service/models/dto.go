package models

import (
	"crypto/rand"
	"fmt"
)

func generateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// 请求DTO
type CreateDownloadTaskRequest struct {
	NovelID string `json:"novel_id" binding:"required"`
	Format  string `json:"format" binding:"required,oneof=txt epub pdf"`
}

type UpdateDownloadTaskRequest struct {
	Status string `json:"status" binding:"required,oneof=pending downloading paused completed failed"`
}

// 响应DTO
type DownloadTaskResponse struct {
	ID                 string                     `json:"id"`
	NovelID            string                     `json:"novel_id"`
	Status             string                     `json:"status"`
	Progress           float64                    `json:"progress"`
	DownloadedChapters int                        `json:"downloaded_chapters"`
	TotalChapters      int                        `json:"total_chapters"`
	FileSize           int64                      `json:"file_size"`
	DownloadSpeed      *string                    `json:"download_speed"`
	ErrorMessage       *string                    `json:"error_message"`
	StartedAt          *string                    `json:"started_at"`
	CompletedAt        *string                    `json:"completed_at"`
	CreatedAt          string                     `json:"created_at"`
	UpdatedAt          string                     `json:"updated_at"`
	Novel              *NovelInfo                 `json:"novel,omitempty"`
	DownloadURL        *string                    `json:"download_url,omitempty"`
	DownloadChapters   []DownloadChapterResponse  `json:"download_chapters,omitempty"`
}

type DownloadChapterResponse struct {
	ID            string  `json:"id"`
	ChapterID     string  `json:"chapter_id"`
	ChapterNumber int     `json:"chapter_number"`
	FileSize      int64   `json:"file_size"`
	IsDownloaded  bool    `json:"is_downloaded"`
	DownloadedAt  *string `json:"downloaded_at"`
	CreatedAt     string  `json:"created_at"`
}

type NovelInfo struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	Author        string  `json:"author"`
	CoverURL      *string `json:"cover_url"`
	Status        string  `json:"status"`
	TotalChapters int     `json:"total_chapters"`
}

type DownloadStatsResponse struct {
	Total       int64 `json:"total"`
	Pending     int64 `json:"pending"`
	Downloading int64 `json:"downloading"`
	Paused      int64 `json:"paused"`
	Completed   int64 `json:"completed"`
	Failed      int64 `json:"failed"`
	TotalSize   int64 `json:"total_size"` // 字节
}

// 下载状态常量
const (
	StatusPending     = "pending"
	StatusDownloading = "downloading"
	StatusPaused      = "paused"
	StatusCompleted   = "completed"
	StatusFailed      = "failed"
)

// 文件格式常量
const (
	FormatTXT  = "txt"
	FormatEPUB = "epub"
	FormatPDF  = "pdf"
)