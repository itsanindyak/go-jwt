package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/itsanindyak/go-jwt/controllers"
)

func AuthRouter(route *gin.RouterGroup) {

	route.POST("/signup", controllers.Signup())
	route.POST("/login", controllers.Login())

}
