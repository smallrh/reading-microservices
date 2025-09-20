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

	"reading-microservices/content-service/handlers"
	"reading-microservices/content-service/models"
	"reading-microservices/content-service/repositories"
	"reading-microservices/content-service/services"
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
		&models.Category{},
		&models.Tag{},
		&models.Novel{},
		&models.NovelTag{},
		&models.Chapter{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化Repository
	contentRepo := repositories.NewContentRepository(db)

	// 初始化Service
	contentService := services.NewContentService(contentRepo)

	// 初始化Handler
	contentHandler := handlers.NewContentHandler(contentService)

	// 初始化路由
	router := setupRouter(contentHandler, cfg.JWT.Secret)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logrus.Infof("Content service starting on %s", addr)

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

func setupRouter(contentHandler *handlers.ContentHandler, jwtSecret string) *gin.Engine {
	router := gin.Default()

	// 中间件
	router.Use(middleware.CORS())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 健康检查
	router.GET("/health", contentHandler.Health)

	// API路由
	v1 := router.Group("/api/v1")
	{
		// 公开API - 不需要认证
		public := v1.Group("/content")
		{
			// 分类
			public.GET("/categories", contentHandler.GetCategories)
			public.GET("/categories/:id", contentHandler.GetCategoryByID)
			public.GET("/categories/:id/novels", contentHandler.GetNovelsByCategory)

			// 标签
			public.GET("/tags", contentHandler.GetTags)
			public.GET("/tags/:id", contentHandler.GetTagByID)

			// 小说
			public.GET("/novels/search", contentHandler.SearchNovels)
			public.GET("/novels/featured", contentHandler.GetFeaturedNovels)
			public.GET("/novels/latest", contentHandler.GetLatestNovels)
			public.GET("/novels/:novel_id", contentHandler.GetNovelByID)
			
			// 章节
			public.GET("/novels/:novel_id/chapters", contentHandler.GetChaptersByNovel)
			public.GET("/novels/:novel_id/chapters/:chapter_number", contentHandler.GetChapterByNumber)
			public.GET("/chapters/:id", contentHandler.GetChapterByID)

		}

		// 管理API - 需要认证
		admin := v1.Group("/admin/content")
		admin.Use(middleware.JWTAuth(jwtSecret))
		{
			// 分类管理
			admin.POST("/categories", contentHandler.CreateCategory)
			admin.PUT("/categories/:id", contentHandler.UpdateCategory)
			admin.DELETE("/categories/:id", contentHandler.DeleteCategory)

			// 标签管理
			admin.POST("/tags", contentHandler.CreateTag)

			// 小说管理
			admin.POST("/novels", contentHandler.CreateNovel)
			admin.PUT("/novels/:id", contentHandler.UpdateNovel)
			admin.DELETE("/novels/:id", contentHandler.DeleteNovel)

			// 章节管理
			admin.POST("/chapters", contentHandler.CreateChapter)
			admin.PUT("/chapters/:id", contentHandler.UpdateChapter)
			admin.DELETE("/chapters/:id", contentHandler.DeleteChapter)
		}
	}

	return router
}
