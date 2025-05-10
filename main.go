package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/itsanindyak/go-jwt/config"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	port := config.PORT

	if port == "" {
		port = "8000"
	}

	// Create router from gin
	router := gin.New()

	// Trust only localhost
	err := router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Fatalf("Error setting trusted proxies: %v", err)
	}

	// add middleware for logging
	router.Use(gin.Logger())

	// Health testing route
	router.GET("/api/v1/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"Success": "Access granted for /api/v1"})
	})

	router.Run(":" + port)
}
