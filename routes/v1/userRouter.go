package routes

import (
	controller "github.com/Clementol/restur-manag/controllers/v1/user"
	"github.com/Clementol/restur-manag/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.RouterGroup) {
	incomingRoutes.POST("/users/signup", controller.Signup())
	incomingRoutes.POST("/users/login", controller.Login())
	incomingRoutes.GET("/users/:user_id", middleware.Authentication(), controller.GetUser())
	incomingRoutes.PATCH("/users/:user_id", middleware.Authentication(), controller.UpdateUser())

	incomingRoutes.GET("/users", controller.GetUsers())
}
