package repositories

import (
	"gorm.io/gorm"
	"reading-microservices/reading-service/models"
	"time"
)

type ReadingRepository interface {
	// Reading Records
	CreateOrUpdateReadingRecord(record *models.ReadingRecord) error
	GetReadingRecord(userID, novelID, chapterID string) (*models.ReadingRecord, error)
	GetUserReadingRecords(userID string, page, size int) ([]models.ReadingRecord, int64, error)
	GetLatestReadingRecords(userID string, limit int) ([]models.ReadingRecord, error)

	// Bookshelf
	AddToBookshelf(bookshelf *models.Bookshelf) error
	RemoveFromBookshelf(userID, novelID, shelfType string) error
	GetBookshelf(userID, shelfType string, page, size int) ([]models.Bookshelf, int64, error)
	GetBookshelfItem(userID, novelID, shelfType string) (*models.Bookshelf, error)
	UpdateBookshelfProgress(userID, novelID string, progress float64) error
	GetBookshelfStats(userID string) (*models.BookshelfStatsResponse, error)

	// Favorites
	AddFavorite(favorite *models.Favorite) error
	RemoveFavorite(userID, novelID string) error
	GetUserFavorites(userID string, page, size int) ([]models.Favorite, int64, error)
	IsFavorite(userID, novelID string) bool

	// Comments
	CreateComment(comment *models.Comment) error
	GetCommentByID(id string) (*models.Comment, error)
	UpdateComment(comment *models.Comment) error
	DeleteComment(id string) error
	GetNovelComments(novelID string, page, size int) ([]models.Comment, int64, error)
	GetChapterComments(chapterID string, page, size int) ([]models.Comment, int64, error)
	GetUserComments(userID string, page, size int) ([]models.Comment, int64, error)
	UpdateCommentLikes(commentID string, increment int) error

	// Search History
	AddSearchHistory(history *models.SearchHistory) error
	GetUserSearchHistory(userID string, limit int) ([]models.SearchHistory, error)
	ClearUserSearchHistory(userID string) error

	// Statistics
	GetReadingStats(userID string) (*models.ReadingStatsResponse, error)
}

type readingRepository struct {
	db *gorm.DB
}

func NewReadingRepository(db *gorm.DB) ReadingRepository {
	return &readingRepository{db: db}
}

// Reading Records
func (r *readingRepository) CreateOrUpdateReadingRecord(record *models.ReadingRecord) error {
	var existing models.ReadingRecord
	err := r.db.Where("user_id = ? AND novel_id = ? AND chapter_id = ?",
		record.UserID, record.NovelID, record.ChapterID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新记录
		return r.db.Create(record).Error
	} else if err != nil {
		return err
	} else {
		// 更新现有记录
		existing.ReadingPosition = record.ReadingPosition
		existing.ReadingProgress = record.ReadingProgress
		existing.ReadingTime += record.ReadingTime
		existing.LastReadAt = record.LastReadAt
		return r.db.Save(&existing).Error
	}
}

func (r *readingRepository) GetReadingRecord(userID, novelID, chapterID string) (*models.ReadingRecord, error) {
	var record models.ReadingRecord
	err := r.db.Where("user_id = ? AND novel_id = ? AND chapter_id = ?",
		userID, novelID, chapterID).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *readingRepository) GetUserReadingRecords(userID string, page, size int) ([]models.ReadingRecord, int64, error) {
	var records []models.ReadingRecord
	var total int64

	query := r.db.Model(&models.ReadingRecord{}).Where("user_id = ?", userID)

	query.Count(&total)

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size

	err := query.Order("last_read_at DESC").Offset(offset).Limit(size).Find(&records).Error
	return records, total, err
}

func (r *readingRepository) GetLatestReadingRecords(userID string, limit int) ([]models.ReadingRecord, error) {
	var records []models.ReadingRecord
	err := r.db.Where("user_id = ?", userID).
		Order("last_read_at DESC").
		Limit(limit).
		Find(&records).Error
	return records, err
}

// Bookshelf
func (r *readingRepository) AddToBookshelf(bookshelf *models.Bookshelf) error {
	// 检查是否已存在
	var existing models.Bookshelf
	err := r.db.Where("user_id = ? AND novel_id = ? AND shelf_type = ?",
		bookshelf.UserID, bookshelf.NovelID, bookshelf.ShelfType).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		return r.db.Create(bookshelf).Error
	}

	return nil // 已存在，不需要重复添加
}

