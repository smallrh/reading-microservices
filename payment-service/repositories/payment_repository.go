package repositories

import (
	"time"
	"gorm.io/gorm"
	"reading-microservices/payment-service/models"
)

type PaymentRepository interface {
	// VIP Membership
	CreateVipMembership(membership *models.VipMembership) error
	GetActiveVipMembership(userID string) (*models.VipMembership, error)
	GetUserVipMemberships(userID string, page, size int) ([]models.VipMembership, int64, error)
	UpdateVipMembership(membership *models.VipMembership) error
	ExpireVipMembership(membershipID string) error

	// Points
	CreatePointsRecord(record *models.PointsRecord) error
	GetUserPointsRecords(userID string, page, size int) ([]models.PointsRecord, int64, error)
	GetUserPointsBalance(userID string) (int, error)
	GetPointsStats(userID string) (*models.PointsStatsResponse, error)

	// Coins
	CreateCoinsRecord(record *models.CoinsRecord) error
	GetUserCoinsRecords(userID string, page, size int) ([]models.CoinsRecord, int64, error)
	GetUserCoinsBalance(userID string) (int, error)
	GetCoinsStats(userID string) (*models.CoinsStatsResponse, error)

	// Checkin
	CreateCheckinRecord(record *models.CheckinRecord) error
	GetLastCheckinRecord(userID string) (*models.CheckinRecord, error)
	GetUserCheckinRecords(userID string, page, size int) ([]models.CheckinRecord, int64, error)
	GetCheckinStats(userID string) (int, error) // 连续签到天数

	// Gifts
	CreateGift(gift *models.Gift) error
	GetGifts(category string, isActive *bool) ([]models.Gift, error)
	GetGiftByID(id string) (*models.Gift, error)
	UpdateGift(gift *models.Gift) error
	DeleteGift(id string) error

	// User Gifts
	CreateUserGift(userGift *models.UserGift) error
	GetUserGifts(userID string, status string, page, size int) ([]models.UserGift, int64, error)
	GetUserGiftByID(userID, giftID string) (*models.UserGift, error)
	UpdateUserGiftStatus(userGiftID, status string, usedAt *time.Time) error

	// Redeem Codes
	CreateRedeemCode(code *models.RedeemCode) error
	GetRedeemCodeByCode(code string) (*models.RedeemCode, error)
	UseRedeemCode(codeID, userID string) error
	GetRedeemCodes(isUsed *bool, page, size int) ([]models.RedeemCode, int64, error)

	// Wallet
	GetUserWallet(userID string) (*models.WalletResponse, error)
	UpdateUserBalance(userID string, pointsDelta, coinsDelta int) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

// VIP Membership
func (r *paymentRepository) CreateVipMembership(membership *models.VipMembership) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 使当前活跃的会员失效
		if err := tx.Model(&models.VipMembership{}).
			Where("user_id = ? AND is_active = true", membership.UserID).
			Update("is_active", false).Error; err != nil {
			return err
		}

		// 创建新会员记录
		return tx.Create(membership).Error
	})
}

func (r *paymentRepository) GetActiveVipMembership(userID string) (*models.VipMembership, error) {
	var membership models.VipMembership
	err := r.db.Where("user_id = ? AND is_active = true AND end_date > ?", userID, time.Now()).
		First(&membership).Error
	if err != nil {
		return nil, err
	}
	return &membership, nil
}

func (r *paymentRepository) GetUserVipMemberships(userID string, page, size int) ([]models.VipMembership, int64, error) {
	var memberships []models.VipMembership
	var total int64

	query := r.db.Model(&models.VipMembership{}).Where("user_id = ?", userID)
	query.Count(&total)

	if page <= 0 { page = 1 }
	if size <= 0 { size = 20 }
	offset := (page - 1) * size

	err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&memberships).Error
	return memberships, total, err
}

func (r *paymentRepository) UpdateVipMembership(membership *models.VipMembership) error {
	return r.db.Save(membership).Error
}

