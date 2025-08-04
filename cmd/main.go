package main

import (
	"github.com/gin-gonic/gin"
	"github.com/itsanindyak/go-jwt/config"
	"github.com/itsanindyak/go-jwt/pkg/logger"
	routes "github.com/itsanindyak/go-jwt/routes"
)

func main() {

	env := config.ENV
	// $env:ENV="production"
	switch env {
	case "production":
		gin.SetMode(gin.ReleaseMode)
		logger.Info("🚀 Running in PRODUCTION mode")
	case "development":
		gin.SetMode(gin.TestMode)
		logger.Info("🧪 Running in TEST mode")
	default:
		gin.SetMode(gin.DebugMode)
		logger.Info("🛠️ Running in DEVELOPMENT mode")
	}
	port := config.PORT

	if port == "" {
		port = "8000"
	}

	// Create router from gin
	router := gin.New()

	// Trust only localhost
	err := router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		logger.Fatal("Error setting trusted proxies: " + err.Error())
	}

	// add middleware for logging
	router.Use(gin.Logger())

	// Add auth router
	authGroup := router.Group("/api/v1/auth")
	routes.AuthRouter(authGroup)

	// add user router
	userGroup := router.Group("/api/v1/user")
	routes.UserRouter(userGroup)

	// Health testing route
	router.GET("/api/v1/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"Success": "Access granted for /api/v1"})
	})

	router.Run(":" + port)
}
