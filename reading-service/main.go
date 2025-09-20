package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"

	"reading-microservices/reading-service/handlers"
	"reading-microservices/reading-service/models"
	"reading-microservices/reading-service/repositories"
	"reading-microservices/reading-service/services"
	"reading-microservices/shared/config"
	"reading-microservices/shared/middleware"
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
		&models.ReadingRecord{},
		&models.Bookshelf{},
		&models.Favorite{},
		&models.Comment{},
		&models.SearchHistory{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化Repository
	readingRepo := repositories.NewReadingRepository(db)

	// 初始化Service
	readingService := services.NewReadingService(readingRepo)

	// 初始化Handler
	readingHandler := handlers.NewReadingHandler(readingService)

	// 初始化路由
	router := setupRouter(readingHandler, cfg.JWT.Secret)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logrus.Infof("Reading service starting on %s", addr)

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

func setupRouter(readingHandler *handlers.ReadingHandler, jwtSecret string) *gin.Engine {
	router := gin.Default()

	// 中间件
	router.Use(middleware.CORS())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 健康检查
	router.GET("/health", readingHandler.Health)

	// API路由 - 需要认证
	v1 := router.Group("/api/v1/reading")
	v1.Use(middleware.JWTAuth(jwtSecret))
	{
		// 阅读进度
		v1.POST("/progress", readingHandler.UpdateReadingProgress)
		v1.GET("/progress/:novel_id/:chapter_id", readingHandler.GetReadingProgress)
		v1.GET("/history", readingHandler.GetReadingHistory)

		// 书架管理
		v1.POST("/bookshelf", readingHandler.AddToBookshelf)
		v1.DELETE("/bookshelf/:novel_id", readingHandler.RemoveFromBookshelf)
		v1.GET("/bookshelf", readingHandler.GetBookshelf)
		v1.GET("/bookshelf/stats", readingHandler.GetBookshelfStats)

		// 收藏管理
		v1.POST("/favorites/:novel_id", readingHandler.AddFavorite)
		v1.DELETE("/favorites/:novel_id", readingHandler.RemoveFavorite)
		v1.GET("/favorites", readingHandler.GetFavorites)
		v1.GET("/favorites/:novel_id/status", readingHandler.CheckFavoriteStatus)

		// 评论系统
		v1.POST("/comments", readingHandler.CreateComment)
		v1.PUT("/comments/:id", readingHandler.UpdateComment)
		v1.DELETE("/comments/:id", readingHandler.DeleteComment)
		v1.GET("/comments/user", readingHandler.GetUserComments)
		v1.POST("/comments/:id/like", readingHandler.LikeComment)
		v1.DELETE("/comments/:id/like", readingHandler.UnlikeComment)

		// 搜索历史
		v1.POST("/search/history", readingHandler.AddSearchHistory)
		v1.GET("/search/history", readingHandler.GetSearchHistory)
		v1.DELETE("/search/history", readingHandler.ClearSearchHistory)

		// 统计信息
		v1.GET("/stats", readingHandler.GetReadingStats)
	}

	// 公开API - 不需要认证
	public := router.Group("/api/v1/reading/public")
	{
		// 小说评论
		public.GET("/novels/:novel_id/comments", readingHandler.GetNovelComments)
		// 章节评论
		public.GET("/chapters/:chapter_id/comments", readingHandler.GetChapterComments)
	}

	return router
}
