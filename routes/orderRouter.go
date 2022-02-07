package routes

import (
	controller "github.com/Clementol/restur-manag/controllers"
	"github.com/Clementol/restur-manag/middleware"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(incommingRoutes *gin.Engine) {

	incommingRoutes.GET("/orders", middleware.Authentication(),  controller.GetUserOrders())
	incommingRoutes.POST("/orders", middleware.Authentication(), controller.CreateOrder())

}
