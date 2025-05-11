package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/itsanindyak/go-jwt/config"
	"github.com/itsanindyak/go-jwt/controllers"
	"github.com/itsanindyak/go-jwt/middleware"
)

func UserRouter(route *gin.RouterGroup) {
	route.Use(middleware.Authenticate())

	route.GET("/:user_id",
		middleware.GrantMiddleware([]string{config.PermissionReadSelf, config.PermissionReadUser}),
		controllers.GetUser())

	route.GET("/",
		middleware.GrantMiddleware([]string{config.PermissionReadUser}),
		controllers.GetUsers())

	route.PATCH("/:user_id",
		middleware.GrantMiddleware([]string{config.PermissionUpdateUser}),
		controllers.UpdateUser())

	route.DELETE("/:user_id",
		middleware.GrantMiddleware([]string{config.PermissionDeleteUser}),
		controllers.DeleteUser())
}
