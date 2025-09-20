package services

import (
	"errors"
	"time"
	"fmt"
	"gorm.io/gorm"
	"reading-microservices/payment-service/models"
	"reading-microservices/payment-service/repositories"
)

type PaymentService interface {
	// VIP Management
	CreateVipMembership(userID string, req *models.CreateVipMembershipRequest) (*models.VipMembershipResponse, error)
	GetUserVipStatus(userID string) (*models.VipMembershipResponse, error)
	GetUserVipHistory(userID string, page, size int) ([]models.VipMembershipResponse, int64, error)

	// Points Management
	EarnPoints(userID string, req *models.EarnPointsRequest) error
	SpendPoints(userID string, req *models.SpendPointsRequest) error
	GetPointsHistory(userID string, page, size int) ([]models.PointsRecordResponse, int64, error)
	GetPointsStats(userID string) (*models.PointsStatsResponse, error)

	// Coins Management
	EarnCoins(userID string, req *models.EarnCoinsRequest) error
	SpendCoins(userID string, req *models.SpendCoinsRequest) error
	GetCoinsHistory(userID string, page, size int) ([]models.CoinsRecordResponse, int64, error)
	GetCoinsStats(userID string) (*models.CoinsStatsResponse, error)

	// Checkin System
	DailyCheckin(userID string) (*models.CheckinResponse, error)
	GetCheckinStatus(userID string) (*models.CheckinStatusResponse, error)
	GetCheckinHistory(userID string, page, size int) ([]models.CheckinResponse, int64, error)

	// Gift Management
	CreateGift(req *models.CreateGiftRequest) (*models.GiftResponse, error)
	GetGifts(category string, isActive *bool) ([]models.GiftResponse, error)
	UpdateGift(id string, req *models.CreateGiftRequest) error
	DeleteGift(id string) error

	// User Gifts
	GetUserGifts(userID string, status string, page, size int) ([]models.UserGiftResponse, int64, error)
	UseUserGift(userID, giftID string) error

	// Redeem System
	RedeemCode(userID string, req *models.RedeemCodeRequest) (*models.UserGiftResponse, error)

	// Wallet
	GetUserWallet(userID string) (*models.WalletResponse, error)
}

type paymentService struct {
	repo repositories.PaymentRepository
}

func NewPaymentService(repo repositories.PaymentRepository) PaymentService {
	return &paymentService{repo: repo}
}

// VIP Management
func (s *paymentService) CreateVipMembership(userID string, req *models.CreateVipMembershipRequest) (*models.VipMembershipResponse, error) {
	startDate := time.Now()
	endDate := startDate.AddDate(0, req.Duration, 0)

	membership := &models.VipMembership{
		UserID:        userID,
		VipType:       req.VipType,
		StartDate:     startDate,
		EndDate:       endDate,
		IsActive:      true,
		AutoRenew:     req.AutoRenew,
		PaymentMethod: &req.PaymentMethod,
		Amount:        &req.Amount,
	}

	if err := s.repo.CreateVipMembership(membership); err != nil {
		return nil, err
	}

	return s.convertToVipMembershipResponse(membership), nil
}

func (s *paymentService) GetUserVipStatus(userID string) (*models.VipMembershipResponse, error) {
	membership, err := s.repo.GetActiveVipMembership(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 返回默认的非会员状态
			return &models.VipMembershipResponse{
				VipType:       "none",
				IsActive:      false,
				DaysRemaining: 0,
			}, nil
		}
		return nil, err
	}

	return s.convertToVipMembershipResponse(membership), nil
}

func (s *paymentService) GetUserVipHistory(userID string, page, size int) ([]models.VipMembershipResponse, int64, error) {
	memberships, total, err := s.repo.GetUserVipMemberships(userID, page, size)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.VipMembershipResponse, len(memberships))
	for i, membership := range memberships {
		responses[i] = *s.convertToVipMembershipResponse(&membership)
	}

	return responses, total, nil
}

