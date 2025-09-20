package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"reading-microservices/user-service/models"
	"reading-microservices/user-service/repositories"
	"reading-microservices/user-service/utils"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	accessSessionPrefix  = "access_session:"
	refreshSessionPrefix = "refresh_session:"
	userSessionsPrefix   = "user_sessions:"
	redisJitterMax       = 300 // 秒
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type SessionManager struct {
	userRepo    repositories.UserRepository
	redisClient *redis.Client
}

func NewSessionManager(userRepo repositories.UserRepository, redisClient *redis.Client) *SessionManager {
	return &SessionManager{
		userRepo:    userRepo,
		redisClient: redisClient,
	}
}

// ----------------- Helpers -----------------

func (s *SessionManager) accessSessionKey(token string) string {
	return accessSessionPrefix + token
}

func (s *SessionManager) refreshSessionKey(token string) string {
	return refreshSessionPrefix + token
}

func (s *SessionManager) userSessionsKey(userID string) string {
	return userSessionsPrefix + userID
}

func ttlWithJitter(expiresAt time.Time) time.Duration {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return time.Second
	}
	return ttl + time.Duration(rand.Intn(redisJitterMax))*time.Second
}

// ----------------- Session CRUD -----------------

// CreateSessionPair 创建双 token 会话对
func (s *SessionManager) CreateSessionPair(accessSession, refreshSession *models.UserSession) error {
	ctx := context.Background()

	// 检查是否已存在
	if existing, _ := s.GetSession(accessSession.SessionToken, true); existing != nil {
		return errors.New("access session already exists")
	}
	if existing, _ := s.GetSession(refreshSession.SessionToken, false); existing != nil {
		return errors.New("refresh session already exists")
	}

	// --- 写数据库（只存 refresh token）---
	maxRetries := 2
	var err error
	for i := 0; i < maxRetries; i++ {
		err = s.userRepo.CreateSession(refreshSession) // 只保存 refresh token 到数据库
		if err == nil {
			break
		}
		if utils.IsDuplicateError(err) {
			return err
		}
		if i < maxRetries-1 {
			time.Sleep(time.Duration(100*(i+1)) * time.Millisecond)
		}
	}
	if err != nil {
		return err
	}

	// --- 写Redis（存储 access 和 refresh token）---
	accessData, err := json.Marshal(accessSession)
	if err != nil {
		return err
	}

	refreshData, err := json.Marshal(refreshSession)
	if err != nil {
		return err
	}

	// 存储 access session（只存 Redis）
	accessKey := s.accessSessionKey(accessSession.SessionToken)
	if err := s.redisClient.Set(ctx, accessKey, accessData, ttlWithJitter(accessSession.ExpiresAt)).Err(); err != nil {
		log.Printf("warning: redis set access session failed: %v", err)
	}

	// 存储 refresh session（Redis 也存一份用于快速验证）
	refreshKey := s.refreshSessionKey(refreshSession.SessionToken)
	if err := s.redisClient.Set(ctx, refreshKey, refreshData, ttlWithJitter(refreshSession.ExpiresAt)).Err(); err != nil {
		log.Printf("warning: redis set refresh session failed: %v", err)
	}

	// 把 session keys 放入 user_sessions:{userID} 集合
	userSessionsKey := s.userSessionsKey(accessSession.UserID)
	if err := s.redisClient.SAdd(ctx, userSessionsKey, accessKey, refreshKey).Err(); err != nil {
		log.Printf("warning: redis add user_sessions failed: %v", err)
	}

	return nil
}

