package routes

import (
	controller "github.com/Clementol/restur-manag/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controller.Signup())
	incomingRoutes.POST("/users/login", controller.Login())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
}
