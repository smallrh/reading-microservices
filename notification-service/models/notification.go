package models

import (
	"time"
	"gorm.io/gorm"
)

type Notification struct {
	ID          string     `gorm:"type:varchar(36);primarykey" json:"id"`
	UserID      string     `gorm:"type:varchar(36);not null;index" json:"user_id"`
	Title       string     `gorm:"type:varchar(200);not null" json:"title"`
	Content     string     `gorm:"type:text;not null" json:"content"`
	Type        string     `gorm:"type:enum('update','gift','like','recommendation','system');not null;index" json:"type"`
	IsRead      bool       `gorm:"default:false;index" json:"is_read"`
	RelatedID   *string    `gorm:"type:varchar(36)" json:"related_id"`
	RelatedType *string    `gorm:"type:varchar(50)" json:"related_type"`
	ActionURL   *string    `gorm:"type:varchar(500)" json:"action_url"`
	CreatedAt   time.Time  `gorm:"index" json:"created_at"`
	ReadAt      *time.Time `json:"read_at"`
}

type NotificationSetting struct {
	ID          string    `gorm:"type:varchar(36);primarykey" json:"id"`
	UserID      string    `gorm:"type:varchar(36);not null;index" json:"user_id"`
	SettingType string    `gorm:"type:varchar(50);not null" json:"setting_type"`
	IsEnabled   bool      `gorm:"default:true" json:"is_enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PushToken struct {
	ID         string    `gorm:"type:varchar(36);primarykey" json:"id"`
	UserID     string    `gorm:"type:varchar(36);not null;index" json:"user_id"`
	DeviceID   string    `gorm:"type:varchar(100);not null" json:"device_id"`
	Token      string    `gorm:"type:varchar(255);not null" json:"token"`
	Platform   string    `gorm:"type:enum('ios','android','web');not null" json:"platform"`
	IsActive   bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	LastUsedAt time.Time `json:"last_used_at"`
}

// Hook functions
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = generateUUID()
	}
	return nil
}

func (ns *NotificationSetting) BeforeCreate(tx *gorm.DB) error {
	if ns.ID == "" {
		ns.ID = generateUUID()
	}
	return nil
}

func (pt *PushToken) BeforeCreate(tx *gorm.DB) error {
	if pt.ID == "" {
		pt.ID = generateUUID()
	}
	return nil
}