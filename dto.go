package ginredisratelimiter

import (
	"context"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

type LimiterResponse struct {
	status bool
	remain int64
	err    error
}

type Limiter interface {
	Take() *LimiterResponse
}

type RedisRateLimiter struct {
	sync.Mutex
	context			context.Context
	scriptSHA1  	string
	client 			*redis.Client
}

type TokenBucketRedisRateLimiter struct {
	RedisRateLimiter

	identifier      string
	interval		time.Duration
	maxRequest		int
}

func (r *TokenBucketRedisRateLimiter) Take(request TokenBucketLuaRequest) *LimiterResponse {
	r.Lock()
	defer r.Unlock()

	result, err := r.client.EvalSha(
		r.context,
		r.scriptSHA1,
		[]string{request.valueKey, request.timestampKey},
		request.limit, request.interval, request.batchSize,
	).Result()

	if err != nil {
		return &LimiterResponse {
			status: false,
			remain: 0,
			err: err,
		}
	} else {
		data := result.([]interface{})
		if len(data) != 2 {
			return &LimiterResponse{
				status: false,
				remain: 0,
				err:    ErrRedisError,
			}
		}
		return &LimiterResponse{
			status: data[0] == nil,
			remain: data[1].(int64),
			err:    nil,
		}
	}
}