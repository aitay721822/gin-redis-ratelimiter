package ginredisratelimiter

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strconv"
	"time"
)

func (r *TokenBucketRedisRateLimiter) Middleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		ip := context.ClientIP()
		if ip == "" {
			_ = context.AbortWithError(http.StatusInternalServerError, ErrIpNotRecognize)
		}
		request := TokenBucketLuaRequest{
			valueKey: 	  fmt.Sprintf("%v_%v_Token", r.identifier, ip),
			timestampKey: fmt.Sprintf("%v_%v_Update_Time", r.identifier, ip),
			limit: 		  int64(r.maxRequest),
			interval:     r.interval.Milliseconds(),
			batchSize: 	  1,
		}
		response := r.Take(request)
		if response.status {
			context.Writer.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(response.remain, 10))
			context.Writer.Header().Set("X-RateLimit-Limit", strconv.Itoa(r.maxRequest))
			context.Next()
		} else {
			_ = context.AbortWithError(http.StatusTooManyRequests, TooManyRequest)
		}
	}
}
func NewRedisRateLimiter(ctx context.Context, identifier string,interval time.Duration, times int, redisClient *redis.Client) *TokenBucketRedisRateLimiter {
	script := TokenBucketLuaScript
	scriptSHA1 := fmt.Sprintf("%x", sha1.Sum([]byte(script)))

	if !redisClient.ScriptExists(ctx, scriptSHA1).Val()[0] {
		redisClient.ScriptLoad(ctx, script).Val()
	}

	return &TokenBucketRedisRateLimiter{
		RedisRateLimiter: RedisRateLimiter{
			context:    ctx,
			scriptSHA1: scriptSHA1,
			client:     redisClient,
		},
		identifier:		identifier,
		interval:       interval,
		maxRequest:     times,
	}
}