package services

import (
	"errors"
	"reading-microservices/reading-service/models"
	"reading-microservices/reading-service/repositories"
	"time"
)

type ReadingService interface {
	// Reading Progress
	UpdateReadingProgress(userID string, req *models.UpdateReadingProgressRequest) error
	GetReadingProgress(userID, novelID, chapterID string) (*models.ReadingRecordResponse, error)
	GetUserReadingHistory(userID string, page, size int) ([]models.ReadingRecordResponse, int64, error)

	// Bookshelf
	AddToBookshelf(userID string, req *models.AddToBookshelfRequest) error
	RemoveFromBookshelf(userID, novelID, shelfType string) error
	GetBookshelf(userID, shelfType string, page, size int) ([]models.BookshelfResponse, int64, error)
	GetBookshelfStats(userID string) (*models.BookshelfStatsResponse, error)

	// Favorites
	AddFavorite(userID, novelID string) error
	RemoveFavorite(userID, novelID string) error
	GetUserFavorites(userID string, page, size int) ([]models.BookshelfResponse, int64, error)
	CheckFavoriteStatus(userID, novelID string) bool

	// Comments
	CreateComment(userID string, req *models.CreateCommentRequest) (*models.CommentResponse, error)
	UpdateComment(userID, commentID string, req *models.UpdateCommentRequest) error
	DeleteComment(userID, commentID string) error
	GetNovelComments(novelID string, page, size int) ([]models.CommentResponse, int64, error)
	GetChapterComments(chapterID string, page, size int) ([]models.CommentResponse, int64, error)
	GetUserComments(userID string, page, size int) ([]models.CommentResponse, int64, error)
	LikeComment(userID, commentID string) error
	UnlikeComment(userID, commentID string) error

	// Search
	AddSearchHistory(userID, keyword, searchType string, resultCount int) error
	GetSearchHistory(userID string) ([]models.SearchHistory, error)
	ClearSearchHistory(userID string) error

	// Statistics
	GetReadingStats(userID string) (*models.ReadingStatsResponse, error)
}

type readingService struct {
	repo repositories.ReadingRepository
}

func NewReadingService(repo repositories.ReadingRepository) ReadingService {
	return &readingService{repo: repo}
}

// Reading Progress
func (s *readingService) UpdateReadingProgress(userID string, req *models.UpdateReadingProgressRequest) error {
	record := &models.ReadingRecord{
		UserID:          userID,
		NovelID:         req.NovelID,
		ChapterID:       req.ChapterID,
		ChapterNumber:   req.ChapterNumber,
		ReadingPosition: req.ReadingPosition,
		ReadingProgress: req.ReadingProgress,
		ReadingTime:     req.ReadingTime,
		LastReadAt:      time.Now(),
	}

	// 更新阅读记录
	if err := s.repo.CreateOrUpdateReadingRecord(record); err != nil {
		return err
	}

	// 更新书架进度
	return s.repo.UpdateBookshelfProgress(userID, req.NovelID, req.ReadingProgress)
}

func (s *readingService) GetReadingProgress(userID, novelID, chapterID string) (*models.ReadingRecordResponse, error) {
	record, err := s.repo.GetReadingRecord(userID, novelID, chapterID)
	if err != nil {
		return nil, err
	}

	return s.convertToReadingRecordResponse(record), nil
}

func (s *readingService) GetUserReadingHistory(userID string, page, size int) ([]models.ReadingRecordResponse, int64, error) {
	records, total, err := s.repo.GetUserReadingRecords(userID, page, size)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.ReadingRecordResponse, len(records))
	for i, record := range records {
		responses[i] = *s.convertToReadingRecordResponse(&record)
	}

	return responses, total, nil
}

// Bookshelf
func (s *readingService) AddToBookshelf(userID string, req *models.AddToBookshelfRequest) error {
	bookshelf := &models.Bookshelf{
		UserID:    userID,
		NovelID:   req.NovelID,
		ShelfType: req.ShelfType,
		AddedAt:   time.Now(),
	}

	return s.repo.AddToBookshelf(bookshelf)
}

func (s *readingService) RemoveFromBookshelf(userID, novelID, shelfType string) error {
	return s.repo.RemoveFromBookshelf(userID, novelID, shelfType)
}

func (s *readingService) GetBookshelf(userID, shelfType string, page, size int) ([]models.BookshelfResponse, int64, error) {
	items, total, err := s.repo.GetBookshelf(userID, shelfType, page, size)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.BookshelfResponse, len(items))
	for i, item := range items {
		responses[i] = *s.convertToBookshelfResponse(&item)
	}

	return responses, total, nil
}

func (s *readingService) GetBookshelfStats(userID string) (*models.BookshelfStatsResponse, error) {
	return s.repo.GetBookshelfStats(userID)
}

// Favorites
func (s *readingService) AddFavorite(userID, novelID string) error {
	favorite := &models.Favorite{
		UserID:    userID,
		NovelID:   novelID,
		CreatedAt: time.Now(),
	}

	// 同时添加到书架的收藏分类
	bookshelfReq := &models.AddToBookshelfRequest{
		NovelID:   novelID,
		ShelfType: "favorite",
	}

	if err := s.AddToBookshelf(userID, bookshelfReq); err != nil {
		return err
	}

	return s.repo.AddFavorite(favorite)
}

func (s *readingService) RemoveFavorite(userID, novelID string) error {
	// 同时从书架移除
	if err := s.RemoveFromBookshelf(userID, novelID, "favorite"); err != nil {
		return err
	}

	return s.repo.RemoveFavorite(userID, novelID)
}

func (s *readingService) GetUserFavorites(userID string, page, size int) ([]models.BookshelfResponse, int64, error) {
	return s.GetBookshelf(userID, "favorite", page, size)
}