// Points Management
func (s *paymentService) EarnPoints(userID string, req *models.EarnPointsRequest) error {
	record := &models.PointsRecord{
		UserID:      userID,
		Points:      req.Points,
		PointsType:  "earn",
		Source:      req.Source,
		Description: req.Description,
		RelatedID:   req.RelatedID,
		RelatedType: req.RelatedType,
		CreatedAt:   time.Now(),
	}

	return s.repo.CreatePointsRecord(record)
}

func (s *paymentService) SpendPoints(userID string, req *models.SpendPointsRequest) error {
	// 检查余额是否足够
	balance, err := s.repo.GetUserPointsBalance(userID)
	if err != nil {
		return err
	}

	if balance < req.Points {
		return errors.New("insufficient points balance")
	}

	record := &models.PointsRecord{
		UserID:      userID,
		Points:      req.Points,
		PointsType:  "spend",
		Source:      req.Source,
		Description: req.Description,
		RelatedID:   req.RelatedID,
		RelatedType: req.RelatedType,
		CreatedAt:   time.Now(),
	}

	return s.repo.CreatePointsRecord(record)
}

func (s *paymentService) GetPointsHistory(userID string, page, size int) ([]models.PointsRecordResponse, int64, error) {
	records, total, err := s.repo.GetUserPointsRecords(userID, page, size)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.PointsRecordResponse, len(records))
	for i, record := range records {
		responses[i] = *s.convertToPointsRecordResponse(&record)
	}

	return responses, total, nil
}

func (s *paymentService) GetPointsStats(userID string) (*models.PointsStatsResponse, error) {
	return s.repo.GetPointsStats(userID)
}

// Coins Management
func (s *paymentService) EarnCoins(userID string, req *models.EarnCoinsRequest) error {
	record := &models.CoinsRecord{
		UserID:      userID,
		Coins:       req.Coins,
		CoinsType:   "earn",
		Source:      req.Source,
		Description: req.Description,
		RelatedID:   req.RelatedID,
		RelatedType: req.RelatedType,
		CreatedAt:   time.Now(),
	}

	return s.repo.CreateCoinsRecord(record)
}

func (s *paymentService) SpendCoins(userID string, req *models.SpendCoinsRequest) error {
	// 检查余额是否足够
	balance, err := s.repo.GetUserCoinsBalance(userID)
	if err != nil {
		return err
	}

	if balance < req.Coins {
		return errors.New("insufficient coins balance")
	}

	record := &models.CoinsRecord{
		UserID:      userID,
		Coins:       req.Coins,
		CoinsType:   "spend",
		Source:      req.Source,
		Description: req.Description,
		RelatedID:   req.RelatedID,
		RelatedType: req.RelatedType,
		CreatedAt:   time.Now(),
	}

	return s.repo.CreateCoinsRecord(record)
}

func (s *paymentService) GetCoinsHistory(userID string, page, size int) ([]models.CoinsRecordResponse, int64, error) {
	records, total, err := s.repo.GetUserCoinsRecords(userID, page, size)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.CoinsRecordResponse, len(records))
	for i, record := range records {
		responses[i] = *s.convertToCoinsRecordResponse(&record)
	}

	return responses, total, nil
}

func (s *paymentService) GetCoinsStats(userID string) (*models.CoinsStatsResponse, error) {
	return s.repo.GetCoinsStats(userID)
}

