package helpers

import (
	"github.com/itsanindyak/go-jwt/config"
)

func HasPermissionByRole(role string, permission string) bool {
	if perms, ok := config.RolePermissions[role]; ok {
		if _, exists := perms[permission]; exists {
			return true
		}
	}
	return false
}