func (s *readingService) CheckFavoriteStatus(userID, novelID string) bool {
	return s.repo.IsFavorite(userID, novelID)
}

// Comments
func (s *readingService) CreateComment(userID string, req *models.CreateCommentRequest) (*models.CommentResponse, error) {
	comment := &models.Comment{
		UserID:    userID,
		NovelID:   req.NovelID,
		ChapterID: req.ChapterID,
		Content:   req.Content,
		Rating:    req.Rating,
		ParentID:  req.ParentID,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateComment(comment); err != nil {
		return nil, err
	}

	return s.convertToCommentResponse(comment), nil
}

func (s *readingService) UpdateComment(userID, commentID string, req *models.UpdateCommentRequest) error {
	comment, err := s.repo.GetCommentByID(commentID)
	if err != nil {
		return err
	}

	if comment.UserID != userID {
		return errors.New("permission denied")
	}

	comment.Content = req.Content
	if req.Rating != nil {
		comment.Rating = req.Rating
	}
	comment.UpdatedAt = time.Now()

	return s.repo.UpdateComment(comment)
}

func (s *readingService) DeleteComment(userID, commentID string) error {
	comment, err := s.repo.GetCommentByID(commentID)
	if err != nil {
		return err
	}

	if comment.UserID != userID {
		return errors.New("permission denied")
	}

	return s.repo.DeleteComment(commentID)
}

func (s *readingService) GetNovelComments(novelID string, page, size int) ([]models.CommentResponse, int64, error) {
	comments, total, err := s.repo.GetNovelComments(novelID, page, size)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = *s.convertToCommentResponse(&comment)
		// 转换回复
		if len(comment.Replies) > 0 {
			replies := make([]models.CommentResponse, len(comment.Replies))
			for j, reply := range comment.Replies {
				replies[j] = *s.convertToCommentResponse(&reply)
			}
			responses[i].Replies = replies
		}
	}

	return responses, total, nil
}

func (s *readingService) GetChapterComments(chapterID string, page, size int) ([]models.CommentResponse, int64, error) {
	comments, total, err := s.repo.GetChapterComments(chapterID, page, size)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = *s.convertToCommentResponse(&comment)
		// 转换回复
		if len(comment.Replies) > 0 {
			replies := make([]models.CommentResponse, len(comment.Replies))
			for j, reply := range comment.Replies {
				replies[j] = *s.convertToCommentResponse(&reply)
			}
			responses[i].Replies = replies
		}
	}

	return responses, total, nil
}

func (s *readingService) GetUserComments(userID string, page, size int) ([]models.CommentResponse, int64, error) {
	comments, total, err := s.repo.GetUserComments(userID, page, size)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = *s.convertToCommentResponse(&comment)
	}

	return responses, total, nil
}

func (s *readingService) LikeComment(userID, commentID string) error {
	// 这里简化实现，实际应该记录用户点赞状态防止重复点赞
	return s.repo.UpdateCommentLikes(commentID, 1)
}

func (s *readingService) UnlikeComment(userID, commentID string) error {
	// 这里简化实现，实际应该检查用户点赞状态
	return s.repo.UpdateCommentLikes(commentID, -1)
}

// Search
func (s *readingService) AddSearchHistory(userID, keyword, searchType string, resultCount int) error {
	if searchType == "" {
		searchType = "novel"
	}

	history := &models.SearchHistory{
		UserID:      userID,
		Keyword:     keyword,
		SearchType:  searchType,
		ResultCount: resultCount,
		CreatedAt:   time.Now(),
	}

	return s.repo.AddSearchHistory(history)
}

func (s *readingService) GetSearchHistory(userID string) ([]models.SearchHistory, error) {
	return s.repo.GetUserSearchHistory(userID, 20)
}

func (s *readingService) ClearSearchHistory(userID string) error {
	return s.repo.ClearUserSearchHistory(userID)
}

// Statistics
func (s *readingService) GetReadingStats(userID string) (*models.ReadingStatsResponse, error) {
	return s.repo.GetReadingStats(userID)
}

// Helper methods
func (s *readingService) convertToReadingRecordResponse(record *models.ReadingRecord) *models.ReadingRecordResponse {
	return &models.ReadingRecordResponse{
		ID:              record.ID,
		NovelID:         record.NovelID,
		ChapterID:       record.ChapterID,
		ChapterNumber:   record.ChapterNumber,
		ReadingPosition: record.ReadingPosition,
		ReadingProgress: record.ReadingProgress,
		ReadingTime:     record.ReadingTime,
		LastReadAt:      record.LastReadAt.Format(time.RFC3339),
	}
}

func (s *readingService) convertToBookshelfResponse(item *models.Bookshelf) *models.BookshelfResponse {
	response := &models.BookshelfResponse{
		ID:              item.ID,
		NovelID:         item.NovelID,
		ShelfType:       item.ShelfType,
		AddedAt:         item.AddedAt.Format(time.RFC3339),
		ReadingProgress: item.ReadingProgress,
		IsArchived:      item.IsArchived,
	}

	if item.LastReadAt != nil {
		lastReadAt := item.LastReadAt.Format(time.RFC3339)
		response.LastReadAt = &lastReadAt
	}

	return response
}

func (s *readingService) convertToCommentResponse(comment *models.Comment) *models.CommentResponse {
	return &models.CommentResponse{
		ID:         comment.ID,
		UserID:     comment.UserID,
		NovelID:    comment.NovelID,
		ChapterID:  comment.ChapterID,
		Content:    comment.Content,
		Rating:     comment.Rating,
		LikeCount:  comment.LikeCount,
		ReplyCount: comment.ReplyCount,
		ParentID:   comment.ParentID,
		CreatedAt:  comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  comment.UpdatedAt.Format(time.RFC3339),
	}
}