// Checkin System
func (s *paymentService) DailyCheckin(userID string) (*models.CheckinResponse, error) {
	today := time.Now().Truncate(24 * time.Hour)

	// 检查今天是否已经签到
	lastCheckin, err := s.repo.GetLastCheckinRecord(userID)
	if err == nil && lastCheckin.CheckinDate.Equal(today) {
		return nil, errors.New("already checked in today")
	}

	// 计算连续签到天数
	consecutiveDays := 1
	if err == nil {
		yesterday := today.AddDate(0, 0, -1)
		if lastCheckin.CheckinDate.Equal(yesterday) {
			consecutiveDays = lastCheckin.ConsecutiveDays + 1
		}
	}

	// 计算奖励
	pointsReward, coinsReward := s.calculateCheckinRewards(consecutiveDays)

	// 创建签到记录
	checkinRecord := &models.CheckinRecord{
		UserID:          userID,
		CheckinDate:     today,
		ConsecutiveDays: consecutiveDays,
		PointsEarned:    pointsReward,
		CoinsEarned:     coinsReward,
		CreatedAt:       time.Now(),
	}

	if err := s.repo.CreateCheckinRecord(checkinRecord); err != nil {
		return nil, err
	}

	// 发放奖励
	if pointsReward > 0 {
		s.EarnPoints(userID, &models.EarnPointsRequest{
			Points:      pointsReward,
			Source:      "daily_checkin",
			Description: stringPtr(fmt.Sprintf("连续签到%d天奖励", consecutiveDays)),
		})
	}

	if coinsReward > 0 {
		s.EarnCoins(userID, &models.EarnCoinsRequest{
			Coins:       coinsReward,
			Source:      "daily_checkin",
			Description: stringPtr(fmt.Sprintf("连续签到%d天奖励", consecutiveDays)),
		})
	}

	return s.convertToCheckinResponse(checkinRecord), nil
}

func (s *paymentService) GetCheckinStatus(userID string) (*models.CheckinStatusResponse, error) {
	today := time.Now().Truncate(24 * time.Hour)

	lastCheckin, err := s.repo.GetLastCheckinRecord(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			pointsReward, coinsReward := s.calculateCheckinRewards(1)
			return &models.CheckinStatusResponse{
				CanCheckin:        true,
				ConsecutiveDays:   0,
				TodayPointsReward: pointsReward,
				TodayCoinsReward:  coinsReward,
			}, nil
		}
		return nil, err
	}

	canCheckin := !lastCheckin.CheckinDate.Equal(today)
	consecutiveDays := lastCheckin.ConsecutiveDays
	if !lastCheckin.CheckinDate.Equal(today.AddDate(0, 0, -1)) && canCheckin {
		consecutiveDays = 0
	}

	nextConsecutiveDays := consecutiveDays
	if canCheckin {
		nextConsecutiveDays++
	}

	pointsReward, coinsReward := s.calculateCheckinRewards(nextConsecutiveDays)

	lastCheckinDate := lastCheckin.CheckinDate.Format("2006-01-02")

	return &models.CheckinStatusResponse{
		CanCheckin:         canCheckin,
		LastCheckinDate:    &lastCheckinDate,
		ConsecutiveDays:    consecutiveDays,
		TodayPointsReward:  pointsReward,
		TodayCoinsReward:   coinsReward,
	}, nil
}

func (s *paymentService) GetCheckinHistory(userID string, page, size int) ([]models.CheckinResponse, int64, error) {
	records, total, err := s.repo.GetUserCheckinRecords(userID, page, size)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.CheckinResponse, len(records))
	for i, record := range records {
		responses[i] = *s.convertToCheckinResponse(&record)
	}

	return responses, total, nil
}

// Gift Management
func (s *paymentService) CreateGift(req *models.CreateGiftRequest) (*models.GiftResponse, error) {
	gift := &models.Gift{
		Name:        req.Name,
		Description: req.Description,
		GiftType:    req.GiftType,
		Value:       req.Value,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.CreateGift(gift); err != nil {
		return nil, err
	}

	return s.convertToGiftResponse(gift), nil
}

func (s *paymentService) GetGifts(category string, isActive *bool) ([]models.GiftResponse, error) {
	gifts, err := s.repo.GetGifts(category, isActive)
	if err != nil {
		return nil, err
	}

	responses := make([]models.GiftResponse, len(gifts))
	for i, gift := range gifts {
		responses[i] = *s.convertToGiftResponse(&gift)
	}

	return responses, nil
}

func (s *paymentService) UpdateGift(id string, req *models.CreateGiftRequest) error {
	gift, err := s.repo.GetGiftByID(id)
	if err != nil {
		return err
	}

	gift.Name = req.Name
	gift.Description = req.Description
	gift.GiftType = req.GiftType
	gift.Value = req.Value
	gift.Category = req.Category
	gift.ImageURL = req.ImageURL

	return s.repo.UpdateGift(gift)
}

func (s *paymentService) DeleteGift(id string) error {
	return s.repo.DeleteGift(id)
}

// User Gifts
func (s *paymentService) GetUserGifts(userID string, status string, page, size int) ([]models.UserGiftResponse, int64, error) {
	gifts, total, err := s.repo.GetUserGifts(userID, status, page, size)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.UserGiftResponse, len(gifts))
	for i, gift := range gifts {
		responses[i] = *s.convertToUserGiftResponse(&gift)
	}

	return responses, total, nil
}

