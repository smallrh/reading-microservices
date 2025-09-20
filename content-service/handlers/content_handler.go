package handlers

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"reading-microservices/shared/utils"
	"reading-microservices/content-service/models"
	"reading-microservices/content-service/services"
)

type ContentHandler struct {
	contentService services.ContentService
}

func NewContentHandler(contentService services.ContentService) *ContentHandler {
	return &ContentHandler{
		contentService: contentService,
	}
}

// Category handlers
func (h *ContentHandler) CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ERROR_INVALID_PARAMS, err.Error())
		return
	}

	category, err := h.contentService.CreateCategory(&req)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, category)
}

func (h *ContentHandler) GetCategories(c *gin.Context) {
	isActiveStr := c.Query("is_active")
	var isActive *bool
	if isActiveStr != "" {
		active := isActiveStr == "true"
		isActive = &active
	}

	categories, err := h.contentService.GetCategories(isActive)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, categories)
}

func (h *ContentHandler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	category, err := h.contentService.GetCategoryByID(id)
	if err != nil {
		utils.Error(c, utils.ERROR_NOT_FOUND, err.Error())
		return
	}

	utils.Success(c, category)
}

func (h *ContentHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	var req models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ERROR_INVALID_PARAMS, err.Error())
		return
	}

	if err := h.contentService.UpdateCategory(id, &req); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Category updated successfully", nil)
}

func (h *ContentHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	if err := h.contentService.DeleteCategory(id); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Category deleted successfully", nil)
}

// Tag handlers
func (h *ContentHandler) CreateTag(c *gin.Context) {
	var req models.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ERROR_INVALID_PARAMS, err.Error())
		return
	}

	tag, err := h.contentService.CreateTag(&req)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, tag)
}

func (h *ContentHandler) GetTags(c *gin.Context) {
	tags, err := h.contentService.GetTags()
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, tags)
}

func (h *ContentHandler) GetTagByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	tag, err := h.contentService.GetTagByID(id)
	if err != nil {
		utils.Error(c, utils.ERROR_NOT_FOUND, err.Error())
		return
	}

	utils.Success(c, tag)
}

// Novel handlers
func (h *ContentHandler) CreateNovel(c *gin.Context) {
	var req models.CreateNovelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ERROR_INVALID_PARAMS, err.Error())
		return
	}

	novel, err := h.contentService.CreateNovel(&req)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, novel)
}

func (h *ContentHandler) GetNovelByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	novel, err := h.contentService.GetNovelByID(id)
	if err != nil {
		utils.Error(c, utils.ERROR_NOT_FOUND, err.Error())
		return
	}

	// 增加浏览量
	views := int64(1)
	h.contentService.UpdateNovelStats(id, &views, nil, nil)

	utils.Success(c, novel)
}

func (h *ContentHandler) GetNovelsByCategory(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	novels, total, err := h.contentService.GetNovelsByCategory(categoryID, page, size)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.PageSuccess(c, novels, total, page, size)
}

func (h *ContentHandler) SearchNovels(c *gin.Context) {
	var params models.NovelSearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.Error(c, utils.ERROR_INVALID_PARAMS, err.Error())
		return
	}

	novels, total, err := h.contentService.SearchNovels(&params)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.PageSuccess(c, novels, total, params.Page, params.Size)
}

func (h *ContentHandler) UpdateNovel(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	var req models.UpdateNovelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ERROR_INVALID_PARAMS, err.Error())
		return
	}

	if err := h.contentService.UpdateNovel(id, &req); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Novel updated successfully", nil)
}

func (h *ContentHandler) DeleteNovel(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	if err := h.contentService.DeleteNovel(id); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Novel deleted successfully", nil)
}

func (h *ContentHandler) GetFeaturedNovels(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	novels, err := h.contentService.GetFeaturedNovels(limit)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, novels)
}

func (h *ContentHandler) GetLatestNovels(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	novels, err := h.contentService.GetLatestNovels(limit)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, novels)
}

// Chapter handlers
func (h *ContentHandler) CreateChapter(c *gin.Context) {
	var req models.CreateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ERROR_INVALID_PARAMS, err.Error())
		return
	}

	chapter, err := h.contentService.CreateChapter(&req)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.Success(c, chapter)
}

func (h *ContentHandler) GetChapterByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	chapter, err := h.contentService.GetChapterByID(id)
	if err != nil {
		utils.Error(c, utils.ERROR_NOT_FOUND, err.Error())
		return
	}

	utils.Success(c, chapter)
}

func (h *ContentHandler) GetChaptersByNovel(c *gin.Context) {
	novelID := c.Param("novel_id")
	if novelID == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "50"))

	chapters, total, err := h.contentService.GetChaptersByNovel(novelID, page, size)
	if err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.PageSuccess(c, chapters, total, page, size)
}

func (h *ContentHandler) GetChapterByNumber(c *gin.Context) {
	novelID := c.Param("novel_id")
	chapterNumberStr := c.Param("chapter_number")

	if novelID == "" || chapterNumberStr == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	chapterNumber, err := strconv.Atoi(chapterNumberStr)
	if err != nil {
		utils.Error(c, utils.ERROR_INVALID_PARAMS, "Invalid chapter number")
		return
	}

	chapter, err := h.contentService.GetChapterByNumber(novelID, chapterNumber)
	if err != nil {
		utils.Error(c, utils.ERROR_NOT_FOUND, err.Error())
		return
	}

	utils.Success(c, chapter)
}

func (h *ContentHandler) UpdateChapter(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	var req models.UpdateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ERROR_INVALID_PARAMS, err.Error())
		return
	}

	if err := h.contentService.UpdateChapter(id, &req); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Chapter updated successfully", nil)
}

func (h *ContentHandler) DeleteChapter(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorWithCode(c, utils.ERROR_INVALID_PARAMS)
		return
	}

	if err := h.contentService.DeleteChapter(id); err != nil {
		utils.Error(c, utils.ERROR, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Chapter deleted successfully", nil)
}

// Health check
func (h *ContentHandler) Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"service": "content-service",
	})
}