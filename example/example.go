package main

import (
	"context"
	"fmt"
	"github.com/aitay721822/gin-redis-ratelimiter"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"time"
)

var ctx = context.Background()

func main() {

	router := gin.Default()

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	lm1 := ginredisratelimiter.NewRedisRateLimiter(ctx, "ping1", time.Minute, 10, client)
	lm2 := ginredisratelimiter.NewRedisRateLimiter(ctx, "ping2", 2 * time.Minute, 5, client)

	router.GET("/ping", lm1.Middleware(), func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"message": "pong",
		})
	})

	router.GET("/ping2", lm2.Middleware(), func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"message": "pong pong",
		})
	})

	err := router.Run(":8080")
	if err != nil {
		fmt.Printf("error occurred during initialization gin: %v", err)
	}
}
