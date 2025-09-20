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
type UpdateReadingProgressRequest struct {
	NovelID         string  `json:"novel_id" binding:"required"`
	ChapterID       string  `json:"chapter_id" binding:"required"`
	ChapterNumber   int     `json:"chapter_number" binding:"required"`
	ReadingPosition int     `json:"reading_position"`
	ReadingProgress float64 `json:"reading_progress"`
	ReadingTime     int     `json:"reading_time"`
}

type AddToBookshelfRequest struct {
	NovelID   string `json:"novel_id" binding:"required"`
	ShelfType string `json:"shelf_type" binding:"required,oneof=reading favorite download"`
}

type CreateCommentRequest struct {
	NovelID   string  `json:"novel_id" binding:"required"`
	ChapterID *string `json:"chapter_id"`
	Content   string  `json:"content" binding:"required,max=2000"`
	Rating    *int    `json:"rating" binding:"omitempty,min=1,max=5"`
	ParentID  *string `json:"parent_id"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,max=2000"`
	Rating  *int   `json:"rating" binding:"omitempty,min=1,max=5"`
}

type SearchRequest struct {
	Keyword    string `form:"keyword" binding:"required"`
	SearchType string `form:"search_type" binding:"omitempty,oneof=novel author tag"`
	Page       int    `form:"page"`
	Size       int    `form:"size"`
}

// 响应DTO
type ReadingRecordResponse struct {
	ID              string `json:"id"`
	NovelID         string `json:"novel_id"`
	ChapterID       string `json:"chapter_id"`
	ChapterNumber   int    `json:"chapter_number"`
	ReadingPosition int    `json:"reading_position"`
	ReadingProgress float64 `json:"reading_progress"`
	ReadingTime     int    `json:"reading_time"`
	LastReadAt      string `json:"last_read_at"`
	Novel           *NovelInfo `json:"novel,omitempty"`
}

type BookshelfResponse struct {
	ID              string     `json:"id"`
	NovelID         string     `json:"novel_id"`
	ShelfType       string     `json:"shelf_type"`
	AddedAt         string     `json:"added_at"`
	LastReadAt      *string    `json:"last_read_at"`
	ReadingProgress float64    `json:"reading_progress"`
	IsArchived      bool       `json:"is_archived"`
	Novel           *NovelInfo `json:"novel,omitempty"`
}

type CommentResponse struct {
	ID         string             `json:"id"`
	UserID     string             `json:"user_id"`
	NovelID    string             `json:"novel_id"`
	ChapterID  *string            `json:"chapter_id"`
	Content    string             `json:"content"`
	Rating     *int               `json:"rating"`
	LikeCount  int                `json:"like_count"`
	ReplyCount int                `json:"reply_count"`
	ParentID   *string            `json:"parent_id"`
	CreatedAt  string             `json:"created_at"`
	UpdatedAt  string             `json:"updated_at"`
	User       *UserInfo          `json:"user,omitempty"`
	Replies    []CommentResponse  `json:"replies,omitempty"`
}

type NovelInfo struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	Author        string  `json:"author"`
	CoverURL      *string `json:"cover_url"`
	Status        string  `json:"status"`
	TotalChapters int     `json:"total_chapters"`
	Rating        float64 `json:"rating"`
}

type UserInfo struct {
	ID        string  `json:"id"`
	Username  string  `json:"username"`
	Nickname  *string `json:"nickname"`
	AvatarURL *string `json:"avatar_url"`
}

type ReadingStatsResponse struct {
	TotalReadingTime    int     `json:"total_reading_time"`    // 总阅读时长(分钟)
	TotalBooksRead      int     `json:"total_books_read"`      // 已读书籍数
	TotalChaptersRead   int     `json:"total_chapters_read"`   // 已读章节数
	AverageReadingSpeed int     `json:"average_reading_speed"` // 平均阅读速度(字/分钟)
	ReadingStreak       int     `json:"reading_streak"`        // 连续阅读天数
	FavoriteCategories  []string `json:"favorite_categories"`  // 最喜欢的分类
}

type BookshelfStatsResponse struct {
	Reading   int `json:"reading"`   // 在读数量
	Favorite  int `json:"favorite"`  // 收藏数量
	Download  int `json:"download"`  // 下载数量
	Archived  int `json:"archived"`  // 归档数量
}