// RefreshSession 刷新会话（更新现有的 refresh session）
func (s *SessionManager) RefreshSession(accessSession *models.UserSession, refreshToken string) error {
	ctx := context.Background()

	// 检查新的 access session 是否已存在
	if existing, _ := s.GetSession(accessSession.SessionToken, true); existing != nil {
		return errors.New("access session already exists")
	}

	// 获取现有的 refresh session
	existingRefreshSession, err := s.GetSession(refreshToken, false)
	if err != nil {
		return fmt.Errorf("refresh session not found: %v", err)
	}

	// 更新 refresh session 的 access token 引用和最后活动时间
	existingRefreshSession.AccessToken = accessSession.SessionToken
	existingRefreshSession.LastActivityAt = time.Now()

	// --- 更新数据库中的 refresh session ---
	if err := s.userRepo.UpdateRefreshSession(refreshToken, accessSession.SessionToken, time.Now()); err != nil {
		return fmt.Errorf("failed to update refresh session in database: %v", err)
	}

	// --- 更新 Redis 中的 refresh session ---
	refreshData, err := json.Marshal(existingRefreshSession)
	if err != nil {
		return err
	}

	refreshKey := s.refreshSessionKey(refreshToken)
	if err := s.redisClient.Set(ctx, refreshKey, refreshData, ttlWithJitter(existingRefreshSession.ExpiresAt)).Err(); err != nil {
		log.Printf("warning: redis update refresh session failed: %v", err)
	}

	// --- 创建新的 access session ---
	accessData, err := json.Marshal(accessSession)
	if err != nil {
		return err
	}

	accessKey := s.accessSessionKey(accessSession.SessionToken)
	if err := s.redisClient.Set(ctx, accessKey, accessData, ttlWithJitter(accessSession.ExpiresAt)).Err(); err != nil {
		log.Printf("warning: redis set access session failed: %v", err)
	}

	// 把新的 access session key 放入 user_sessions:{userID} 集合
	userSessionsKey := s.userSessionsKey(accessSession.UserID)
	if err := s.redisClient.SAdd(ctx, userSessionsKey, accessKey).Err(); err != nil {
		log.Printf("warning: redis add user_sessions failed: %v", err)
	}

	return nil
}

// GetSession 根据 token 获取会话，isAccessToken 标识是否为 access token
func (s *SessionManager) GetSession(token string, isAccessToken bool) (*models.UserSession, error) {
	ctx := context.Background()

	var key string
	if isAccessToken {
		key = s.accessSessionKey(token)
	} else {
		key = s.refreshSessionKey(token)
	}

	// 先尝试从Redis获取
	data, err := s.redisClient.Get(ctx, key).Result()
	if err == nil {
		var session models.UserSession
		if err := json.Unmarshal([]byte(data), &session); err == nil {
			if time.Now().After(session.ExpiresAt) {
				s.redisClient.Del(ctx, key).Err()
				return nil, redis.Nil
			}
			return &session, nil
		}
		log.Printf("warning: unmarshal redis session failed: %v", err)
	}

	// 如果是 refresh token 且 Redis 中没有，从数据库获取
	if !isAccessToken {
		session, err := s.userRepo.GetActiveSession(token)
		if err != nil {
			return nil, err
		}
		if time.Now().After(session.ExpiresAt) {
			s.userRepo.InvalidateSession(token)
			return nil, redis.Nil
		}

		// 回写 Redis
		b, _ := json.Marshal(session)
		_ = s.redisClient.Set(ctx, key, b, ttlWithJitter(session.ExpiresAt)).Err()
		_ = s.redisClient.SAdd(ctx, s.userSessionsKey(session.UserID), key).Err()

		return session, nil
	}

	// access token 只在 Redis 中，如果 Redis 没有就返回错误
	return nil, redis.Nil
}

// InvalidateSession 让单个session失效
func (s *SessionManager) InvalidateSession(token string, isAccessToken bool) error {
	ctx := context.Background()

	var key string
	if isAccessToken {
		key = s.accessSessionKey(token)
	} else {
		key = s.refreshSessionKey(token)
	}

	// 从 Redis 获取 session 信息
	data, _ := s.redisClient.Get(ctx, key).Result()
	if data != "" {
		var session models.UserSession
		if err := json.Unmarshal([]byte(data), &session); err == nil {
			userSessionsKey := s.userSessionsKey(session.UserID)
			s.redisClient.SRem(ctx, userSessionsKey, key).Err()
		}
	}

	// 删除 Redis 里的 session
	_ = s.redisClient.Del(ctx, key).Err()

	// 如果是 refresh token，还需要更新数据库
	if !isAccessToken {
		return s.userRepo.InvalidateSession(token)
	}

	return nil
}

