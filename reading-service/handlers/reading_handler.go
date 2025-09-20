package handlers

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"reading-microservices/shared/utils"
	"reading-microservices/reading-service/models"
	"reading-microservices/reading-service/services"
)

type ReadingHandler struct {
	readingService services.ReadingService
}

func NewReadingHandler(readingService services.ReadingService) *ReadingHandler {
	return &ReadingHandler{
		readingService: readingService,
	}
}

// Reading Progress Handlers
func (h *ReadingHandler) UpdateReadingProgress(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	var req models.UpdateReadingProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ERROR_INVALID_PARAMS, err.Error())
		return
	}

	if err := h.readingService.UpdateReadingProgress(userID, &req); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Reading progress updated successfully", nil)
}

func (h *ReadingHandler) GetReadingProgress(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	novelID := c.Param("novel_id")
	chapterID := c.Param("chapter_id")

	if novelID == "" || chapterID == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	progress, err := h.readingService.GetReadingProgress(userID, novelID, chapterID)
	if err != nil {
		utils.Error(c, utils.ERROR_NOT_FOUND, err.Error())
		return
	}

	utils.Success(c, progress)
}

func (h *ReadingHandler) GetReadingHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	history, total, err := h.readingService.GetUserReadingHistory(userID, page, size)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.PageSuccess(c, history, total, page, size)
}

// Bookshelf Handlers
func (h *ReadingHandler) AddToBookshelf(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	var req models.AddToBookshelfRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ERROR_INVALID_PARAMS, err.Error())
		return
	}

	if err := h.readingService.AddToBookshelf(userID, &req); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Added to bookshelf successfully", nil)
}

func (h *ReadingHandler) RemoveFromBookshelf(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	novelID := c.Param("novel_id")
	shelfType := c.Query("shelf_type")

	if novelID == "" || shelfType == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	if err := h.readingService.RemoveFromBookshelf(userID, novelID, shelfType); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Removed from bookshelf successfully", nil)
}

func (h *ReadingHandler) GetBookshelf(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	shelfType := c.Query("shelf_type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	bookshelf, total, err := h.readingService.GetBookshelf(userID, shelfType, page, size)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.PageSuccess(c, bookshelf, total, page, size)
}

func (h *ReadingHandler) GetBookshelfStats(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	stats, err := h.readingService.GetBookshelfStats(userID)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, stats)
}

// Favorite Handlers
func (h *ReadingHandler) AddFavorite(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	novelID := c.Param("novel_id")
	if novelID == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	if err := h.readingService.AddFavorite(userID, novelID); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Added to favorites successfully", nil)
}

func (h *ReadingHandler) RemoveFavorite(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	novelID := c.Param("novel_id")
	if novelID == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	if err := h.readingService.RemoveFavorite(userID, novelID); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Removed from favorites successfully", nil)
}

func (h *ReadingHandler) GetFavorites(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	favorites, total, err := h.readingService.GetUserFavorites(userID, page, size)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.PageSuccess(c, favorites, total, page, size)
}

func (h *ReadingHandler) CheckFavoriteStatus(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	novelID := c.Param("novel_id")
	if novelID == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	isFavorite := h.readingService.CheckFavoriteStatus(userID, novelID)

	utils.Success(c, gin.H{
		"is_favorite": isFavorite,
	})
}

// Comment Handlers
func (h *ReadingHandler) CreateComment(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ERROR_INVALID_PARAMS, err.Error())
		return
	}

	comment, err := h.readingService.CreateComment(userID, &req)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, comment)
}

func (h *ReadingHandler) UpdateComment(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	commentID := c.Param("id")
	if commentID == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	var req models.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ERROR_INVALID_PARAMS, err.Error())
		return
	}

	if err := h.readingService.UpdateComment(userID, commentID, &req); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Comment updated successfully", nil)
}

func (h *ReadingHandler) DeleteComment(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	commentID := c.Param("id")
	if commentID == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	if err := h.readingService.DeleteComment(userID, commentID); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Comment deleted successfully", nil)
}

func (h *ReadingHandler) GetNovelComments(c *gin.Context) {
	novelID := c.Param("novel_id")
	if novelID == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	comments, total, err := h.readingService.GetNovelComments(novelID, page, size)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.PageSuccess(c, comments, total, page, size)
}

func (h *ReadingHandler) GetChapterComments(c *gin.Context) {
	chapterID := c.Param("chapter_id")
	if chapterID == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	comments, total, err := h.readingService.GetChapterComments(chapterID, page, size)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.PageSuccess(c, comments, total, page, size)
}

func (h *ReadingHandler) GetUserComments(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	comments, total, err := h.readingService.GetUserComments(userID, page, size)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.PageSuccess(c, comments, total, page, size)
}

func (h *ReadingHandler) LikeComment(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	commentID := c.Param("id")
	if commentID == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	if err := h.readingService.LikeComment(userID, commentID); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Comment liked successfully", nil)
}

func (h *ReadingHandler) UnlikeComment(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	commentID := c.Param("id")
	if commentID == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	if err := h.readingService.UnlikeComment(userID, commentID); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Comment unliked successfully", nil)
}

// Search Handlers
func (h *ReadingHandler) AddSearchHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	var req struct {
		Keyword     string `json:"keyword" binding:"required"`
		SearchType  string `json:"search_type"`
		ResultCount int    `json:"result_count"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ERROR_INVALID_PARAMS, err.Error())
		return
	}

	if err := h.readingService.AddSearchHistory(userID, req.Keyword, req.SearchType, req.ResultCount); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Search history added", nil)
}

func (h *ReadingHandler) GetSearchHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	history, err := h.readingService.GetSearchHistory(userID)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, history)
}

func (h *ReadingHandler) ClearSearchHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	if err := h.readingService.ClearSearchHistory(userID); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Search history cleared", nil)
}

// Statistics Handlers
func (h *ReadingHandler) GetReadingStats(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
		return
	}

	stats, err := h.readingService.GetReadingStats(userID)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, stats)
}

// Health check
func (h *ReadingHandler) Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"service": "reading-service",
	})
}