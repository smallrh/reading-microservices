package main

import (
	"fmt"
	"log"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"reading-microservices/shared/config"
	"reading-microservices/shared/middleware"
	"reading-microservices/payment-service/handlers"
	"reading-microservices/payment-service/models"
	"reading-microservices/payment-service/repositories"
	"reading-microservices/payment-service/services"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化日志
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// 连接数据库
	db, err := initDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// 连接Redis
	rdb, err := initRedis(cfg)
	if err != nil {
		log.Fatal("Failed to connect redis:", err)
	}
	_ = rdb // Redis client available for future use

	// 自动迁移
	if err := db.AutoMigrate(
		&models.VipMembership{},
		&models.PointsRecord{},
		&models.CoinsRecord{},
		&models.CheckinRecord{},
		&models.Gift{},
		&models.UserGift{},
		&models.RedeemCode{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化Repository
	paymentRepo := repositories.NewPaymentRepository(db)

	// 初始化Service
	paymentService := services.NewPaymentService(paymentRepo)

	// 初始化Handler
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// 初始化路由
	router := setupRouter(paymentHandler, cfg.JWT.Secret)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logrus.Infof("Payment service starting on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database,
		cfg.Database.Charset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}

func initRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 测试连接
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func setupRouter(paymentHandler *handlers.PaymentHandler, jwtSecret string) *gin.Engine {
	router := gin.Default()

	// 中间件
	router.Use(middleware.CORS())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 健康检查
	router.GET("/health", paymentHandler.Health)

	// API路由 - 需要认证
	v1 := router.Group("/api/v1/payment")
	v1.Use(middleware.JWTAuth(jwtSecret))
	{
		// VIP管理
		v1.POST("/vip", paymentHandler.CreateVipMembership)
		v1.GET("/vip/status", paymentHandler.GetVipStatus)
		v1.GET("/vip/history", paymentHandler.GetVipHistory)

		// 积分管理
		v1.POST("/points/earn", paymentHandler.EarnPoints)
		v1.POST("/points/spend", paymentHandler.SpendPoints)
		v1.GET("/points/history", paymentHandler.GetPointsHistory)
		v1.GET("/points/stats", paymentHandler.GetPointsStats)

		// 阅读币管理
		v1.POST("/coins/earn", paymentHandler.EarnCoins)
		v1.POST("/coins/spend", paymentHandler.SpendCoins)
		v1.GET("/coins/history", paymentHandler.GetCoinsHistory)
		v1.GET("/coins/stats", paymentHandler.GetCoinsStats)

		// 签到系统
		v1.POST("/checkin", paymentHandler.DailyCheckin)
		v1.GET("/checkin/status", paymentHandler.GetCheckinStatus)
		v1.GET("/checkin/history", paymentHandler.GetCheckinHistory)

		// 用户礼品
		v1.GET("/gifts", paymentHandler.GetUserGifts)
		v1.POST("/gifts/:id/use", paymentHandler.UseUserGift)

		// 兑换码
		v1.POST("/redeem", paymentHandler.RedeemCode)

		// 钱包
		v1.GET("/wallet", paymentHandler.GetWallet)
	}

	// 管理接口
	admin := router.Group("/api/v1/admin/payment")
	admin.Use(middleware.JWTAuth(jwtSecret))
	{
		// 礼品管理
		admin.POST("/gifts", paymentHandler.CreateGift)
		admin.GET("/gifts", paymentHandler.GetGifts)
		admin.PUT("/gifts/:id", paymentHandler.UpdateGift)
		admin.DELETE("/gifts/:id", paymentHandler.DeleteGift)
	}

	return router
}