// InvalidateAllUserSessions 让用户所有session失效
func (s *SessionManager) InvalidateAllUserSessions(userID string, invalidateDB bool) error {
	ctx := context.Background()
	userSessionsKey := s.userSessionsKey(userID)

	// 找到该用户所有 session keys
	sessionKeys, err := s.redisClient.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		return err
	}

	// 删除所有 session
	for _, key := range sessionKeys {
		if err := s.redisClient.Del(ctx, key).Err(); err != nil {
			log.Printf("warning: redis del session failed: %v", err)
		}
	}

	// 删除 user_sessions:{userID}
	if err := s.redisClient.Del(ctx, userSessionsKey).Err(); err != nil {
		log.Printf("warning: redis del user_sessions failed: %v", err)
	}

	if !invalidateDB {
		return nil
	}
	// 更新数据库（使所有 refresh token 失效）
	return s.userRepo.InvalidateAllUserSessions(userID)
}

// RefreshSessionPair 刷新会话对
func (s *SessionManager) RefreshSessionPair(newAccessSession, newRefreshSession *models.UserSession) error {
	// 使旧的会话对失效
	if err := s.InvalidateAllUserSessions(newRefreshSession.UserID, false); err != nil {
		log.Printf("warning: invalidate old session pair failed: %v", err)
	}

	// 更新sql session
	return s.RefreshSession(newAccessSession, newRefreshSession.SessionToken)
}

// UpdateRefreshSession 更新 refresh session 的 access token 引用和最后活动时间
func (s *SessionManager) UpdateRefreshSession(refreshToken, newAccessToken string) error {
	ctx := context.Background()

	// 获取当前的 refresh session
	refreshKey := s.refreshSessionKey(refreshToken)
	data, err := s.redisClient.Get(ctx, refreshKey).Result()
	if err != nil {
		return errors.New("refresh session not found in redis")
	}

	var refreshSession models.UserSession
	if err := json.Unmarshal([]byte(data), &refreshSession); err != nil {
		return errors.New("invalid refresh session data")
	}

	// 更新字段
	refreshSession.AccessToken = newAccessToken
	refreshSession.LastActivityAt = time.Now()

	// 更新 Redis
	updatedData, err := json.Marshal(refreshSession)
	if err != nil {
		return err
	}

	if err := s.redisClient.Set(ctx, refreshKey, updatedData, ttlWithJitter(refreshSession.ExpiresAt)).Err(); err != nil {
		return fmt.Errorf("redis update refresh session failed: %v", err)
	}

	// 更新数据库（只更新必要的字段）
	return s.userRepo.UpdateRefreshSession(refreshToken, newAccessToken, time.Now())
}

// InvalidateSessionPairByAccessToken 通过 access token 找到对应的用户并失效所有会话
func (s *SessionManager) InvalidateSessionPairByAccessToken(accessToken string) error {
	ctx := context.Background()

	// 获取 access session 信息
	accessKey := s.accessSessionKey(accessToken)
	data, err := s.redisClient.Get(ctx, accessKey).Result()
	if err != nil {
		if err == redis.Nil {
			return errors.New("access session not found")
		}
		return fmt.Errorf("failed to get access session: %v", err)
	}

	var accessSession models.UserSession
	if err := json.Unmarshal([]byte(data), &accessSession); err != nil {
		return errors.New("invalid access session data")
	}

	// 获取用户ID
	userID := accessSession.UserID
	if userID == "" {
		return errors.New("user ID not found in access session")
	}

	// 失效该用户的所有会话
	return s.InvalidateAllUserSessions(userID, true)
}
