package models

import (
	"time"
	"gorm.io/gorm"
)

type Category struct {
	ID          string    `gorm:"type:varchar(36);primarykey" json:"id"`
	Name        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Description *string   `gorm:"type:text" json:"description"`
	Icon        *string   `gorm:"type:varchar(100)" json:"icon"`
	Color       *string   `gorm:"type:varchar(20)" json:"color"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`

	// 关联
	Novels []Novel `gorm:"foreignKey:CategoryID" json:"novels,omitempty"`
}

type Tag struct {
	ID         string    `gorm:"type:varchar(36);primarykey" json:"id"`
	Name       string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Color      *string   `gorm:"type:varchar(20)" json:"color"`
	UsageCount int       `gorm:"default:0" json:"usage_count"`
	CreatedAt  time.Time `json:"created_at"`

	// 关联
	NovelTags []NovelTag `gorm:"foreignKey:TagID" json:"novel_tags,omitempty"`
}

type Novel struct {
	ID             string     `gorm:"type:varchar(36);primarykey" json:"id"`
	Title          string     `gorm:"type:varchar(200);not null;index" json:"title"`
	Author         string     `gorm:"type:varchar(100);not null;index" json:"author"`
	CoverURL       *string    `gorm:"type:varchar(500)" json:"cover_url"`
	Description    *string    `gorm:"type:text" json:"description"`
	CategoryID     string     `gorm:"type:varchar(36);not null;index" json:"category_id"`
	Status         string     `gorm:"type:enum('ongoing','completed','paused');default:'ongoing';index" json:"status"`
	TotalChapters  int        `gorm:"default:0" json:"total_chapters"`
	WordCount      int64      `gorm:"default:0" json:"word_count"`
	ViewsCount     int64      `gorm:"default:0;index" json:"views_count"`
	Rating         float64    `gorm:"type:decimal(3,2);default:0.00;index" json:"rating"`
	RatingCount    int        `gorm:"default:0" json:"rating_count"`
	IsFeatured     bool       `gorm:"default:false" json:"is_featured"`
	IsFree         bool       `gorm:"default:true" json:"is_free"`
	Price          float64    `gorm:"type:decimal(10,2);default:0.00" json:"price"`
	PublishDate    *time.Time `gorm:"type:date" json:"publish_date"`
	LastUpdatedAt  *time.Time `json:"last_updated_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// 关联
	Category  Category    `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Chapters  []Chapter   `gorm:"foreignKey:NovelID" json:"chapters,omitempty"`
	NovelTags []NovelTag  `gorm:"foreignKey:NovelID" json:"novel_tags,omitempty"`
	Tags      []Tag       `gorm:"many2many:novel_tags;" json:"tags,omitempty"`
}

type NovelTag struct {
	ID        string    `gorm:"type:varchar(36);primarykey" json:"id"`
	NovelID   string    `gorm:"type:varchar(36);not null;uniqueIndex:idx_novel_tag" json:"novel_id"`
	TagID     string    `gorm:"type:varchar(36);not null;uniqueIndex:idx_novel_tag" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`

	// 关联
	Novel Novel `gorm:"foreignKey:NovelID" json:"novel,omitempty"`
	Tag   Tag   `gorm:"foreignKey:TagID" json:"tag,omitempty"`
}

type Chapter struct {
	ID            string     `gorm:"type:varchar(36);primarykey" json:"id"`
	NovelID       string     `gorm:"type:varchar(36);not null;index:idx_novel_chapter" json:"novel_id"`
	Title         string     `gorm:"type:varchar(200);not null" json:"title"`
	Content       string     `gorm:"type:longtext;not null" json:"content"`
	ChapterNumber int        `gorm:"not null;index:idx_novel_chapter" json:"chapter_number"`
	WordCount     int        `gorm:"default:0" json:"word_count"`
	IsFree        bool       `gorm:"default:true" json:"is_free"`
	Price         float64    `gorm:"type:decimal(10,2);default:0.00" json:"price"`
	PublishDate   *time.Time `gorm:"index" json:"publish_date"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// 关联
	Novel Novel `gorm:"foreignKey:NovelID" json:"novel,omitempty"`
}

// 钩子函数
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = generateUUID()
	}
	return nil
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = generateUUID()
	}
	return nil
}

func (n *Novel) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = generateUUID()
	}
	return nil
}

func (nt *NovelTag) BeforeCreate(tx *gorm.DB) error {
	if nt.ID == "" {
		nt.ID = generateUUID()
	}
	return nil
}

func (c *Chapter) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = generateUUID()
	}
	return nil
}