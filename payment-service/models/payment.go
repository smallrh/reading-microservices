package models

import (
	"gorm.io/gorm"
	"time"
)

type VipMembership struct {
	ID            string    `gorm:"type:varchar(36);primarykey" json:"id"`
	UserID        string    `gorm:"type:varchar(36);not null;index" json:"user_id"`
	VipType       string    `gorm:"type:enum('vip','svip');not null" json:"vip_type"`
	StartDate     time.Time `gorm:"not null" json:"start_date"`
	EndDate       time.Time `gorm:"not null;index" json:"end_date"`
	IsActive      bool      `gorm:"default:true;index" json:"is_active"`
	AutoRenew     bool      `gorm:"default:false" json:"auto_renew"`
	PaymentMethod *string   `gorm:"type:varchar(50)" json:"payment_method"`
	Amount        *float64  `gorm:"type:decimal(10,2)" json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

type PointsRecord struct {
	ID          string    `gorm:"type:varchar(36);primarykey" json:"id"`
	UserID      string    `gorm:"type:varchar(36);not null;index" json:"user_id"`
	Points      int       `gorm:"not null" json:"points"`
	PointsType  string    `gorm:"type:enum('earn','spend');not null;index" json:"points_type"`
	Source      string    `gorm:"type:varchar(100);not null" json:"source"`
	Description *string   `gorm:"type:text" json:"description"`
	RelatedID   *string   `gorm:"type:varchar(36)" json:"related_id"`
	RelatedType *string   `gorm:"type:varchar(50)" json:"related_type"`
	CreatedAt   time.Time `json:"created_at"`
}

type CoinsRecord struct {
	ID          string    `gorm:"type:varchar(36);primarykey" json:"id"`
	UserID      string    `gorm:"type:varchar(36);not null;index" json:"user_id"`
	Coins       int       `gorm:"not null" json:"coins"`
	CoinsType   string    `gorm:"type:enum('earn','spend');not null;index" json:"coins_type"`
	Source      string    `gorm:"type:varchar(100);not null" json:"source"`
	Description *string   `gorm:"type:text" json:"description"`
	RelatedID   *string   `gorm:"type:varchar(36)" json:"related_id"`
	RelatedType *string   `gorm:"type:varchar(50)" json:"related_type"`
	CreatedAt   time.Time `json:"created_at"`
}

type CheckinRecord struct {
	ID              string    `gorm:"type:varchar(36);primarykey" json:"id"`
	UserID          string    `gorm:"type:varchar(36);not null;index" json:"user_id"`
	CheckinDate     time.Time `gorm:"type:date;not null;index" json:"checkin_date"`
	ConsecutiveDays int       `gorm:"default:1" json:"consecutive_days"`
	PointsEarned    int       `gorm:"default:0" json:"points_earned"`
	CoinsEarned     int       `gorm:"default:0" json:"coins_earned"`
	CreatedAt       time.Time `json:"created_at"`
}

type Gift struct {
	ID          string    `gorm:"type:varchar(36);primarykey" json:"id"`
	Name        string    `gorm:"type:varchar(200);not null" json:"name"`
	Description *string   `gorm:"type:text" json:"description"`
	GiftType    string    `gorm:"type:enum('reading_coins','vip_time','physical_gift','points','other');not null;index" json:"gift_type"`
	Value       string    `gorm:"type:varchar(100);not null" json:"value"`
	Category    *string   `gorm:"type:varchar(50);index" json:"category"`
	ImageURL    *string   `gorm:"type:varchar(500)" json:"image_url"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserGift struct {
	ID         string     `gorm:"type:varchar(36);primarykey" json:"id"`
	UserID     string     `gorm:"type:varchar(36);not null;index" json:"user_id"`
	GiftID     string     `gorm:"type:varchar(36);not null" json:"gift_id"`
	Status     string     `gorm:"type:enum('unused','used','expired');default:'unused';index" json:"status"`
	ObtainedAt time.Time  `gorm:"type:datetime(3);autoCreateTime;autoUpdateTime" json:"obtained_at"`
	UsedAt     *time.Time `json:"used_at"`
	ExpiresAt  *time.Time `gorm:"index" json:"expires_at"`
	Source     *string    `gorm:"type:varchar(100)" json:"source"`
	CreatedAt  time.Time  `json:"created_at"`

	// 关联
	Gift Gift `gorm:"foreignKey:GiftID" json:"gift,omitempty"`
}

type RedeemCode struct {
	ID        string     `gorm:"type:varchar(36);primarykey" json:"id"`
	Code      string     `gorm:"type:varchar(100);not null;uniqueIndex" json:"code"`
	GiftID    string     `gorm:"type:varchar(36);not null;index" json:"gift_id"`
	IsUsed    bool       `gorm:"default:false;index" json:"is_used"`
	UsedBy    *string    `gorm:"type:varchar(36)" json:"used_by"`
	UsedAt    *time.Time `json:"used_at"`
	ExpiresAt *time.Time `gorm:"index" json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`

	// 关联
	Gift Gift `gorm:"foreignKey:GiftID" json:"gift,omitempty"`
}

// Hook functions
func (v *VipMembership) BeforeCreate(tx *gorm.DB) error {
	if v.ID == "" {
		v.ID = generateUUID()
	}
	return nil
}

func (p *PointsRecord) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = generateUUID()
	}
	return nil
}

func (c *CoinsRecord) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = generateUUID()
	}
	return nil
}

func (c *CheckinRecord) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = generateUUID()
	}
	return nil
}

func (g *Gift) BeforeCreate(tx *gorm.DB) error {
	if g.ID == "" {
		g.ID = generateUUID()
	}
	return nil
}

func (u *UserGift) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = generateUUID()
	}
	return nil
}

func (r *RedeemCode) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = generateUUID()
	}
	return nil
}
