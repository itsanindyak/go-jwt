package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/itsanindyak/go-jwt/controllers"
	"github.com/itsanindyak/go-jwt/middleware"
)

func UserRouter(route *gin.RouterGroup){
	route.Use(middleware.Authenticate())
	route.GET("/",controllers.GetUsers())
	route.GET("/:user_id",controllers.GetUser())
	route.PATCH("/:user_id",controllers.UpdateUser())
	route.DELETE("/:user_id",controllers.DeleteUser())
}