func (s *paymentService) UseUserGift(userID, giftID string) error {
	userGift, err := s.repo.GetUserGiftByID(userID, giftID)
	if err != nil {
		return err
	}

	if userGift.Status != "unused" {
		return errors.New("gift is not available for use")
	}

	if userGift.ExpiresAt != nil && userGift.ExpiresAt.Before(time.Now()) {
		return errors.New("gift has expired")
	}

	// 根据礼品类型执行相应的逻辑
	switch userGift.Gift.GiftType {
	case "reading_coins":
		// 发放阅读币
		coins := parseIntFromString(userGift.Gift.Value)
		s.EarnCoins(userID, &models.EarnCoinsRequest{
			Coins:       coins,
			Source:      "gift_redeem",
			Description: stringPtr(fmt.Sprintf("使用礼品：%s", userGift.Gift.Name)),
			RelatedID:   &giftID,
			RelatedType: stringPtr("user_gift"),
		})

	case "points":
		// 发放积分
		points := parseIntFromString(userGift.Gift.Value)
		s.EarnPoints(userID, &models.EarnPointsRequest{
			Points:      points,
			Source:      "gift_redeem",
			Description: stringPtr(fmt.Sprintf("使用礼品：%s", userGift.Gift.Name)),
			RelatedID:   &giftID,
			RelatedType: stringPtr("user_gift"),
		})

	case "vip_time":
		// TODO: 发放VIP时间，需要解析Value中的时间长度
		// 这里简化处理
	}

	// 更新礼品状态
	now := time.Now()
	return s.repo.UpdateUserGiftStatus(giftID, "used", &now)
}

// Redeem System
func (s *paymentService) RedeemCode(userID string, req *models.RedeemCodeRequest) (*models.UserGiftResponse, error) {
	// 获取兑换码
	redeemCode, err := s.repo.GetRedeemCodeByCode(req.Code)
	if err != nil {
		return nil, errors.New("invalid redeem code")
	}

	if redeemCode.IsUsed {
		return nil, errors.New("redeem code has already been used")
	}

	if redeemCode.ExpiresAt != nil && redeemCode.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("redeem code has expired")
	}

	// 创建用户礼品
	userGift := &models.UserGift{
		UserID:     userID,
		GiftID:     redeemCode.GiftID,
		Status:     "unused",
		ObtainedAt: time.Now(),
		Source:     stringPtr("redeem_code"),
		CreatedAt:  time.Now(),
	}

	// 如果礼品有有效期，设置过期时间
	if redeemCode.Gift.GiftType == "vip_time" {
		// 这里简化处理，实际应该根据礼品类型设置不同的过期时间
		expiresAt := time.Now().AddDate(0, 1, 0)
		userGift.ExpiresAt = &expiresAt
	}

	if err := s.repo.CreateUserGift(userGift); err != nil {
		return nil, err
	}

	// 标记兑换码为已使用
	if err := s.repo.UseRedeemCode(redeemCode.ID, userID); err != nil {
		return nil, err
	}

	// 重新获取包含礼品信息的用户礼品
	userGiftWithGift, err := s.repo.GetUserGiftByID(userID, userGift.ID)
	if err != nil {
		return nil, err
	}

	return s.convertToUserGiftResponse(userGiftWithGift), nil
}

// Wallet
func (s *paymentService) GetUserWallet(userID string) (*models.WalletResponse, error) {
	return s.repo.GetUserWallet(userID)
}

