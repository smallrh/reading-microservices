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
	"reading-microservices/download-service/handlers"
	"reading-microservices/download-service/models"
	"reading-microservices/download-service/repositories"
	"reading-microservices/download-service/services"
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
		&models.DownloadTask{},
		&models.DownloadChapter{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化Repository
	downloadRepo := repositories.NewDownloadRepository(db)

	// 初始化Service
	downloadService := services.NewDownloadService(downloadRepo)

	// 初始化Handler
	downloadHandler := handlers.NewDownloadHandler(downloadService)

	// 初始化路由
	router := setupRouter(downloadHandler, cfg.JWT.Secret)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logrus.Infof("Download service starting on %s", addr)

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

func setupRouter(downloadHandler *handlers.DownloadHandler, jwtSecret string) *gin.Engine {
	router := gin.Default()

	// 中间件
	router.Use(middleware.CORS())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 健康检查
	router.GET("/health", downloadHandler.Health)

	// API路由 - 需要认证
	v1 := router.Group("/api/v1/download")
	v1.Use(middleware.JWTAuth(jwtSecret))
	{
		// 下载任务管理
		v1.POST("/tasks", downloadHandler.CreateDownloadTask)
		v1.GET("/tasks", downloadHandler.GetDownloadTasks)
		v1.GET("/tasks/:id", downloadHandler.GetDownloadTask)
		v1.PUT("/tasks/:id", downloadHandler.UpdateDownloadTask)
		v1.DELETE("/tasks/:id", downloadHandler.DeleteDownloadTask)
		v1.POST("/tasks/:id/start", downloadHandler.StartDownload)
		v1.POST("/tasks/:id/pause", downloadHandler.PauseDownload)
		v1.POST("/tasks/:id/resume", downloadHandler.ResumeDownload)

		// 文件下载
		v1.GET("/tasks/:id/file", downloadHandler.DownloadFile)

		// 统计信息
		v1.GET("/stats", downloadHandler.GetDownloadStats)
	}

	return router
}