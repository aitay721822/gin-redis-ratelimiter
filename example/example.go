package main

import(
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/aitay721822/gin-redis-ratelimiter"
	"net/http"
)

func main() {

	router := gin.Default()

	router.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"message": "pong",
		})
	})

	router.GET("/ping2", func(context *gin.Context) {
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
