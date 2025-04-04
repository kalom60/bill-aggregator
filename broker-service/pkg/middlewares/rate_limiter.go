package middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	address     = os.Getenv("REDIS_ADDRESS")
	password    = os.Getenv("REDIS_PASSWORD")
	ctx         = context.Background()
	redisClient *redis.Client
)

func InitRedis() {
	fullAddress := fmt.Sprintf("%s:%s", address, "6379")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fullAddress,
		Password: password,
		DB:       0,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

func RateLimiterMiddleware(limit int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", userID)

		count, err := redisClient.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limiter error"})
			c.Abort()
			return
		}

		if count >= limit {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		pipe := redisClient.TxPipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, duration)
		_, _ = pipe.Exec(ctx)

		c.Next()
	}
}
