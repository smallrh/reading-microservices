package services

import (
	"errors"
	"reading-microservices/content-service/models"
	"reading-microservices/content-service/repositories"
	"strings"
	"time"
)

type ContentService interface {
	// Category
	CreateCategory(req *models.CreateCategoryRequest) (*models.Category, error)
	GetCategories(isActive *bool) ([]models.Category, error)
	GetCategoryByID(id string) (*models.Category, error)
	UpdateCategory(id string, req *models.UpdateCategoryRequest) error
	DeleteCategory(id string) error

	// Tag
	CreateTag(req *models.CreateTagRequest) (*models.Tag, error)
	GetTags() ([]models.Tag, error)
	GetTagByID(id string) (*models.Tag, error)

	// Novel
	CreateNovel(req *models.CreateNovelRequest) (*models.Novel, error)
	GetNovelByID(id string) (*models.NovelDetailResponse, error)
	GetNovelsByCategory(categoryID string, page, size int) ([]models.NovelListResponse, int64, error)
	SearchNovels(params *models.NovelSearchParams) ([]models.NovelListResponse, int64, error)
	UpdateNovel(id string, req *models.UpdateNovelRequest) error
	DeleteNovel(id string) error
	UpdateNovelStats(novelID string, views *int64, rating *float64, ratingCount *int) error
	GetFeaturedNovels(limit int) ([]models.NovelListResponse, error)
	GetLatestNovels(limit int) ([]models.NovelListResponse, error)

	// Chapter
	CreateChapter(req *models.CreateChapterRequest) (*models.Chapter, error)
	GetChapterByID(id string) (*models.ChapterDetailResponse, error)
	GetChaptersByNovel(novelID string, page, size int) ([]models.ChapterSummary, int64, error)
	GetChapterByNumber(novelID string, chapterNumber int) (*models.ChapterDetailResponse, error)
	UpdateChapter(id string, req *models.UpdateChapterRequest) error
	DeleteChapter(id string) error
}

type contentService struct {
	repo repositories.ContentRepository
}

func NewContentService(repo repositories.ContentRepository) ContentService {
	return &contentService{repo: repo}
}

// Category methods
func (s *contentService) CreateCategory(req *models.CreateCategoryRequest) (*models.Category, error) {
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Color:       req.Color,
	}

	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}

	if err := s.repo.CreateCategory(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *contentService) GetCategories(isActive *bool) ([]models.Category, error) {
	return s.repo.GetCategories(isActive)
}

func (s *contentService) GetCategoryByID(id string) (*models.Category, error) {
	return s.repo.GetCategoryByID(id)
}

