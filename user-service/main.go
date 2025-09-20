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
	"reading-microservices/shared/config"
	"reading-microservices/shared/middleware"
	"reading-microservices/user-service/handlers"
	middleware2 "reading-microservices/user-service/middleware"
	"reading-microservices/user-service/models"
	"reading-microservices/user-service/repositories"
	"reading-microservices/user-service/services"
	authServices "reading-microservices/user-service/services/auth"
	"time"
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

	// 自动迁移
	if err := db.AutoMigrate(
		&models.User{},
		&models.ThirdPartyAccount{},
		&models.UserSession{},
		&models.LoginLog{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化Repository
	userRepo := repositories.NewUserRepository(db)

	// 初始化 AuthManager
	// 初始化 AuthManager（不再需要默认过期时间，因为我们会分别设置）
	authManager := authServices.NewAuthManager(cfg.JWT.Secret)

	// 初始化 SessionManager
	sessionManager := authServices.NewSessionManager(userRepo, rdb)

	// 初始化 LoginLogger
	loginLogger := authServices.NewLoginLogger(userRepo)

	// 配置双 token 的过期时间
	accessExpiresIn := 7200    // 例如：7200 (2小时)
	refreshExpiresIn := 604800 // 例如：604800 (7天)

	// 初始化 UserService，传入双 token 的过期时间配置
	var userService services.UserServiceInterface
	userService = services.NewUserService(
		userRepo,
		authManager,
		sessionManager,
		loginLogger,
		accessExpiresIn,
		refreshExpiresIn,
	)
	// 初始化 Handler
	userHandler := handlers.NewUserHandler(userService)

	// 初始化路由
	router := setupRouter(userHandler, cfg.JWT.Secret, rdb)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logrus.Infof("User service starting on %s", addr)

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
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func initRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: 20,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func setupRouter(userHandler *handlers.UserHandler, jwtSecret string, rdb *redis.Client) *gin.Engine {
	router := gin.Default()

	// 中间件
	router.Use(middleware.CORS())

	router.Use(middleware2.SecurityHeaders())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 健康检查
	router.GET("/health", userHandler.Health)
	router.GET("/healthz", userHandler.Health) // 添加另一个健康检查端点

	// API 路由
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// 公开认证接口
			auth := v1.Group("/auth")
			auth.Use(middleware2.RateLimit(rdb, 10, time.Minute)) // 登录注册更高
			{
				auth.POST("/register", userHandler.Register)
				auth.POST("/login", userHandler.Login)
				auth.POST("/refresh", userHandler.RefreshToken)
				auth.POST("/validate", userHandler.ValidateToken)
			}

			// 需要认证的用户接口
			user := v1.Group("/user")
			user.Use(middleware.JWTAuth(jwtSecret))
			user.Use(middleware2.RateLimit(rdb, 100, time.Minute)) // 增加限制到每分钟100次
			{
				user.GET("/profile", userHandler.GetProfile)
				user.PUT("/profile", userHandler.UpdateProfile)
				user.POST("/change-password", userHandler.ChangePassword)
				user.POST("/logout", userHandler.Logout)
			}
		}
	}

	return router
}
