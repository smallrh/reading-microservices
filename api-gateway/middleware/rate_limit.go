package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"time"
)

type RateLimiter struct {
	rdb               *redis.Client
	requestsPerMinute int
	burst             int
}

func NewRateLimiter(rdb *redis.Client, rpm, burst int) *RateLimiter {
	return &RateLimiter{rdb: rdb, requestsPerMinute: rpm, burst: burst}
}

func (rl *RateLimiter) IPLimit(requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rl.rdb == nil {
			c.Next()
			return
		}
		ip := c.ClientIP()
		key := fmt.Sprintf("ip_rate:%s", ip)
		ctx := context.Background()
		now := time.Now().Unix()
		window := int64(60)

		rl.rdb.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-window))
		count, _ := rl.rdb.ZCard(ctx, key).Result()
		if int(count) >= requestsPerMinute {
			c.JSON(429, gin.H{"code": 429, "message": "IP rate limit exceeded"})
			c.Abort()
			return
		}
		rl.rdb.ZAdd(ctx, key, &redis.Z{Score: float64(now), Member: now})
		rl.rdb.Expire(ctx, key, time.Duration(window)*time.Second)
		c.Next()
	}
}

func (rl *RateLimiter) UserLimit(requestsPerHour int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rl.rdb == nil {
			c.Next()
			return
		}
		userID := c.GetString("user_id")
		if userID == "" {
			c.Next()
			return
		}
		key := fmt.Sprintf("user_rate:%s", userID)
		ctx := context.Background()
		now := time.Now().Unix()
		window := int64(3600)

		rl.rdb.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-window))
		count, _ := rl.rdb.ZCard(ctx, key).Result()
		if int(count) >= requestsPerHour {
			c.JSON(429, gin.H{"code": 429, "message": "User rate limit exceeded"})
			c.Abort()
			return
		}
		rl.rdb.ZAdd(ctx, key, &redis.Z{Score: float64(now), Member: now})
		rl.rdb.Expire(ctx, key, time.Duration(window)*time.Second)
		c.Next()
	}
}
