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
type CreateVipMembershipRequest struct {
	VipType       string  `json:"vip_type" binding:"required,oneof=vip svip"`
	Duration      int     `json:"duration" binding:"required,min=1"` // 月数
	PaymentMethod string  `json:"payment_method" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	AutoRenew     bool    `json:"auto_renew"`
}

type EarnPointsRequest struct {
	Points      int     `json:"points" binding:"required,gt=0"`
	Source      string  `json:"source" binding:"required"`
	Description *string `json:"description"`
	RelatedID   *string `json:"related_id"`
	RelatedType *string `json:"related_type"`
}

type SpendPointsRequest struct {
	Points      int     `json:"points" binding:"required,gt=0"`
	Source      string  `json:"source" binding:"required"`
	Description *string `json:"description"`
	RelatedID   *string `json:"related_id"`
	RelatedType *string `json:"related_type"`
}

type EarnCoinsRequest struct {
	Coins       int     `json:"coins" binding:"required,gt=0"`
	Source      string  `json:"source" binding:"required"`
	Description *string `json:"description"`
	RelatedID   *string `json:"related_id"`
	RelatedType *string `json:"related_type"`
}

type SpendCoinsRequest struct {
	Coins       int     `json:"coins" binding:"required,gt=0"`
	Source      string  `json:"source" binding:"required"`
	Description *string `json:"description"`
	RelatedID   *string `json:"related_id"`
	RelatedType *string `json:"related_type"`
}

type CreateGiftRequest struct {
	Name        string  `json:"name" binding:"required,max=200"`
	Description *string `json:"description"`
	GiftType    string  `json:"gift_type" binding:"required,oneof=reading_coins vip_time physical_gift points other"`
	Value       string  `json:"value" binding:"required"`
	Category    *string `json:"category"`
	ImageURL    *string `json:"image_url"`
}

type RedeemCodeRequest struct {
	Code string `json:"code" binding:"required"`
}

// 响应DTO
type VipMembershipResponse struct {
	ID            string  `json:"id"`
	VipType       string  `json:"vip_type"`
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	IsActive      bool    `json:"is_active"`
	AutoRenew     bool    `json:"auto_renew"`
	PaymentMethod *string `json:"payment_method"`
	Amount        *float64 `json:"amount"`
	DaysRemaining int     `json:"days_remaining"`
	CreatedAt     string  `json:"created_at"`
}

type PointsRecordResponse struct {
	ID          string  `json:"id"`
	Points      int     `json:"points"`
	PointsType  string  `json:"points_type"`
	Source      string  `json:"source"`
	Description *string `json:"description"`
	RelatedID   *string `json:"related_id"`
	RelatedType *string `json:"related_type"`
	CreatedAt   string  `json:"created_at"`
}

type CoinsRecordResponse struct {
	ID          string  `json:"id"`
	Coins       int     `json:"coins"`
	CoinsType   string  `json:"coins_type"`
	Source      string  `json:"source"`
	Description *string `json:"description"`
	RelatedID   *string `json:"related_id"`
	RelatedType *string `json:"related_type"`
	CreatedAt   string  `json:"created_at"`
}

type CheckinResponse struct {
	ID             string `json:"id"`
	CheckinDate    string `json:"checkin_date"`
	ConsecutiveDays int    `json:"consecutive_days"`
	PointsEarned   int    `json:"points_earned"`
	CoinsEarned    int    `json:"coins_earned"`
	CreatedAt      string `json:"created_at"`
}

type UserGiftResponse struct {
	ID         string       `json:"id"`
	Status     string       `json:"status"`
	ObtainedAt string       `json:"obtained_at"`
	UsedAt     *string      `json:"used_at"`
	ExpiresAt  *string      `json:"expires_at"`
	Source     *string      `json:"source"`
	Gift       GiftResponse `json:"gift"`
}

type GiftResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	GiftType    string  `json:"gift_type"`
	Value       string  `json:"value"`
	Category    *string `json:"category"`
	ImageURL    *string `json:"image_url"`
}

type WalletResponse struct {
	UserID             string `json:"user_id"`
	Points             int    `json:"points"`
	ReadingCoins       int    `json:"reading_coins"`
	VipType            string `json:"vip_type"`
	VipExpiresAt       *string `json:"vip_expires_at"`
	ConsecutiveCheckins int    `json:"consecutive_checkins"`
}

type CheckinStatusResponse struct {
	CanCheckin         bool   `json:"can_checkin"`
	LastCheckinDate    *string `json:"last_checkin_date"`
	ConsecutiveDays    int    `json:"consecutive_days"`
	TodayPointsReward  int    `json:"today_points_reward"`
	TodayCoinsReward   int    `json:"today_coins_reward"`
}

type PointsStatsResponse struct {
	TotalEarned    int `json:"total_earned"`
	TotalSpent     int `json:"total_spent"`
	CurrentBalance int `json:"current_balance"`
	ThisMonthEarned int `json:"this_month_earned"`
	ThisMonthSpent  int `json:"this_month_spent"`
}

type CoinsStatsResponse struct {
	TotalEarned     int `json:"total_earned"`
	TotalSpent      int `json:"total_spent"`
	CurrentBalance  int `json:"current_balance"`
	ThisMonthEarned int `json:"this_month_earned"`
	ThisMonthSpent  int `json:"this_month_spent"`
}