func (s *contentService) UpdateCategory(id string, req *models.UpdateCategoryRequest) error {
	category, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return err
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Description != nil {
		category.Description = req.Description
	}
	if req.Icon != nil {
		category.Icon = req.Icon
	}
	if req.Color != nil {
		category.Color = req.Color
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	return s.repo.UpdateCategory(category)
}

func (s *contentService) DeleteCategory(id string) error {
	return s.repo.DeleteCategory(id)
}

// Tag methods
func (s *contentService) CreateTag(req *models.CreateTagRequest) (*models.Tag, error) {
	// 检查标签是否已存在
	if _, err := s.repo.GetTagByName(req.Name); err == nil {
		return nil, errors.New("tag already exists")
	}

	tag := &models.Tag{
		Name:  req.Name,
		Color: req.Color,
	}

	if err := s.repo.CreateTag(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

func (s *contentService) GetTags() ([]models.Tag, error) {
	return s.repo.GetTags()
}

func (s *contentService) GetTagByID(id string) (*models.Tag, error) {
	return s.repo.GetTagByID(id)
}

// Novel methods
func (s *contentService) CreateNovel(req *models.CreateNovelRequest) (*models.Novel, error) {
	// 验证分类是否存在
	if _, err := s.repo.GetCategoryByID(req.CategoryID); err != nil {
		return nil, errors.New("category not found")
	}

	novel := &models.Novel{
		Title:       req.Title,
		Author:      req.Author,
		CoverURL:    req.CoverURL,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		Status:      "ongoing",
		IsFree:      true,
	}

	if req.Status != nil {
		novel.Status = *req.Status
	}
	if req.IsFeatured != nil {
		novel.IsFeatured = *req.IsFeatured
	}
	if req.IsFree != nil {
		novel.IsFree = *req.IsFree
	}
	if req.Price != nil {
		novel.Price = *req.Price
	}

	if err := s.repo.CreateNovel(novel); err != nil {
		return nil, err
	}

	// 添加标签
	if len(req.TagIDs) > 0 {
		s.repo.AddNovelTags(novel.ID, req.TagIDs)
	}

	return novel, nil
}

func (s *contentService) GetNovelByID(id string) (*models.NovelDetailResponse, error) {
	novel, err := s.repo.GetNovelByID(id)
	if err != nil {
		return nil, err
	}

	// 获取最新章节
	latestChapters, _ := s.repo.GetLatestChapters(id, 5)

	response := &models.NovelDetailResponse{
		NovelListResponse: s.convertToNovelListResponse(novel),
		LatestChapters:    s.convertToChapterSummaries(latestChapters),
	}

	return response, nil
}

func (s *contentService) GetNovelsByCategory(categoryID string, page, size int) ([]models.NovelListResponse, int64, error) {
	novels, total, err := s.repo.GetNovelsByCategory(categoryID, page, size)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.NovelListResponse, len(novels))
	for i, novel := range novels {
		responses[i] = *s.convertToNovelListResponse(&novel)
	}

	return responses, total, nil
}

func (s *contentService) SearchNovels(params *models.NovelSearchParams) ([]models.NovelListResponse, int64, error) {
	novels, total, err := s.repo.SearchNovels(params)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.NovelListResponse, len(novels))
	for i, novel := range novels {
		responses[i] = *s.convertToNovelListResponse(&novel)
	}

	return responses, total, nil
}

func (s *contentService) UpdateNovel(id string, req *models.UpdateNovelRequest) error {
	novel, err := s.repo.GetNovelByID(id)
	if err != nil {
		return err
	}

	if req.Title != nil {
		novel.Title = *req.Title
	}
	if req.Author != nil {
		novel.Author = *req.Author
	}
	if req.CoverURL != nil {
		novel.CoverURL = req.CoverURL
	}
	if req.Description != nil {
		novel.Description = req.Description
	}
	if req.CategoryID != nil {
		novel.CategoryID = *req.CategoryID
	}
	if req.Status != nil {
		novel.Status = *req.Status
	}
	if req.IsFeatured != nil {
		novel.IsFeatured = *req.IsFeatured
	}
	if req.IsFree != nil {
		novel.IsFree = *req.IsFree
	}
	if req.Price != nil {
		novel.Price = *req.Price
	}

	if err := s.repo.UpdateNovel(novel); err != nil {
		return err
	}

	// 更新标签
	if req.TagIDs != nil {
		// 移除所有现有标签
		currentTags, _ := s.repo.GetNovelTags(id)
		currentTagIDs := make([]string, len(currentTags))
		for i, tag := range currentTags {
			currentTagIDs[i] = tag.ID
		}
		s.repo.RemoveNovelTags(id, currentTagIDs)

		// 添加新标签
		if len(req.TagIDs) > 0 {
			s.repo.AddNovelTags(id, req.TagIDs)
		}
	}

	return nil
}

func (s *contentService) DeleteNovel(id string) error {
	return s.repo.DeleteNovel(id)
}

func (s *contentService) UpdateNovelStats(novelID string, views *int64, rating *float64, ratingCount *int) error {
	return s.repo.UpdateNovelStats(novelID, views, rating, ratingCount)
}

func (s *contentService) GetFeaturedNovels(limit int) ([]models.NovelListResponse, error) {
	novels, err := s.repo.GetFeaturedNovels(limit)
	if err != nil {
		return nil, err
	}

	responses := make([]models.NovelListResponse, len(novels))
	for i, novel := range novels {
		responses[i] = *s.convertToNovelListResponse(&novel)
	}

	return responses, nil
}

func (s *contentService) GetLatestNovels(limit int) ([]models.NovelListResponse, error) {
	novels, err := s.repo.GetLatestNovels(limit)
	if err != nil {
		return nil, err
	}

	responses := make([]models.NovelListResponse, len(novels))
	for i, novel := range novels {
		responses[i] = *s.convertToNovelListResponse(&novel)
	}

	return responses, nil
}

// Chapter methods
func (s *contentService) CreateChapter(req *models.CreateChapterRequest) (*models.Chapter, error) {
	// 验证小说是否存在
	if _, err := s.repo.GetNovelByID(req.NovelID); err != nil {
		return nil, errors.New("novel not found")
	}

	// 检查章节号是否已存在
	if _, err := s.repo.GetChapterByNumber(req.NovelID, req.ChapterNumber); err == nil {
		return nil, errors.New("chapter number already exists")
	}

	chapter := &models.Chapter{
		NovelID:       req.NovelID,
		Title:         req.Title,
		Content:       req.Content,
		ChapterNumber: req.ChapterNumber,
		WordCount:     len(strings.Fields(req.Content)),
		IsFree:        true,
	}

	if req.IsFree != nil {
		chapter.IsFree = *req.IsFree
	}
	if req.Price != nil {
		chapter.Price = *req.Price
	}

	if err := s.repo.CreateChapter(chapter); err != nil {
		return nil, err
	}

	return chapter, nil
}

func (s *contentService) GetChapterByID(id string) (*models.ChapterDetailResponse, error) {
	chapter, err := s.repo.GetChapterByID(id)
	if err != nil {
		return nil, err
	}

	return s.convertToChapterDetailResponse(chapter), nil
}

func (s *contentService) GetChaptersByNovel(novelID string, page, size int) ([]models.ChapterSummary, int64, error) {
	chapters, total, err := s.repo.GetChaptersByNovel(novelID, page, size)
	if err != nil {
		return nil, 0, err
	}

	summaries := s.convertToChapterSummaries(chapters)
	return summaries, total, nil
}

func (s *contentService) GetChapterByNumber(novelID string, chapterNumber int) (*models.ChapterDetailResponse, error) {
	chapter, err := s.repo.GetChapterByNumber(novelID, chapterNumber)
	if err != nil {
		return nil, err
	}

	return s.convertToChapterDetailResponse(chapter), nil
}

func (s *contentService) UpdateChapter(id string, req *models.UpdateChapterRequest) error {
	chapter, err := s.repo.GetChapterByID(id)
	if err != nil {
		return err
	}

	if req.Title != nil {
		chapter.Title = *req.Title
	}
	if req.Content != nil {
		chapter.Content = *req.Content
		chapter.WordCount = len(strings.Fields(*req.Content))
	}
	if req.IsFree != nil {
		chapter.IsFree = *req.IsFree
	}
	if req.Price != nil {
		chapter.Price = *req.Price
	}

	return s.repo.UpdateChapter(chapter)
}

func (s *contentService) DeleteChapter(id string) error {
	return s.repo.DeleteChapter(id)
}

// Helper methods
func (s *contentService) convertToNovelListResponse(novel *models.Novel) *models.NovelListResponse {
	return &models.NovelListResponse{
		ID:            novel.ID,
		Title:         novel.Title,
		Author:        novel.Author,
		CoverURL:      novel.CoverURL,
		Description:   novel.Description,
		Category:      novel.Category,
		Status:        novel.Status,
		TotalChapters: novel.TotalChapters,
		WordCount:     novel.WordCount,
		ViewsCount:    novel.ViewsCount,
		Rating:        novel.Rating,
		RatingCount:   novel.RatingCount,
		IsFeatured:    novel.IsFeatured,
		IsFree:        novel.IsFree,
		Price:         novel.Price,
		Tags:          novel.Tags,
		CreatedAt:     novel.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     novel.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *contentService) convertToChapterSummaries(chapters []models.Chapter) []models.ChapterSummary {
	summaries := make([]models.ChapterSummary, len(chapters))
	for i, chapter := range chapters {
		summaries[i] = models.ChapterSummary{
			ID:            chapter.ID,
			Title:         chapter.Title,
			ChapterNumber: chapter.ChapterNumber,
			WordCount:     chapter.WordCount,
			IsFree:        chapter.IsFree,
			Price:         chapter.Price,
			CreatedAt:     chapter.CreatedAt.Format(time.RFC3339),
		}
	}
	return summaries
}

func (s *contentService) convertToChapterDetailResponse(chapter *models.Chapter) *models.ChapterDetailResponse {
	response := &models.ChapterDetailResponse{
		ID:            chapter.ID,
		NovelID:       chapter.NovelID,
		Title:         chapter.Title,
		Content:       chapter.Content,
		ChapterNumber: chapter.ChapterNumber,
		WordCount:     chapter.WordCount,
		IsFree:        chapter.IsFree,
		Price:         chapter.Price,
		CreatedAt:     chapter.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     chapter.UpdatedAt.Format(time.RFC3339),
	}

	// 获取前后章节
	if prevChapter, err := s.repo.GetChapterByNumber(chapter.NovelID, chapter.ChapterNumber-1); err == nil {
		response.PrevChapter = &models.ChapterNav{
			ID:            prevChapter.ID,
			Title:         prevChapter.Title,
			ChapterNumber: prevChapter.ChapterNumber,
		}
	}

	if nextChapter, err := s.repo.GetChapterByNumber(chapter.NovelID, chapter.ChapterNumber+1); err == nil {
		response.NextChapter = &models.ChapterNav{
			ID:            nextChapter.ID,
			Title:         nextChapter.Title,
			ChapterNumber: nextChapter.ChapterNumber,
		}
	}

	return response
}
