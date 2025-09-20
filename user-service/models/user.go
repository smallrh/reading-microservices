package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID                  string     `gorm:"type:varchar(36);primarykey" json:"id"`
	Username            string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email               *string    `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Phone               *string    `gorm:"type:varchar(20);uniqueIndex" json:"phone"`
	PasswordHash        string     `gorm:"type:varchar(255);not null" json:"-"`
	AvatarURL           *string    `gorm:"type:varchar(500)" json:"avatar_url"`
	Nickname            *string    `gorm:"type:varchar(50)" json:"nickname"`
	Bio                 *string    `gorm:"type:text" json:"bio"`
	Gender              string     `gorm:"type:enum('male','female','other');default:'other'" json:"gender"`
	BirthDate           *time.Time `gorm:"type:date" json:"birth_date"`
	Level               int        `gorm:"default:1" json:"level"`
	ExperiencePoints    int        `gorm:"default:0" json:"experience_points"`
	ReadingCoins        int        `gorm:"default:0" json:"reading_coins"`
	VipLevel            string     `gorm:"type:enum('none','vip','svip');default:'none'" json:"vip_level"`
	VipExpiresAt        *time.Time `json:"vip_expires_at"`
	LoginType           string     `gorm:"type:enum('password','wechat','qq','weibo','apple');default:'password'" json:"login_type"`
	IsPhoneVerified     bool       `gorm:"default:false" json:"is_phone_verified"`
	IsEmailVerified     bool       `gorm:"default:false" json:"is_email_verified"`
	PhoneVerifiedAt     *time.Time `json:"phone_verified_at"`
	EmailVerifiedAt     *time.Time `json:"email_verified_at"`
	LastThirdPartyLogin *time.Time `json:"last_third_party_login"`
	IsActive            bool       `gorm:"default:true" json:"is_active"`
	LastLoginAt         *time.Time `json:"last_login_at"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`

	// 关联
	ThirdPartyAccounts []ThirdPartyAccount `gorm:"foreignKey:UserID" json:"third_party_accounts,omitempty"`
	Sessions           []UserSession       `gorm:"foreignKey:UserID" json:"sessions,omitempty"`
	LoginLogs          []LoginLog          `gorm:"foreignKey:UserID" json:"login_logs,omitempty"`
}

type ThirdPartyAccount struct {
	ID               string     `gorm:"type:varchar(36);primarykey" json:"id"`
	UserID           string     `gorm:"type:varchar(36);not null;index" json:"user_id"`
	Platform         string     `gorm:"type:enum('wechat','qq','weibo','apple');not null" json:"platform"`
	PlatformUserID   string     `gorm:"type:varchar(100);not null" json:"platform_user_id"`
	PlatformUsername *string    `gorm:"type:varchar(100)" json:"platform_username"`
	PlatformNickname *string    `gorm:"type:varchar(100)" json:"platform_nickname"`
	PlatformAvatar   *string    `gorm:"type:varchar(500)" json:"platform_avatar"`
	PlatformEmail    *string    `gorm:"type:varchar(100)" json:"platform_email"`
	PlatformPhone    *string    `gorm:"type:varchar(20)" json:"platform_phone"`
	UnionID          *string    `gorm:"type:varchar(100);uniqueIndex" json:"union_id"`
	OpenID           *string    `gorm:"type:varchar(100);index" json:"open_id"`
	AccessToken      *string    `gorm:"type:text" json:"-"`
	RefreshToken     *string    `gorm:"type:text" json:"-"`
	TokenExpiresAt   *time.Time `json:"token_expires_at"`
	IsPrimary        bool       `gorm:"default:false" json:"is_primary"`
	IsActive         bool       `gorm:"default:true" json:"is_active"`
	LastLoginAt      *time.Time `json:"last_login_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type UserSession struct {
	ID             string    `gorm:"primaryKey" json:"id"`
	UserID         string    `gorm:"index;not null" json:"user_id"`
	SessionToken   string    `gorm:"type:varchar(512);uniqueIndex;not null" json:"session_token"` // 现在是 refresh token
	SessionType    string    `gorm:"not null;default:'refresh'" json:"session_type"`              // refresh 或 access（历史兼容）
	Platform       string    `gorm:"not null" json:"platform"`
	DeviceID       *string   `gorm:"index" json:"device_id"`
	IPAddress      *string   `json:"ip_address"`
	UserAgent      *string   `json:"user_agent"`
	ExpiresAt      time.Time `gorm:"not null" json:"expires_at"`
	LastActivityAt time.Time `gorm:"not null" json:"last_activity_at"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联字段
	AccessToken     string `gorm:"-" json:"access_token"`      // 仅用于响应，不存数据库
	AccessExpiresIn int    `gorm:"-" json:"access_expires_in"` // 仅用于响应
}

type LoginLog struct {
	ID            string    `gorm:"type:varchar(36);primarykey" json:"id"`
	UserID        string    `gorm:"type:varchar(36);not null;index" json:"user_id"`
	LoginType     string    `gorm:"type:enum('password','wechat','qq','weibo','apple');not null" json:"login_type"`
	Platform      string    `gorm:"type:enum('ios','android','web','h5');not null" json:"platform"`
	DeviceID      *string   `gorm:"type:varchar(100)" json:"device_id"`
	DeviceInfo    *string   `gorm:"type:text" json:"device_info"`
	IPAddress     *string   `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent     *string   `gorm:"type:text" json:"user_agent"`
	Location      *string   `gorm:"type:varchar(200)" json:"location"`
	IsSuccess     bool      `gorm:"default:true" json:"is_success"`
	FailureReason *string   `gorm:"type:varchar(200)" json:"failure_reason"`
	SessionID     *string   `gorm:"type:varchar(100)" json:"session_id"`
	CreatedAt     time.Time `json:"created_at"`
}

// BeforeCreate 钩子函数
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = generateUUID()
	}
	return nil
}

func (t *ThirdPartyAccount) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = generateUUID()
	}
	return nil
}

func (s *UserSession) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = generateUUID()
	}
	return nil
}

func (l *LoginLog) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = generateUUID()
	}
	return nil
}
