package models

import (
	"gorm.io/gorm"
	"time"
)

type ReadingRecord struct {
	ID              string    `gorm:"type:varchar(36);primarykey" json:"id"`
	UserID          string    `gorm:"type:varchar(36);not null;index:idx_user_novel" json:"user_id"`
	NovelID         string    `gorm:"type:varchar(36);not null;index:idx_user_novel" json:"novel_id"`
	ChapterID       string    `gorm:"type:varchar(36);not null" json:"chapter_id"`
	ChapterNumber   int       `gorm:"not null" json:"chapter_number"`
	ReadingPosition int       `gorm:"default:0" json:"reading_position"`
	ReadingProgress float64   `gorm:"type:decimal(5,2);default:0.00" json:"reading_progress"`
	ReadingTime     int       `gorm:"default:0" json:"reading_time"` // 阅读时长(秒)
	LastReadAt      time.Time `gorm:"type:datetime(3);autoCreateTime;autoUpdateTime" json:"last_read_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Bookshelf struct {
	ID              string     `gorm:"type:varchar(36);primarykey" json:"id"`
	UserID          string     `gorm:"type:varchar(36);not null;index" json:"user_id"`
	NovelID         string     `gorm:"type:varchar(36);not null" json:"novel_id"`
	ShelfType       string     `gorm:"type:enum('reading','favorite','download');default:'reading';index" json:"shelf_type"`
	AddedAt         time.Time  `gorm:"type:datetime(3);autoCreateTime;autoUpdateTime" json:"added_at"`
	LastReadAt      *time.Time `json:"last_read_at"`
	ReadingProgress float64    `gorm:"type:decimal(5,2);default:0.00" json:"reading_progress"`
	IsArchived      bool       `gorm:"default:false" json:"is_archived"`
}

type Favorite struct {
	ID        string    `gorm:"type:varchar(36);primarykey" json:"id"`
	UserID    string    `gorm:"type:varchar(36);not null;index" json:"user_id"`
	NovelID   string    `gorm:"type:varchar(36);not null" json:"novel_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID         string    `gorm:"type:varchar(36);primarykey" json:"id"`
	UserID     string    `gorm:"type:varchar(36);not null;index" json:"user_id"`
	NovelID    string    `gorm:"type:varchar(36);not null;index" json:"novel_id"`
	ChapterID  *string   `gorm:"type:varchar(36);index" json:"chapter_id"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	Rating     *int      `gorm:"check:rating >= 1 AND rating <= 5" json:"rating"`
	LikeCount  int       `gorm:"default:0" json:"like_count"`
	ReplyCount int       `gorm:"default:0" json:"reply_count"`
	ParentID   *string   `gorm:"type:varchar(36);index" json:"parent_id"`
	IsDeleted  bool      `gorm:"default:false" json:"is_deleted"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// 关联
	Replies []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

type SearchHistory struct {
	ID          string    `gorm:"type:varchar(36);primarykey" json:"id"`
	UserID      string    `gorm:"type:varchar(36);not null;index" json:"user_id"`
	Keyword     string    `gorm:"type:varchar(200);not null;index" json:"keyword"`
	SearchType  string    `gorm:"type:enum('novel','author','tag');default:'novel'" json:"search_type"`
	ResultCount int       `gorm:"default:0" json:"result_count"`
	CreatedAt   time.Time `json:"created_at"`
}

// Hook functions
func (r *ReadingRecord) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = generateUUID()
	}
	return nil
}

func (b *Bookshelf) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = generateUUID()
	}
	return nil
}

func (f *Favorite) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = generateUUID()
	}
	return nil
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = generateUUID()
	}
	return nil
}

func (s *SearchHistory) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = generateUUID()
	}
	return nil
}
