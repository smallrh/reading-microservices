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
type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required,max=50"`
	Description *string `json:"description"`
	Icon        *string `json:"icon"`
	Color       *string `json:"color"`
	SortOrder   *int    `json:"sort_order"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=50"`
	Description *string `json:"description"`
	Icon        *string `json:"icon"`
	Color       *string `json:"color"`
	SortOrder   *int    `json:"sort_order"`
	IsActive    *bool   `json:"is_active"`
}

type CreateTagRequest struct {
	Name  string  `json:"name" binding:"required,max=50"`
	Color *string `json:"color"`
}

type CreateNovelRequest struct {
	Title       string  `json:"title" binding:"required,max=200"`
	Author      string  `json:"author" binding:"required,max=100"`
	CoverURL    *string `json:"cover_url"`
	Description *string `json:"description"`
	CategoryID  string  `json:"category_id" binding:"required"`
	Status      *string `json:"status" binding:"omitempty,oneof=ongoing completed paused"`
	IsFeatured  *bool   `json:"is_featured"`
	IsFree      *bool   `json:"is_free"`
	Price       *float64 `json:"price"`
	TagIDs      []string `json:"tag_ids"`
}

type UpdateNovelRequest struct {
	Title       *string  `json:"title" binding:"omitempty,max=200"`
	Author      *string  `json:"author" binding:"omitempty,max=100"`
	CoverURL    *string  `json:"cover_url"`
	Description *string  `json:"description"`
	CategoryID  *string  `json:"category_id"`
	Status      *string  `json:"status" binding:"omitempty,oneof=ongoing completed paused"`
	IsFeatured  *bool    `json:"is_featured"`
	IsFree      *bool    `json:"is_free"`
	Price       *float64 `json:"price"`
	TagIDs      []string `json:"tag_ids"`
}

type CreateChapterRequest struct {
	NovelID       string  `json:"novel_id" binding:"required"`
	Title         string  `json:"title" binding:"required,max=200"`
	Content       string  `json:"content" binding:"required"`
	ChapterNumber int     `json:"chapter_number" binding:"required,min=1"`
	IsFree        *bool   `json:"is_free"`
	Price         *float64 `json:"price"`
}

type UpdateChapterRequest struct {
	Title   *string  `json:"title" binding:"omitempty,max=200"`
	Content *string  `json:"content"`
	IsFree  *bool    `json:"is_free"`
	Price   *float64 `json:"price"`
}

// 查询参数
type NovelSearchParams struct {
	Keyword    string `form:"keyword"`
	CategoryID string `form:"category_id"`
	TagIDs     string `form:"tag_ids"`
	Status     string `form:"status"`
	IsFree     string `form:"is_free"`
	OrderBy    string `form:"order_by"`
	Page       int    `form:"page"`
	Size       int    `form:"size"`
}

type ChapterListParams struct {
	NovelID string `form:"novel_id" binding:"required"`
	Page    int    `form:"page"`
	Size    int    `form:"size"`
}

// 响应DTO
type NovelListResponse struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Author        string   `json:"author"`
	CoverURL      *string  `json:"cover_url"`
	Description   *string  `json:"description"`
	Category      Category `json:"category"`
	Status        string   `json:"status"`
	TotalChapters int      `json:"total_chapters"`
	WordCount     int64    `json:"word_count"`
	ViewsCount    int64    `json:"views_count"`
	Rating        float64  `json:"rating"`
	RatingCount   int      `json:"rating_count"`
	IsFeatured    bool     `json:"is_featured"`
	IsFree        bool     `json:"is_free"`
	Price         float64  `json:"price"`
	Tags          []Tag    `json:"tags"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}

type NovelDetailResponse struct {
	*NovelListResponse
	LatestChapters []ChapterSummary `json:"latest_chapters"`
}

type ChapterSummary struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	ChapterNumber int     `json:"chapter_number"`
	WordCount     int     `json:"word_count"`
	IsFree        bool    `json:"is_free"`
	Price         float64 `json:"price"`
	CreatedAt     string  `json:"created_at"`
}

type ChapterDetailResponse struct {
	ID            string  `json:"id"`
	NovelID       string  `json:"novel_id"`
	Title         string  `json:"title"`
	Content       string  `json:"content"`
	ChapterNumber int     `json:"chapter_number"`
	WordCount     int     `json:"word_count"`
	IsFree        bool    `json:"is_free"`
	Price         float64 `json:"price"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	PrevChapter   *ChapterNav `json:"prev_chapter"`
	NextChapter   *ChapterNav `json:"next_chapter"`
}

type ChapterNav struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	ChapterNumber int    `json:"chapter_number"`
}