package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Random   string `json:"rand"` // 确保唯一性
	jwt.RegisteredClaims
}

// 添加随机后缀确保token唯一性
func GenerateToken(userID, username, secret string, expiresIn int) (string, error) {
	// 生成随机后缀
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	randomSuffix := hex.EncodeToString(randomBytes)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresIn) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Random: randomSuffix, // 添加随机字段
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(token string, secret string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

func RefreshToken(token string, secret string, expiresIn int) (string, error) {
	claims, err := ParseToken(token, secret)
	if err != nil {
		return "", err
	}

	// 如果token还有超过1小时的有效期，就不刷新
	if time.Until(claims.ExpiresAt.Time) > time.Hour {
		return "", errors.New("token doesn't need refresh")
	}

	return GenerateToken(claims.UserID, claims.Username, secret, expiresIn)
}