func (r *readingRepository) RemoveFromBookshelf(userID, novelID, shelfType string) error {
	return r.db.Where("user_id = ? AND novel_id = ? AND shelf_type = ?",
		userID, novelID, shelfType).Delete(&models.Bookshelf{}).Error
}

func (r *readingRepository) GetBookshelf(userID, shelfType string, page, size int) ([]models.Bookshelf, int64, error) {
	var items []models.Bookshelf
	var total int64

	query := r.db.Model(&models.Bookshelf{}).Where("user_id = ?", userID)

	if shelfType != "" {
		query = query.Where("shelf_type = ?", shelfType)
	}

	query.Count(&total)

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size

	err := query.Order("added_at DESC").Offset(offset).Limit(size).Find(&items).Error
	return items, total, err
}

func (r *readingRepository) GetBookshelfItem(userID, novelID, shelfType string) (*models.Bookshelf, error) {
	var item models.Bookshelf
	err := r.db.Where("user_id = ? AND novel_id = ? AND shelf_type = ?",
		userID, novelID, shelfType).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *readingRepository) UpdateBookshelfProgress(userID, novelID string, progress float64) error {
	now := time.Now()
	return r.db.Model(&models.Bookshelf{}).
		Where("user_id = ? AND novel_id = ? AND shelf_type = 'reading'", userID, novelID).
		Updates(map[string]interface{}{
			"reading_progress": progress,
			"last_read_at":     &now,
		}).Error
}

func (r *readingRepository) GetBookshelfStats(userID string) (*models.BookshelfStatsResponse, error) {
	var stats models.BookshelfStatsResponse
	var reading, favorite, download, archived int64

	// 统计各类型数量
	r.db.Model(&models.Bookshelf{}).Where("user_id = ? AND shelf_type = 'reading'", userID).Count(&reading)
	r.db.Model(&models.Bookshelf{}).Where("user_id = ? AND shelf_type = 'favorite'", userID).Count(&favorite)
	r.db.Model(&models.Bookshelf{}).Where("user_id = ? AND shelf_type = 'download'", userID).Count(&download)
	r.db.Model(&models.Bookshelf{}).Where("user_id = ? AND is_archived = true", userID).Count(&archived)

	// 转换为 int 类型
	stats.Reading = int(reading)
	stats.Favorite = int(favorite)
	stats.Download = int(download)
	stats.Archived = int(archived)

	return &stats, nil
}

// Favorites
func (r *readingRepository) AddFavorite(favorite *models.Favorite) error {
	// 检查是否已收藏
	var existing models.Favorite
	err := r.db.Where("user_id = ? AND novel_id = ?", favorite.UserID, favorite.NovelID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		return r.db.Create(favorite).Error
	}

	return nil // 已收藏，不需要重复添加
}

func (r *readingRepository) RemoveFavorite(userID, novelID string) error {
	return r.db.Where("user_id = ? AND novel_id = ?", userID, novelID).Delete(&models.Favorite{}).Error
}

func (r *readingRepository) GetUserFavorites(userID string, page, size int) ([]models.Favorite, int64, error) {
	var favorites []models.Favorite
	var total int64

	query := r.db.Model(&models.Favorite{}).Where("user_id = ?", userID)

	query.Count(&total)

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size

	err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&favorites).Error
	return favorites, total, err
}

func (r *readingRepository) IsFavorite(userID, novelID string) bool {
	var count int64
	r.db.Model(&models.Favorite{}).Where("user_id = ? AND novel_id = ?", userID, novelID).Count(&count)
	return count > 0
}

