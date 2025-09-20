package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"reading-microservices/download-service/services"
)

type DownloadHandler struct {
	downloadService *services.DownloadService
}

func NewDownloadHandler(downloadService *services.DownloadService) *DownloadHandler {
	return &DownloadHandler{
		downloadService: downloadService,
	}
}

func (h *DownloadHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *DownloadHandler) CreateDownloadTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "CreateDownloadTask not implemented"})
}

func (h *DownloadHandler) GetDownloadTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetDownloadTask not implemented"})
}

func (h *DownloadHandler) GetDownloadTasks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetDownloadTasks not implemented"})
}

func (h *DownloadHandler) DownloadFile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "DownloadFile not implemented"})
}

func (h *DownloadHandler) UpdateDownloadTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "UpdateDownloadTask not implemented"})
}

func (h *DownloadHandler) DeleteDownloadTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "DeleteDownloadTask not implemented"})
}

func (h *DownloadHandler) StartDownload(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "StartDownload not implemented"})
}

func (h *DownloadHandler) PauseDownload(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PauseDownload not implemented"})
}

func (h *DownloadHandler) ResumeDownload(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ResumeDownload not implemented"})
}

func (h *DownloadHandler) GetDownloadStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetDownloadStats not implemented"})
}