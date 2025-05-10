package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type GrantOptions struct {
	RequestedUserID string
	AdminOnly       bool
}

func CheckGrant(c *gin.Context, opts GrantOptions) error {
	userType := c.GetString("user_type")
	loggedInUserId := c.GetString("uid")

	if userType != "USER" && userType != "ADMIN" {
		return errors.New("unauthorized: invalid user type")
	}

	if opts.AdminOnly && userType != "ADMIN" {
		return errors.New("unauthorized: admin access required")
	}

	if userType == "USER" && opts.RequestedUserID != "" && loggedInUserId != opts.RequestedUserID {
		return errors.New("unauthorized: user can only access their own data")
	}

	return nil
}
