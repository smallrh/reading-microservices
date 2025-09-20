package utils

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountExists      = errors.New("account already exists")
	ErrSessionExpired     = errors.New("session expired")
)

// IsDuplicateError 检查是否为重复错误
func IsDuplicateError(err error) bool {
	return err != nil && (contains(err.Error(), "1062") ||
		contains(err.Error(), "Duplicate entry") ||
		contains(err.Error(), "unique constraint"))
}
func contains(s, substr string) bool {
	// 简单的字符串包含检查
	// 在实际项目中可以使用更复杂的方法
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
