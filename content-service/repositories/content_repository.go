package repositories

import (
	"gorm.io/gorm"
	"reading-microservices/content-service/models"
	"strings"
)

type ContentRepository interface {
	// CreateCategory Category
	CreateCategory(category *models.Category) error
	GetCategoryByID(id string) (*models.Category, error)
	GetCategories(isActive *bool) ([]models.Category, error)
	UpdateCategory(category *models.Category) error
	DeleteCategory(id string) error

	// CreateTag Tag
	CreateTag(tag *models.Tag) error
	GetTagByID(id string) (*models.Tag, error)
	GetTagByName(name string) (*models.Tag, error)
	GetTags() ([]models.Tag, error)
	UpdateTag(tag *models.Tag) error
	DeleteTag(id string) error

	// CreateNovel Novel
	CreateNovel(novel *models.Novel) error
	GetNovelByID(id string) (*models.Novel, error)
	GetNovelsByCategory(categoryID string, page, size int) ([]models.Novel, int64, error)
	SearchNovels(params *models.NovelSearchParams) ([]models.Novel, int64, error)
	UpdateNovel(novel *models.Novel) error
	DeleteNovel(id string) error
	UpdateNovelStats(novelID string, views *int64, rating *float64, ratingCount *int) error
	GetFeaturedNovels(limit int) ([]models.Novel, error)
	GetLatestNovels(limit int) ([]models.Novel, error)

	// CreateChapter Chapter
	CreateChapter(chapter *models.Chapter) error
	GetChapterByID(id string) (*models.Chapter, error)
	GetChaptersByNovel(novelID string, page, size int) ([]models.Chapter, int64, error)
	GetChapterByNumber(novelID string, chapterNumber int) (*models.Chapter, error)
	UpdateChapter(chapter *models.Chapter) error
	DeleteChapter(id string) error
	GetLatestChapters(novelID string, limit int) ([]models.Chapter, error)

	// AddNovelTags NovelTag
	AddNovelTags(novelID string, tagIDs []string) error
	RemoveNovelTags(novelID string, tagIDs []string) error
	GetNovelTags(novelID string) ([]models.Tag, error)
}

type contentRepository struct {
	db *gorm.DB
}

func NewContentRepository(db *gorm.DB) ContentRepository {
	return &contentRepository{db: db}
}

// Category methods
func (r *contentRepository) CreateCategory(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *contentRepository) GetCategoryByID(id string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *contentRepository) GetCategories(isActive *bool) ([]models.Category, error) {
	var categories []models.Category
	query := r.db.Order("sort_order ASC, created_at ASC")

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Find(&categories).Error
	return categories, err
}

func (r *contentRepository) UpdateCategory(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *contentRepository) DeleteCategory(id string) error {
	return r.db.Delete(&models.Category{}, "id = ?", id).Error
}

// Tag methods
func (r *contentRepository) CreateTag(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

func (r *contentRepository) GetTagByID(id string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("id = ?", id).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *contentRepository) GetTagByName(name string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *contentRepository) GetTags() ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Order("usage_count DESC, created_at ASC").Find(&tags).Error
	return tags, err
}

func (r *contentRepository) UpdateTag(tag *models.Tag) error {
	return r.db.Save(tag).Error
}

func (r *contentRepository) DeleteTag(id string) error {
	return r.db.Delete(&models.Tag{}, "id = ?", id).Error
}

// Novel methods
func (r *contentRepository) CreateNovel(novel *models.Novel) error {
	return r.db.Create(novel).Error
}

func (r *contentRepository) GetNovelByID(id string) (*models.Novel, error) {
	var novel models.Novel
	err := r.db.Preload("Category").Preload("Tags").Where("id = ?", id).First(&novel).Error
	if err != nil {
		return nil, err
	}
	return &novel, nil
}

func (r *contentRepository) GetNovelsByCategory(categoryID string, page, size int) ([]models.Novel, int64, error) {
	var novels []models.Novel
	var total int64

	query := r.db.Model(&models.Novel{}).Where("category_id = ?", categoryID)

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * size
	err := query.Preload("Category").Preload("Tags").
		Offset(offset).Limit(size).
		Order("created_at DESC").Find(&novels).Error

	return novels, total, err
}

func (r *contentRepository) SearchNovels(params *models.NovelSearchParams) ([]models.Novel, int64, error) {
	var novels []models.Novel
	var total int64

	query := r.db.Model(&models.Novel{})

	// 关键词搜索
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("title LIKE ? OR author LIKE ? OR description LIKE ?", keyword, keyword, keyword)
	}

	// 分类筛选
	if params.CategoryID != "" {
		query = query.Where("category_id = ?", params.CategoryID)
	}

	// 标签筛选
	if params.TagIDs != "" {
		tagIDs := strings.Split(params.TagIDs, ",")
		query = query.Joins("JOIN novel_tags ON novels.id = novel_tags.novel_id").
			Where("novel_tags.tag_id IN ?", tagIDs)
	}

	// 状态筛选
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// 免费筛选
	if params.IsFree != "" {
		query = query.Where("is_free = ?", params.IsFree == "true")
	}

	// 获取总数
	query.Count(&total)

	// 排序
	orderBy := "created_at DESC"
	switch params.OrderBy {
	case "views":
		orderBy = "views_count DESC"
	case "rating":
		orderBy = "rating DESC"
	case "updated":
		orderBy = "last_updated_at DESC"
	case "chapters":
		orderBy = "total_chapters DESC"
	}

	// 分页
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Size <= 0 {
		params.Size = 20
	}

	offset := (params.Page - 1) * params.Size
	err := query.Preload("Category").Preload("Tags").
		Offset(offset).Limit(params.Size).
		Order(orderBy).Find(&novels).Error

	return novels, total, err
}

