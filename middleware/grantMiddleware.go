package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsanindyak/go-jwt/helpers"
)

func GrantMiddleware(Permissions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("uid")
		userType := c.GetString("user_type")
		userParamID := c.Param("user_id")

		hasPermission := false

		for _, permission := range Permissions {
			log.Printf("User with ID %s and type %s is attempting to access permission %s", userID, userType, permission)

			// Special case: 'read_self' for regular users
			if permission == "read_self" && userType == "USER" {
				if userParamID == userID {
					hasPermission = true
					break
				}
				continue
			}

			// General case: permission based on role
			if helpers.HasPermissionByRole(userType, permission) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permission"})
			return
		}

		c.Next()
	}
}

