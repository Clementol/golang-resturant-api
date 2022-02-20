package routes

import (
	controller "github.com/Clementol/restur-manag/controllers/v1/order"
	"github.com/Clementol/restur-manag/middleware"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(incommingRoutes *gin.RouterGroup) {

	incommingRoutes.GET("/orders", middleware.Authentication(),  controller.GetUserOrders())
	incommingRoutes.POST("/orders", middleware.Authentication(), controller.CreateOrder())

}
