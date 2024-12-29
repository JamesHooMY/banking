package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func RateLimitMiddleware(redisClient *redis.Client, limit int64, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing API key"})
			c.Abort()
			return
		}

		redisKey := fmt.Sprintf("rate_limit:%s", key)
		ctx := c.Request.Context()

		luaScript := redis.NewScript(`
			local limit = tonumber(ARGV[1])
			local ttl = ARGV[2]
			local current = redis.call("GET", KEYS[1])
			if not current then
				redis.call("SET", KEYS[1], limit - 1, "PX", ttl)
				return {limit - 1, ttl}
			end
			if tonumber(current) <= 0 then
				return {-1, redis.call("PTTL", KEYS[1])}
			end
			redis.call("DECR", KEYS[1])
			return {tonumber(current) - 1, redis.call("PTTL", KEYS[1])}
		`)

		result, err := luaScript.Run(ctx, redisClient, []string{redisKey}, limit, duration.Milliseconds()).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}
		fmt.Printf("result: %v\n", result)

		data := result.([]interface{})
		remaining, err := strconv.ParseInt(fmt.Sprintf("%v", data[0]), 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}
		reset, err := strconv.ParseInt(fmt.Sprintf("%v", data[1]), 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}

		// Rate limit exceeded
		if remaining < 0 {
			c.Header("Retry-After", fmt.Sprintf("%d", reset/1000))
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Duration(reset)*time.Millisecond).Unix()))
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Next()
	}
}