func (r *paymentRepository) ExpireVipMembership(membershipID string) error {
	return r.db.Model(&models.VipMembership{}).
		Where("id = ?", membershipID).
		Update("is_active", false).Error
}

// Points
func (r *paymentRepository) CreatePointsRecord(record *models.PointsRecord) error {
	return r.db.Create(record).Error
}

func (r *paymentRepository) GetUserPointsRecords(userID string, page, size int) ([]models.PointsRecord, int64, error) {
	var records []models.PointsRecord
	var total int64

	query := r.db.Model(&models.PointsRecord{}).Where("user_id = ?", userID)
	query.Count(&total)

	if page <= 0 { page = 1 }
	if size <= 0 { size = 20 }
	offset := (page - 1) * size

	err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&records).Error
	return records, total, err
}

func (r *paymentRepository) GetUserPointsBalance(userID string) (int, error) {
	var earned, spent int64

	r.db.Model(&models.PointsRecord{}).
		Where("user_id = ? AND points_type = 'earn'", userID).
		Select("COALESCE(SUM(points), 0)").Scan(&earned)

	r.db.Model(&models.PointsRecord{}).
		Where("user_id = ? AND points_type = 'spend'", userID).
		Select("COALESCE(SUM(points), 0)").Scan(&spent)

	return int(earned - spent), nil
}

func (r *paymentRepository) GetPointsStats(userID string) (*models.PointsStatsResponse, error) {
	stats := &models.PointsStatsResponse{}

	// 总收入和支出
	var totalEarned, totalSpent int64
	r.db.Model(&models.PointsRecord{}).
		Where("user_id = ? AND points_type = 'earn'", userID).
		Select("COALESCE(SUM(points), 0)").Scan(&totalEarned)

	r.db.Model(&models.PointsRecord{}).
		Where("user_id = ? AND points_type = 'spend'", userID).
		Select("COALESCE(SUM(points), 0)").Scan(&totalSpent)

	stats.TotalEarned = int(totalEarned)
	stats.TotalSpent = int(totalSpent)
	stats.CurrentBalance = int(totalEarned - totalSpent)

	// 本月收入和支出
	startOfMonth := time.Now().Truncate(24 * time.Hour).AddDate(0, 0, -time.Now().Day()+1)
	var monthEarned, monthSpent int64

	r.db.Model(&models.PointsRecord{}).
		Where("user_id = ? AND points_type = 'earn' AND created_at >= ?", userID, startOfMonth).
		Select("COALESCE(SUM(points), 0)").Scan(&monthEarned)

	r.db.Model(&models.PointsRecord{}).
		Where("user_id = ? AND points_type = 'spend' AND created_at >= ?", userID, startOfMonth).
		Select("COALESCE(SUM(points), 0)").Scan(&monthSpent)

	stats.ThisMonthEarned = int(monthEarned)
	stats.ThisMonthSpent = int(monthSpent)

	return stats, nil
}

// Coins
func (r *paymentRepository) CreateCoinsRecord(record *models.CoinsRecord) error {
	return r.db.Create(record).Error
}

func (r *paymentRepository) GetUserCoinsRecords(userID string, page, size int) ([]models.CoinsRecord, int64, error) {
	var records []models.CoinsRecord
	var total int64

	query := r.db.Model(&models.CoinsRecord{}).Where("user_id = ?", userID)
	query.Count(&total)

	if page <= 0 { page = 1 }
	if size <= 0 { size = 20 }
	offset := (page - 1) * size

	err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&records).Error
	return records, total, err
}

func (r *paymentRepository) GetUserCoinsBalance(userID string) (int, error) {
	var earned, spent int64

	r.db.Model(&models.CoinsRecord{}).
		Where("user_id = ? AND coins_type = 'earn'", userID).
		Select("COALESCE(SUM(coins), 0)").Scan(&earned)

	r.db.Model(&models.CoinsRecord{}).
		Where("user_id = ? AND coins_type = 'spend'", userID).
		Select("COALESCE(SUM(coins), 0)").Scan(&spent)

	return int(earned - spent), nil
}