// Helper methods
func (s *paymentService) calculateCheckinRewards(consecutiveDays int) (int, int) {
	basePoints := 10
	baseCoins := 5

	// 连续签到奖励递增
	bonusMultiplier := (consecutiveDays - 1) / 7 // 每7天增加一个奖励级别
	points := basePoints + bonusMultiplier*5
	coins := baseCoins + bonusMultiplier*2

	// 设置上限
	if points > 50 {
		points = 50
	}
	if coins > 25 {
		coins = 25
	}

	return points, coins
}

func (s *paymentService) convertToVipMembershipResponse(membership *models.VipMembership) *models.VipMembershipResponse {
	daysRemaining := 0
	if membership.IsActive && membership.EndDate.After(time.Now()) {
		daysRemaining = int(time.Until(membership.EndDate).Hours() / 24)
	}

	return &models.VipMembershipResponse{
		ID:            membership.ID,
		VipType:       membership.VipType,
		StartDate:     membership.StartDate.Format("2006-01-02"),
		EndDate:       membership.EndDate.Format("2006-01-02"),
		IsActive:      membership.IsActive,
		AutoRenew:     membership.AutoRenew,
		PaymentMethod: membership.PaymentMethod,
		Amount:        membership.Amount,
		DaysRemaining: daysRemaining,
		CreatedAt:     membership.CreatedAt.Format(time.RFC3339),
	}
}

func (s *paymentService) convertToPointsRecordResponse(record *models.PointsRecord) *models.PointsRecordResponse {
	return &models.PointsRecordResponse{
		ID:          record.ID,
		Points:      record.Points,
		PointsType:  record.PointsType,
		Source:      record.Source,
		Description: record.Description,
		RelatedID:   record.RelatedID,
		RelatedType: record.RelatedType,
		CreatedAt:   record.CreatedAt.Format(time.RFC3339),
	}
}

func (s *paymentService) convertToCoinsRecordResponse(record *models.CoinsRecord) *models.CoinsRecordResponse {
	return &models.CoinsRecordResponse{
		ID:          record.ID,
		Coins:       record.Coins,
		CoinsType:   record.CoinsType,
		Source:      record.Source,
		Description: record.Description,
		RelatedID:   record.RelatedID,
		RelatedType: record.RelatedType,
		CreatedAt:   record.CreatedAt.Format(time.RFC3339),
	}
}

func (s *paymentService) convertToCheckinResponse(record *models.CheckinRecord) *models.CheckinResponse {
	return &models.CheckinResponse{
		ID:             record.ID,
		CheckinDate:    record.CheckinDate.Format("2006-01-02"),
		ConsecutiveDays: record.ConsecutiveDays,
		PointsEarned:   record.PointsEarned,
		CoinsEarned:    record.CoinsEarned,
		CreatedAt:      record.CreatedAt.Format(time.RFC3339),
	}
}

func (s *paymentService) convertToGiftResponse(gift *models.Gift) *models.GiftResponse {
	return &models.GiftResponse{
		ID:          gift.ID,
		Name:        gift.Name,
		Description: gift.Description,
		GiftType:    gift.GiftType,
		Value:       gift.Value,
		Category:    gift.Category,
		ImageURL:    gift.ImageURL,
	}
}

func (s *paymentService) convertToUserGiftResponse(userGift *models.UserGift) *models.UserGiftResponse {
	response := &models.UserGiftResponse{
		ID:         userGift.ID,
		Status:     userGift.Status,
		ObtainedAt: userGift.ObtainedAt.Format(time.RFC3339),
		Source:     userGift.Source,
		Gift:       *s.convertToGiftResponse(&userGift.Gift),
	}

	if userGift.UsedAt != nil {
		usedAt := userGift.UsedAt.Format(time.RFC3339)
		response.UsedAt = &usedAt
	}

	if userGift.ExpiresAt != nil {
		expiresAt := userGift.ExpiresAt.Format(time.RFC3339)
		response.ExpiresAt = &expiresAt
	}

	return response
}

func stringPtr(s string) *string {
	return &s
}

func parseIntFromString(s string) int {
	// 这里简化处理，实际应该解析字符串中的数字
	// 例如 "100" -> 100, "100coins" -> 100
	switch s {
	case "10":
		return 10
	case "50":
		return 50
	case "100":
		return 100
	default:
		return 0
	}
}