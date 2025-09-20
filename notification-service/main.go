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
	"reading-microservices/notification-service/handlers"
	"reading-microservices/notification-service/models"
	"reading-microservices/notification-service/repositories"
	"reading-microservices/notification-service/services"
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
		&models.Notification{},
		&models.NotificationSetting{},
		&models.PushToken{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化Repository
	notificationRepo := repositories.NewNotificationRepository(db)

	// 初始化Service
	notificationService := services.NewNotificationService(notificationRepo)

	// 初始化Handler
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// 初始化路由
	router := setupRouter(notificationHandler, cfg.JWT.Secret)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logrus.Infof("Notification service starting on %s", addr)

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

func setupRouter(notificationHandler *handlers.NotificationHandler, jwtSecret string) *gin.Engine {
	router := gin.Default()

	// 中间件
	router.Use(middleware.CORS())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 健康检查
	router.GET("/health", notificationHandler.Health)

	// API路由 - 需要认证
	v1 := router.Group("/api/v1/notification")
	v1.Use(middleware.JWTAuth(jwtSecret))
	{
		// 通知管理
		v1.GET("/", notificationHandler.GetNotifications)
		v1.GET("/stats", notificationHandler.GetNotificationStats)
		v1.PUT("/:id/read", notificationHandler.MarkAsRead)
		v1.PUT("/read-all", notificationHandler.MarkAllAsRead)
		v1.DELETE("/:id", notificationHandler.DeleteNotification)

		// 通知设置
		v1.GET("/settings", notificationHandler.GetNotificationSettings)
		v1.PUT("/settings", notificationHandler.UpdateNotificationSetting)

		// 推送Token管理
		v1.POST("/push-token", notificationHandler.RegisterPushToken)
		v1.DELETE("/push-token/:device_id", notificationHandler.UnregisterPushToken)
	}

	// 内部API - 供其他服务调用
	internal := router.Group("/api/v1/internal/notification")
	{
		// 创建通知
		internal.POST("/", notificationHandler.CreateNotification)
		internal.POST("/batch", notificationHandler.CreateBatchNotifications)

		// 推送通知
		internal.POST("/push", notificationHandler.PushNotification)
	}

	return router
}