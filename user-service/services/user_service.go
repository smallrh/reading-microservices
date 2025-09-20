package services

import (
	"errors"
	"log"
	"reading-microservices/user-service/models"
	"reading-microservices/user-service/repositories"
	services "reading-microservices/user-service/services/auth"
	"reading-microservices/user-service/utils"
	"sync"
	"time"
)

// UserService 组合 AuthManager + SessionManager + LoginLogger
// UserService 组合 AuthManager + SessionManager + LoginLogger
type UserService struct {
	userRepo       repositories.UserRepository
	authManager    *services.AuthManager
	sessionManager *services.SessionManager
	loginLogger    *services.LoginLogger
	loginLocks     sync.Map // 用户登录锁

	// 双 token 配置
	accessExpiresIn  int
	refreshExpiresIn int
}

func NewUserService(
	userRepo repositories.UserRepository,
	authManager *services.AuthManager,
	sessionManager *services.SessionManager,
	loginLogger *services.LoginLogger,
	accessExpiresIn int,
	refreshExpiresIn int,
) *UserService {
	return &UserService{
		userRepo:         userRepo,
		authManager:      authManager,
		sessionManager:   sessionManager,
		loginLogger:      loginLogger,
		accessExpiresIn:  accessExpiresIn,
		refreshExpiresIn: refreshExpiresIn,
	}
}

// ------------------- Register -------------------

func (s *UserService) Register(req *models.RegisterRequest) (*models.LoginResponse, error) {
	// 检查唯一性略，可调用 userRepo.GetByUsername/GetByEmail/GetByPhone

	// 加密密码
	hashedPassword, _ := s.authManager.HashPassword(req.Password)
	user := &models.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
		LoginType:    "password",
	}
	if req.Email != "" {
		user.Email = &req.Email
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 生成双 token
	accessToken, _ := s.authManager.GenerateToken(user.ID, user.Username, s.accessExpiresIn)
	refreshToken, _ := s.authManager.GenerateToken(user.ID, user.Username, s.refreshExpiresIn)

	// 创建 refresh session（存数据库）
	refreshSession := &models.UserSession{
		UserID:         user.ID,
		SessionToken:   refreshToken,
		SessionType:    "refresh",
		Platform:       req.Platform,
		ExpiresAt:      time.Now().Add(time.Duration(s.refreshExpiresIn) * time.Second),
		LastActivityAt: time.Now(),
	}
	if req.DeviceID != "" {
		refreshSession.DeviceID = &req.DeviceID
	}

	// 创建 access session（只存 Redis）
	accessSession := &models.UserSession{
		UserID:         user.ID,
		SessionToken:   accessToken,
		SessionType:    "access",
		Platform:       req.Platform,
		ExpiresAt:      time.Now().Add(time.Duration(s.accessExpiresIn) * time.Second),
		LastActivityAt: time.Now(),
	}
	if req.DeviceID != "" {
		accessSession.DeviceID = &req.DeviceID
	}

	// 保存双 session
	if err := s.sessionManager.CreateSessionPair(accessSession, refreshSession); err != nil {
		return nil, err
	}

	// 登录日志
	s.loginLogger.LogSuccess(user.ID, &refreshSession.ID, req.Platform, req.DeviceID)

	return &models.LoginResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresIn:  s.accessExpiresIn,
		RefreshExpiresIn: s.refreshExpiresIn,
		User:             convertToUserInfo(user),
	}, nil
}

// ------------------- Login -------------------

func (s *UserService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// 获取用户锁，防止并发登录
	lock, _ := s.loginLocks.LoadOrStore(req.Username, &sync.Mutex{})
	mutex := lock.(*sync.Mutex)
	mutex.Lock()
	defer mutex.Unlock()
	defer s.loginLocks.Delete(req.Username)

	// 1. 查询用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// 2. 校验密码
	if err := s.authManager.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		s.loginLogger.LogFailure(user.ID, "invalid password", req.Platform, req.DeviceID)
		return nil, errors.New("invalid username or password")
	}

	// 3. 生成双 token
	accessToken, err := s.authManager.GenerateToken(user.ID, user.Username, s.accessExpiresIn)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := s.authManager.GenerateToken(user.ID, user.Username, s.refreshExpiresIn)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// 4. 清理旧的 session（不中断登录流程）
	if err := s.sessionManager.InvalidateAllUserSessions(user.ID, true); err != nil {
		s.loginLogger.LogFailure(user.ID, "warning: failed to invalidate old sessions", req.Platform, req.DeviceID)
	}

	// 5. 创建新的 session 对
	refreshSession := &models.UserSession{
		UserID:         user.ID,
		SessionToken:   refreshToken,
		SessionType:    "refresh",
		Platform:       req.Platform,
		ExpiresAt:      time.Now().Add(time.Duration(s.refreshExpiresIn) * time.Second),
		LastActivityAt: time.Now(),
	}
	if req.DeviceID != "" {
		refreshSession.DeviceID = &req.DeviceID
	}

	accessSession := &models.UserSession{
		UserID:         user.ID,
		SessionToken:   accessToken,
		SessionType:    "access",
		Platform:       req.Platform,
		ExpiresAt:      time.Now().Add(time.Duration(s.accessExpiresIn) * time.Second),
		LastActivityAt: time.Now(),
	}
	if req.DeviceID != "" {
		accessSession.DeviceID = &req.DeviceID
	}
	refreshSession.AccessToken = accessSession.SessionToken

	// 使用重试机制创建 session 对
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err = s.sessionManager.CreateSessionPair(accessSession, refreshSession)
		if err == nil {
			break
		}

		// 如果是唯一性冲突，重新生成 token
		if utils.IsDuplicateError(err) && i < maxRetries-1 {
			accessToken, _ = s.authManager.GenerateToken(user.ID, user.Username, s.accessExpiresIn)
			refreshToken, _ = s.authManager.GenerateToken(user.ID, user.Username, s.refreshExpiresIn)
			accessSession.SessionToken = accessToken
			refreshSession.SessionToken = refreshToken
			continue
		}

		return nil, err
	}

	// 6. 记录登录成功日志
	s.loginLogger.LogSuccess(user.ID, &refreshSession.ID, req.Platform, req.DeviceID)

	// 7. 返回响应
	return &models.LoginResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresIn:  s.accessExpiresIn,
		RefreshExpiresIn: s.refreshExpiresIn,
		User:             convertToUserInfo(user),
	}, nil
}

