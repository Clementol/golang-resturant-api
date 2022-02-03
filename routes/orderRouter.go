package routes

import (
	controller "github.com/Clementol/restur-manag/controllers"
	"github.com/Clementol/restur-manag/middleware"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(incommingRoutes *gin.Engine) {

	incommingRoutes.GET("/orders",  controller.GetOrders())
	incommingRoutes.POST("/orders", middleware.Authentication(), controller.CreateOrder())
	incommingRoutes.GET("/order/:order_id", controller.GetOrder())
	incommingRoutes.PATCH("/order/:order_id", controller.UpdateOrder())
}
