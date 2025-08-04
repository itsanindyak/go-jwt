package helpers

import (
	"github.com/itsanindyak/go-jwt/config"
)

var RolePermissions = map[string]map[string]struct{}{
	config.RoleAdmin: {
		config.PermissionReadUser:   {},
		config.PermissionDeleteUser: {},
		config.PermissionUpdateUser: {},
	},
	config.RoleModerator: {
		config.PermissionReadUser:   {},
		config.PermissionUpdateUser: {},
	},
	config.RoleUser: {
		config.PermissionReadSelf: {},
	},
}

func HasPermissionByRole(role string, permission string) bool {
	if perms, ok := RolePermissions[role]; ok {
		if _, exists := perms[permission]; exists {
			return true
		}
	}
	return false
}