// Comments
func (r *readingRepository) CreateComment(comment *models.Comment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 创建评论
		if err := tx.Create(comment).Error; err != nil {
			return err
		}

		// 更新父评论的回复数
		if comment.ParentID != nil {
			if err := tx.Model(&models.Comment{}).Where("id = ?", *comment.ParentID).
				UpdateColumn("reply_count", gorm.Expr("reply_count + 1")).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *readingRepository) GetCommentByID(id string) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.Where("id = ? AND is_deleted = false", id).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *readingRepository) UpdateComment(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

func (r *readingRepository) DeleteComment(id string) error {
	return r.db.Model(&models.Comment{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *readingRepository) GetNovelComments(novelID string, page, size int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	query := r.db.Model(&models.Comment{}).Where("novel_id = ? AND parent_id IS NULL AND is_deleted = false", novelID)

	query.Count(&total)

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size

	err := query.Preload("Replies", "is_deleted = false").
		Order("created_at DESC").
		Offset(offset).Limit(size).
		Find(&comments).Error

	return comments, total, err
}

func (r *readingRepository) GetChapterComments(chapterID string, page, size int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	query := r.db.Model(&models.Comment{}).Where("chapter_id = ? AND parent_id IS NULL AND is_deleted = false", chapterID)

	query.Count(&total)

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size

	err := query.Preload("Replies", "is_deleted = false").
		Order("created_at DESC").
		Offset(offset).Limit(size).
		Find(&comments).Error

	return comments, total, err
}

func (r *readingRepository) GetUserComments(userID string, page, size int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	query := r.db.Model(&models.Comment{}).Where("user_id = ? AND is_deleted = false", userID)

	query.Count(&total)

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	offset := (page - 1) * size

	err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&comments).Error
	return comments, total, err
}

func (r *readingRepository) UpdateCommentLikes(commentID string, increment int) error {
	return r.db.Model(&models.Comment{}).Where("id = ?", commentID).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", increment)).Error
}

// Search History
func (r *readingRepository) AddSearchHistory(history *models.SearchHistory) error {
	// 检查是否已存在相同关键词
	var existing models.SearchHistory
	err := r.db.Where("user_id = ? AND keyword = ? AND search_type = ?",
		history.UserID, history.Keyword, history.SearchType).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		return r.db.Create(history).Error
	} else if err != nil {
		return err
	} else {
		// 更新时间和结果数量
		existing.ResultCount = history.ResultCount
		existing.CreatedAt = history.CreatedAt
		return r.db.Save(&existing).Error
	}
}

func (r *readingRepository) GetUserSearchHistory(userID string, limit int) ([]models.SearchHistory, error) {
	var history []models.SearchHistory
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&history).Error
	return history, err
}

func (r *readingRepository) ClearUserSearchHistory(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.SearchHistory{}).Error
}

// Statistics
func (r *readingRepository) GetReadingStats(userID string) (*models.ReadingStatsResponse, error) {
	stats := &models.ReadingStatsResponse{}

	// 总阅读时长(转换为分钟)
	var totalSeconds int64
	r.db.Model(&models.ReadingRecord{}).Where("user_id = ?", userID).
		Select("COALESCE(SUM(reading_time), 0)").Scan(&totalSeconds)
	stats.TotalReadingTime = int(totalSeconds / 60)

	// 已读书籍数(去重)
	var totalBooks int64
	r.db.Model(&models.ReadingRecord{}).Where("user_id = ?", userID).
		Select("COUNT(DISTINCT novel_id)").Scan(&totalBooks)
	stats.TotalBooksRead = int(totalBooks)

	// 已读章节数
	var totalChapters int64
	r.db.Model(&models.ReadingRecord{}).Where("user_id = ?", userID).Count(&totalChapters)
	stats.TotalChaptersRead = int(totalChapters)

	// 平均阅读速度(假设每千字需要3分钟阅读)
	if stats.TotalReadingTime > 0 {
		stats.AverageReadingSpeed = 300 // 字/分钟，这里是估算值
	}

	// 连续阅读天数(简化实现)
	stats.ReadingStreak = r.calculateReadingStreak(userID)

	return stats, nil
}

func (r *readingRepository) calculateReadingStreak(userID string) int {
	// 简化的连续阅读天数计算
	// 实际实现应该基于每日阅读记录
	var recentDays []string
	r.db.Model(&models.ReadingRecord{}).
		Where("user_id = ? AND last_read_at >= ?", userID, time.Now().AddDate(0, 0, -30)).
		Select("DISTINCT DATE(last_read_at) as date").
		Order("date DESC").
		Limit(30).
		Pluck("date", &recentDays)

	// 计算连续天数
	streak := 0
	today := time.Now().Format("2006-01-02")

	for i, date := range recentDays {
		expectedDate := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		if date == expectedDate || (i == 0 && date == today) {
			streak++
		} else {
			break
		}
	}

	return streak
}