// ------------------- Logout -------------------

func (s *UserService) Logout(accessToken string) error {
	// 通过 access token 找到对应的 refresh token 并一起失效
	return s.sessionManager.InvalidateSessionPairByAccessToken(accessToken)
}

// ------------------- RefreshToken -------------------

func (s *UserService) RefreshToken(refreshToken string) (*models.LoginResponse, error) {
	// 验证 refresh token
	claims, err := s.authManager.ParseToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// 获取 refresh session（从数据库）
	refreshSession, err := s.sessionManager.GetSession(refreshToken, false)
	if err != nil {
		return nil, errors.New("refresh session expired or invalid")
	}

	// 检查 refresh token 是否即将过期（比如7天内过期）
	if time.Until(refreshSession.ExpiresAt) < 1*time.Hour {
		// 如果 refresh token 即将过期，可以在这里选择：
		// 1. 返回错误提示用户重新登录
		// 2. 或者生成新的 refresh token（轮换）
		// 这里采用方案1：提示用户重新登录
		return nil, errors.New("refresh token is about to expire, please login again")
	}

	// 只生成新的 access token
	newAccessToken, _ := s.authManager.GenerateToken(claims.UserID, "", s.accessExpiresIn)
	// 创建新的 access session
	newAccessSession := &models.UserSession{
		UserID:         claims.UserID,
		SessionToken:   newAccessToken,
		SessionType:    "access",
		Platform:       refreshSession.Platform,
		ExpiresAt:      time.Now().Add(time.Duration(s.accessExpiresIn) * time.Second),
		LastActivityAt: time.Now(),
		DeviceID:       refreshSession.DeviceID,
	}

	// 更新 refresh session 的 access token 引用和最后活动时间
	refreshSession.AccessToken = newAccessToken
	refreshSession.LastActivityAt = time.Now()

	if err := s.sessionManager.RefreshSessionPair(newAccessSession, refreshSession); err != nil {
		log.Printf("warning: failed to refresh session: %v", err)
		return nil, errors.New("refresh session failed")
	}

	// 更新 refresh session 的最后活动时间（在数据库中）
	if err := s.sessionManager.UpdateRefreshSession(refreshToken, newAccessToken); err != nil {
		log.Printf("warning: failed to update refresh session: %v", err)
	}

	// 获取用户信息
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &models.LoginResponse{
		AccessToken:      newAccessToken,
		RefreshToken:     refreshToken, // 返回相同的 refresh token
		AccessExpiresIn:  s.accessExpiresIn,
		RefreshExpiresIn: int(time.Until(refreshSession.ExpiresAt).Seconds()),
		User:             convertToUserInfo(user),
	}, nil
}

// ------------------- ValidateToken -------------------

func (s *UserService) ValidateToken(accessToken string) (*models.UserSession, error) {
	session, err := s.sessionManager.GetSession(accessToken, true)
	if err != nil {
		return nil, errors.New("invalid or expired access token")
	}
	return session, nil
}

// ------------------- GetProfile -------------------

func (s *UserService) GetProfile(userID string) (*models.UserInfo, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	return convertToUserInfo(user), nil
}

// ------------------- UpdateProfile -------------------

func (s *UserService) UpdateProfile(userID string, req *models.UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if req.Nickname != nil {
		user.Nickname = req.Nickname
	}
	if req.Bio != nil {
		user.Bio = req.Bio
	}
	if req.Gender != nil {
		user.Gender = *req.Gender
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}
	if req.BirthDate != nil {
		if bd, err := time.Parse("2006-01-02", *req.BirthDate); err == nil {
			user.BirthDate = &bd
		}
	}
	return s.userRepo.Update(user)
}

// ------------------- ChangePassword -------------------

func (s *UserService) ChangePassword(userID string, req *models.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if err := s.authManager.VerifyPassword(user.PasswordHash, req.OldPassword); err != nil {
		return errors.New("invalid old password")
	}
	hashed, _ := s.authManager.HashPassword(req.NewPassword)
	user.PasswordHash = hashed
	return s.userRepo.Update(user)
}

// ------------------- Helper -------------------

func convertToUserInfo(user *models.User) *models.UserInfo {
	return &models.UserInfo{
		ID:               user.ID,
		Username:         user.Username,
		Email:            user.Email,
		Phone:            user.Phone,
		AvatarURL:        user.AvatarURL,
		Nickname:         user.Nickname,
		Bio:              user.Bio,
		Gender:           user.Gender,
		Level:            user.Level,
		ExperiencePoints: user.ExperiencePoints,
		ReadingCoins:     user.ReadingCoins,
		VipLevel:         user.VipLevel,
		IsPhoneVerified:  user.IsPhoneVerified,
		IsEmailVerified:  user.IsEmailVerified,
	}
}
