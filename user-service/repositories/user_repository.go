package repositories

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"reading-microservices/user-service/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByPhone(phone string) (*models.User, error)
	Update(user *models.User) error
	UpdateLastLogin(userID string) error
	CreateSession(session *models.UserSession) error
	GetActiveSession(token string) (*models.UserSession, error)
	InvalidateSession(token string) error
	InvalidateAllUserSessions(userID string) error // ⭐ 新增
	CreateLoginLog(log *models.LoginLog) error
	GetThirdPartyAccount(platform, platformUserID string) (*models.ThirdPartyAccount, error)
	CreateThirdPartyAccount(account *models.ThirdPartyAccount) error
	LinkThirdPartyAccount(userID string, account *models.ThirdPartyAccount) error
	UpdateRefreshSession(refreshToken, newAccessToken string, lastActivityAt time.Time) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ? AND is_active = ?", id, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ? AND is_active = ?", username, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? AND is_active = ?", email, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("phone = ? AND is_active = ?", phone, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) UpdateLastLogin(userID string) error {
	now := time.Now()
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("last_login_at", now).Error
}

func (r *userRepository) CreateSession(session *models.UserSession) error {
	return r.db.Create(session).Error
}

func (r *userRepository) GetActiveSession(token string) (*models.UserSession, error) {
	var session models.UserSession
	err := r.db.Where("session_token = ? AND is_active = ? AND expires_at > ?",
		token, true, time.Now()).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *userRepository) InvalidateSession(token string) error {
	return r.db.Model(&models.UserSession{}).
		Where("session_token = ?", token).
		Update("is_active", false).Error
}

func (r *userRepository) InvalidateAllUserSessions(userID string) error {
	return r.db.Model(&models.UserSession{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Update("is_active", false).Error
}

func (r *userRepository) CreateLoginLog(log *models.LoginLog) error {
	return r.db.Create(log).Error
}

func (r *userRepository) GetThirdPartyAccount(platform, platformUserID string) (*models.ThirdPartyAccount, error) {
	var account models.ThirdPartyAccount
	err := r.db.Where("platform = ? AND platform_user_id = ? AND is_active = ?",
		platform, platformUserID, true).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *userRepository) CreateThirdPartyAccount(account *models.ThirdPartyAccount) error {
	return r.db.Create(account).Error
}

func (r *userRepository) LinkThirdPartyAccount(userID string, account *models.ThirdPartyAccount) error {
	account.UserID = userID
	return r.db.Create(account).Error
}

// UpdateRefreshSession 更新 refresh session
func (r *userRepository) UpdateRefreshSession(refreshToken, newAccessToken string, lastActivityAt time.Time) error {
	result := r.db.Model(&models.UserSession{}).
		Where("session_token = ? AND session_type = ?", refreshToken, "refresh").
		Updates(map[string]interface{}{
			"last_activity_at": lastActivityAt,
			"updated_at":       time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("refresh session not found")
	}

	return nil
}