func (r *paymentRepository) GetCoinsStats(userID string) (*models.CoinsStatsResponse, error) {
	stats := &models.CoinsStatsResponse{}

	// 总收入和支出
	var totalEarned, totalSpent int64
	r.db.Model(&models.CoinsRecord{}).
		Where("user_id = ? AND coins_type = 'earn'", userID).
		Select("COALESCE(SUM(coins), 0)").Scan(&totalEarned)

	r.db.Model(&models.CoinsRecord{}).
		Where("user_id = ? AND coins_type = 'spend'", userID).
		Select("COALESCE(SUM(coins), 0)").Scan(&totalSpent)

	stats.TotalEarned = int(totalEarned)
	stats.TotalSpent = int(totalSpent)
	stats.CurrentBalance = int(totalEarned - totalSpent)

	// 本月收入和支出
	startOfMonth := time.Now().Truncate(24 * time.Hour).AddDate(0, 0, -time.Now().Day()+1)
	var monthEarned, monthSpent int64

	r.db.Model(&models.CoinsRecord{}).
		Where("user_id = ? AND coins_type = 'earn' AND created_at >= ?", userID, startOfMonth).
		Select("COALESCE(SUM(coins), 0)").Scan(&monthEarned)

	r.db.Model(&models.CoinsRecord{}).
		Where("user_id = ? AND coins_type = 'spend' AND created_at >= ?", userID, startOfMonth).
		Select("COALESCE(SUM(coins), 0)").Scan(&monthSpent)

	stats.ThisMonthEarned = int(monthEarned)
	stats.ThisMonthSpent = int(monthSpent)

	return stats, nil
}

// Checkin
func (r *paymentRepository) CreateCheckinRecord(record *models.CheckinRecord) error {
	return r.db.Create(record).Error
}

func (r *paymentRepository) GetLastCheckinRecord(userID string) (*models.CheckinRecord, error) {
	var record models.CheckinRecord
	err := r.db.Where("user_id = ?", userID).
		Order("checkin_date DESC").
		First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *paymentRepository) GetUserCheckinRecords(userID string, page, size int) ([]models.CheckinRecord, int64, error) {
	var records []models.CheckinRecord
	var total int64

	query := r.db.Model(&models.CheckinRecord{}).Where("user_id = ?", userID)
	query.Count(&total)

	if page <= 0 { page = 1 }
	if size <= 0 { size = 20 }
	offset := (page - 1) * size

	err := query.Order("checkin_date DESC").Offset(offset).Limit(size).Find(&records).Error
	return records, total, err
}

func (r *paymentRepository) GetCheckinStats(userID string) (int, error) {
	lastRecord, err := r.GetLastCheckinRecord(userID)
	if err != nil {
		return 0, nil
	}
	return lastRecord.ConsecutiveDays, nil
}

// Gifts
func (r *paymentRepository) CreateGift(gift *models.Gift) error {
	return r.db.Create(gift).Error
}

func (r *paymentRepository) GetGifts(category string, isActive *bool) ([]models.Gift, error) {
	var gifts []models.Gift
	query := r.db.Model(&models.Gift{})

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Order("created_at DESC").Find(&gifts).Error
	return gifts, err
}

func (r *paymentRepository) GetGiftByID(id string) (*models.Gift, error) {
	var gift models.Gift
	err := r.db.Where("id = ?", id).First(&gift).Error
	if err != nil {
		return nil, err
	}
	return &gift, nil
}

func (r *paymentRepository) UpdateGift(gift *models.Gift) error {
	return r.db.Save(gift).Error
}

func (r *paymentRepository) DeleteGift(id string) error {
	return r.db.Delete(&models.Gift{}, "id = ?", id).Error
}

