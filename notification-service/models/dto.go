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
type CreateNotificationRequest struct {
	UserID      string  `json:"user_id" binding:"required"`
	Title       string  `json:"title" binding:"required,max=200"`
	Content     string  `json:"content" binding:"required"`
	Type        string  `json:"type" binding:"required,oneof=update gift like recommendation system"`
	RelatedID   *string `json:"related_id"`
	RelatedType *string `json:"related_type"`
	ActionURL   *string `json:"action_url"`
}

type BatchNotificationRequest struct {
	UserIDs     []string `json:"user_ids" binding:"required"`
	Title       string   `json:"title" binding:"required,max=200"`
	Content     string   `json:"content" binding:"required"`
	Type        string   `json:"type" binding:"required,oneof=update gift like recommendation system"`
	RelatedID   *string  `json:"related_id"`
	RelatedType *string  `json:"related_type"`
	ActionURL   *string  `json:"action_url"`
}

type UpdateNotificationSettingRequest struct {
	SettingType string `json:"setting_type" binding:"required"`
	IsEnabled   bool   `json:"is_enabled"`
}

type RegisterPushTokenRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
	Token    string `json:"token" binding:"required"`
	Platform string `json:"platform" binding:"required,oneof=ios android web"`
}

type PushNotificationRequest struct {
	UserIDs []string `json:"user_ids"`
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Data    map[string]interface{} `json:"data"`
}

// 响应DTO
type NotificationResponse struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Content     string  `json:"content"`
	Type        string  `json:"type"`
	IsRead      bool    `json:"is_read"`
	RelatedID   *string `json:"related_id"`
	RelatedType *string `json:"related_type"`
	ActionURL   *string `json:"action_url"`
	CreatedAt   string  `json:"created_at"`
	ReadAt      *string `json:"read_at"`
}

type NotificationStatsResponse struct {
	Total  int64 `json:"total"`
	Unread int64 `json:"unread"`
	Read   int64 `json:"read"`
}

type NotificationSettingResponse struct {
	ID          string `json:"id"`
	SettingType string `json:"setting_type"`
	IsEnabled   bool   `json:"is_enabled"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// 通知类型常量
const (
	NotificationTypeUpdate        = "update"
	NotificationTypeGift          = "gift"
	NotificationTypeLike          = "like"
	NotificationTypeRecommendation = "recommendation"
	NotificationTypeSystem        = "system"
)

// 通知设置类型常量
const (
	SettingTypeUpdate        = "novel_update"
	SettingTypeComment       = "comment_reply"
	SettingTypeLike          = "comment_like"
	SettingTypeRecommendation = "recommendation"
	SettingTypeSystem        = "system_notice"
	SettingTypeMarketing     = "marketing"
)