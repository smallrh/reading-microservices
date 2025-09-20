package services

import (
	"reading-microservices/user-service/models"
)

// UserServiceInterface 定义了用户服务需要实现的所有方法
type UserServiceInterface interface {
	// Register 用户注册
	Register(req *models.RegisterRequest) (*models.LoginResponse, error)

	// Login 用户登录
	Login(req *models.LoginRequest) (*models.LoginResponse, error)

	// Logout 用户登出
	Logout(accessToken string) error

	// RefreshToken 刷新令牌（双token模式）
	RefreshToken(refreshToken string) (*models.LoginResponse, error)

	// ValidateToken 验证令牌有效性
	ValidateToken(accessToken string) (*models.UserSession, error)

	// GetProfile 获取用户资料
	GetProfile(userID string) (*models.UserInfo, error)

	// UpdateProfile 更新用户资料
	UpdateProfile(userID string, req *models.UpdateProfileRequest) error

	// ChangePassword 修改密码
	ChangePassword(userID string, req *models.ChangePasswordRequest) error
}

// 可选：验证 UserService 是否实现了接口
var _ UserServiceInterface = (*UserService)(nil)
