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

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Platform string `json:"platform" binding:"required,oneof=ios android web h5"`
	DeviceID string `json:"device_id"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty"`
	Platform string `json:"platform" binding:"required,oneof=ios android web h5"`
	DeviceID string `json:"device_id"`
}

type UpdateProfileRequest struct {
	Nickname  *string `json:"nickname"`
	Bio       *string `json:"bio"`
	Gender    *string `json:"gender" binding:"omitempty,oneof=male female other"`
	BirthDate *string `json:"birth_date"`
	AvatarURL *string `json:"avatar_url"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type LoginResponse struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	AccessExpiresIn  int       `json:"access_expires_in"`
	RefreshExpiresIn int       `json:"refresh_expires_in"`
	User             *UserInfo `json:"user"`
}

type UserInfo struct {
	ID               string  `json:"id"`
	Username         string  `json:"username"`
	Email            *string `json:"email"`
	Phone            *string `json:"phone"`
	AvatarURL        *string `json:"avatar_url"`
	Nickname         *string `json:"nickname"`
	Bio              *string `json:"bio"`
	Gender           string  `json:"gender"`
	Level            int     `json:"level"`
	ExperiencePoints int     `json:"experience_points"`
	ReadingCoins     int     `json:"reading_coins"`
	VipLevel         string  `json:"vip_level"`
	IsPhoneVerified  bool    `json:"is_phone_verified"`
	IsEmailVerified  bool    `json:"is_email_verified"`
}
