package services

import (
	"golang.org/x/crypto/bcrypt"
	"reading-microservices/shared/utils"
)

// AuthManager 只负责用户认证和 Token 生成
type AuthManager struct {
	jwtSecret string
}

func NewAuthManager(jwtSecret string) *AuthManager {
	return &AuthManager{
		jwtSecret: jwtSecret,
	}
}

// VerifyPassword 验证明文密码与哈希是否匹配
func (a *AuthManager) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// HashPassword 加密密码
func (a *AuthManager) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// GenerateToken 根据用户 ID 生成 JWT，支持自定义过期时间
func (a *AuthManager) GenerateToken(userID, username string, expiresIn int) (string, error) {
	return utils.GenerateToken(userID, username, a.jwtSecret, expiresIn)
}

// ParseToken 验证并解析 JWT
func (a *AuthManager) ParseToken(token string) (*utils.Claims, error) {
	return utils.ParseToken(token, a.jwtSecret)
}
