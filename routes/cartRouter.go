package routes

import (
	controller "github.com/Clementol/restur-manag/controllers"
	"github.com/Clementol/restur-manag/middleware"
	"github.com/gin-gonic/gin"
)

func CartRoutes(incommingRoutes *gin.Engine) {

	incommingRoutes.GET("/carts", middleware.Authentication(), controller.GetCartItems())
	incommingRoutes.POST("/cart", middleware.Authentication(), controller.AddItemToCart())
	incommingRoutes.PATCH("/cart", middleware.Authentication(), controller.RemoveCartItem())

}
