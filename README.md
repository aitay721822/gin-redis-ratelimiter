# gin-redis-ratelimiter

This is an **example** about using Redis + Token Bucket algorithm for Rate Limiting on Gin Framework.

# Usage

```go
 client := redis.NewClient(&redis.Options{
  Addr: "localhost:6379",
  Password: "",
  DB: 0,
 })

lm := ginredisratelimiter.NewRedisRateLimiter(ctx, "ping", time.Minute, 10, client)

r.GET("/ping", lm.Middleware(), func(c *gin.Context) {
  c.JSON(200, gin.H{
   "message": "pong",
  })
})
```
That means the URI `/ping` only allow 10 request per minute by IP Address.