func (r *contentRepository) UpdateNovel(novel *models.Novel) error {
	return r.db.Save(novel).Error
}

func (r *contentRepository) DeleteNovel(id string) error {
	return r.db.Delete(&models.Novel{}, "id = ?", id).Error
}

func (r *contentRepository) UpdateNovelStats(novelID string, views *int64, rating *float64, ratingCount *int) error {
	updates := make(map[string]interface{})

	if views != nil {
		updates["views_count"] = gorm.Expr("views_count + ?", *views)
	}
	if rating != nil {
		updates["rating"] = *rating
	}
	if ratingCount != nil {
		updates["rating_count"] = *ratingCount
	}

	return r.db.Model(&models.Novel{}).Where("id = ?", novelID).Updates(updates).Error
}

func (r *contentRepository) GetFeaturedNovels(limit int) ([]models.Novel, error) {
	var novels []models.Novel
	err := r.db.Where("is_featured = ?", true).
		Preload("Category").Preload("Tags").
		Order("views_count DESC").Limit(limit).Find(&novels).Error
	return novels, err
}

func (r *contentRepository) GetLatestNovels(limit int) ([]models.Novel, error) {
	var novels []models.Novel
	err := r.db.Preload("Category").Preload("Tags").
		Order("created_at DESC").Limit(limit).Find(&novels).Error
	return novels, err
}

// Chapter methods
func (r *contentRepository) CreateChapter(chapter *models.Chapter) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 创建章节
		if err := tx.Create(chapter).Error; err != nil {
			return err
		}

		// 更新小说统计
		return tx.Model(&models.Novel{}).
			Where("id = ?", chapter.NovelID).
			Updates(map[string]interface{}{
				"total_chapters":  gorm.Expr("total_chapters + 1"),
				"word_count":      gorm.Expr("word_count + ?", chapter.WordCount),
				"last_updated_at": gorm.Expr("NOW()"),
			}).Error
	})
}

func (r *contentRepository) GetChapterByID(id string) (*models.Chapter, error) {
	var chapter models.Chapter
	err := r.db.Where("id = ?", id).First(&chapter).Error
	if err != nil {
		return nil, err
	}
	return &chapter, nil
}

func (r *contentRepository) GetChaptersByNovel(novelID string, page, size int) ([]models.Chapter, int64, error) {
	var chapters []models.Chapter
	var total int64

	query := r.db.Model(&models.Chapter{}).Where("novel_id = ?", novelID)

	// 获取总数
	query.Count(&total)

	// 分页查询
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 50
	}

	offset := (page - 1) * size
	err := query.Select("id, novel_id, title, chapter_number, word_count, is_free, price, created_at").
		Offset(offset).Limit(size).
		Order("chapter_number ASC").Find(&chapters).Error

	return chapters, total, err
}

func (r *contentRepository) GetChapterByNumber(novelID string, chapterNumber int) (*models.Chapter, error) {
	var chapter models.Chapter
	err := r.db.Where("novel_id = ? AND chapter_number = ?", novelID, chapterNumber).First(&chapter).Error
	if err != nil {
		return nil, err
	}
	return &chapter, nil
}

func (r *contentRepository) UpdateChapter(chapter *models.Chapter) error {
	return r.db.Save(chapter).Error
}

func (r *contentRepository) DeleteChapter(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 获取章节信息
		var chapter models.Chapter
		if err := tx.Where("id = ?", id).First(&chapter).Error; err != nil {
			return err
		}

		// 删除章节
		if err := tx.Delete(&chapter).Error; err != nil {
			return err
		}

		// 更新小说统计
		return tx.Model(&models.Novel{}).
			Where("id = ?", chapter.NovelID).
			Updates(map[string]interface{}{
				"total_chapters": gorm.Expr("total_chapters - 1"),
				"word_count":     gorm.Expr("word_count - ?", chapter.WordCount),
			}).Error
	})
}

func (r *contentRepository) GetLatestChapters(novelID string, limit int) ([]models.Chapter, error) {
	var chapters []models.Chapter
	err := r.db.Where("novel_id = ?", novelID).
		Select("id, novel_id, title, chapter_number, word_count, is_free, price, created_at").
		Order("chapter_number DESC").Limit(limit).Find(&chapters).Error
	return chapters, err
}

// NovelTag methods
func (r *contentRepository) AddNovelTags(novelID string, tagIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, tagID := range tagIDs {
			novelTag := &models.NovelTag{
				NovelID: novelID,
				TagID:   tagID,
			}

			// 检查是否已存在
			var count int64
			tx.Model(&models.NovelTag{}).Where("novel_id = ? AND tag_id = ?", novelID, tagID).Count(&count)
			if count == 0 {
				if err := tx.Create(novelTag).Error; err != nil {
					return err
				}

				// 更新标签使用次数
				tx.Model(&models.Tag{}).Where("id = ?", tagID).Update("usage_count", gorm.Expr("usage_count + 1"))
			}
		}
		return nil
	})
}

func (r *contentRepository) RemoveNovelTags(novelID string, tagIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("novel_id = ? AND tag_id IN ?", novelID, tagIDs).Delete(&models.NovelTag{}).Error; err != nil {
			return err
		}

		// 更新标签使用次数
		for _, tagID := range tagIDs {
			tx.Model(&models.Tag{}).Where("id = ?", tagID).Update("usage_count", gorm.Expr("usage_count - 1"))
		}

		return nil
	})
}

func (r *contentRepository) GetNovelTags(novelID string) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Joins("JOIN novel_tags ON tags.id = novel_tags.tag_id").
		Where("novel_tags.novel_id = ?", novelID).Find(&tags).Error
	return tags, err
}