// User Gifts
func (r *paymentRepository) CreateUserGift(userGift *models.UserGift) error {
	return r.db.Create(userGift).Error
}

func (r *paymentRepository) GetUserGifts(userID string, status string, page, size int) ([]models.UserGift, int64, error) {
	var gifts []models.UserGift
	var total int64

	query := r.db.Model(&models.UserGift{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	if page <= 0 { page = 1 }
	if size <= 0 { size = 20 }
	offset := (page - 1) * size

	err := query.Preload("Gift").Order("created_at DESC").Offset(offset).Limit(size).Find(&gifts).Error
	return gifts, total, err
}

func (r *paymentRepository) GetUserGiftByID(userID, giftID string) (*models.UserGift, error) {
	var userGift models.UserGift
	err := r.db.Where("user_id = ? AND id = ?", userID, giftID).
		Preload("Gift").First(&userGift).Error
	if err != nil {
		return nil, err
	}
	return &userGift, nil
}

func (r *paymentRepository) UpdateUserGiftStatus(userGiftID, status string, usedAt *time.Time) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if usedAt != nil {
		updates["used_at"] = usedAt
	}

	return r.db.Model(&models.UserGift{}).Where("id = ?", userGiftID).Updates(updates).Error
}

// Redeem Codes
func (r *paymentRepository) CreateRedeemCode(code *models.RedeemCode) error {
	return r.db.Create(code).Error
}

func (r *paymentRepository) GetRedeemCodeByCode(code string) (*models.RedeemCode, error) {
	var redeemCode models.RedeemCode
	err := r.db.Where("code = ?", code).Preload("Gift").First(&redeemCode).Error
	if err != nil {
		return nil, err
	}
	return &redeemCode, nil
}

func (r *paymentRepository) UseRedeemCode(codeID, userID string) error {
	now := time.Now()
	return r.db.Model(&models.RedeemCode{}).Where("id = ?", codeID).Updates(map[string]interface{}{
		"is_used": true,
		"used_by": userID,
		"used_at": &now,
	}).Error
}

func (r *paymentRepository) GetRedeemCodes(isUsed *bool, page, size int) ([]models.RedeemCode, int64, error) {
	var codes []models.RedeemCode
	var total int64

	query := r.db.Model(&models.RedeemCode{})

	if isUsed != nil {
		query = query.Where("is_used = ?", *isUsed)
	}

	query.Count(&total)

	if page <= 0 { page = 1 }
	if size <= 0 { size = 20 }
	offset := (page - 1) * size

	err := query.Preload("Gift").Order("created_at DESC").Offset(offset).Limit(size).Find(&codes).Error
	return codes, total, err
}

// Wallet
func (r *paymentRepository) GetUserWallet(userID string) (*models.WalletResponse, error) {
	wallet := &models.WalletResponse{
		UserID: userID,
	}

	// 获取积分余额
	pointsBalance, _ := r.GetUserPointsBalance(userID)
	wallet.Points = pointsBalance

	// 获取阅读币余额
	coinsBalance, _ := r.GetUserCoinsBalance(userID)
	wallet.ReadingCoins = coinsBalance

	// 获取VIP信息
	vipMembership, err := r.GetActiveVipMembership(userID)
	if err == nil {
		wallet.VipType = vipMembership.VipType
		expiresAt := vipMembership.EndDate.Format(time.RFC3339)
		wallet.VipExpiresAt = &expiresAt
	} else {
		wallet.VipType = "none"
	}

	// 获取连续签到天数
	consecutiveDays, _ := r.GetCheckinStats(userID)
	wallet.ConsecutiveCheckins = consecutiveDays

	return wallet, nil
}

func (r *paymentRepository) UpdateUserBalance(userID string, pointsDelta, coinsDelta int) error {
	// 这里应该与用户服务同步更新用户表中的余额
	// 由于这是微服务架构，实际实现中可能需要通过消息队列或API调用来同步
	return nil
}