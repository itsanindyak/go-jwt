package config


const (
	PermissionReadUser       = "read_user"
	PermissionDeleteUser     = "delete_user"
	PermissionUpdateUser     = "update_user"
	PermissionUpdateUserType = "update_user_type"
	PermissionReadSelf       = "read_self"

	RoleAdmin     = "ADMIN"
	RoleModerator = "MODERATOR"
	RoleUser      = "USER"
)

var RolePermissions = map[string]map[string]struct{}{
	RoleAdmin: {
		PermissionReadUser:       {},
		PermissionDeleteUser:     {},
		PermissionUpdateUser:     {},
		PermissionUpdateUserType: {},
	},
	RoleModerator: {
		PermissionReadUser:   {},
		PermissionUpdateUser: {},
	},
	RoleUser: {
		PermissionReadSelf: {},
	},
}
