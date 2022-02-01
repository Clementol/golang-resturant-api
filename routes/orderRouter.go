package routes

import (
	controller "github.com/Clementol/restur-manag/controllers"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(incommingRoutes *gin.Engine) {

	incommingRoutes.GET("/orders", controller.GetOrders())
	incommingRoutes.POST("/orders", controller.CreateOrder())
	incommingRoutes.GET("/orders/:order_id", controller.GetOrder())
	incommingRoutes.PATCH("/orders/:order_id", controller.UpdateOrder())
}
