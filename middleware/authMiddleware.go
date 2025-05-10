package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/itsanindyak/go-jwt/helpers"
)

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string

		authHeader := ctx.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if token == "" {
			cookieToken, err := ctx.Request.Cookie("token")
			if err == nil {
				token = cookieToken.Value
			}
		}

		if token == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No cookie found"})
			ctx.Abort()
			return
		}

		tokenData, msg := helpers.ParseToken(token)

		if msg != "" {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": msg})
			ctx.Abort()
			return
		}

		ctx.Set("uid", tokenData.UID)

		ctx.Set("user_type", tokenData.UserType)
		ctx.Next()

	}

}
