package routes

import (
	controller "github.com/Clementol/restur-manag/controllers"
	"github.com/Clementol/restur-manag/middleware"
	"github.com/gin-gonic/gin"
)

func CartRoutes(incommingRoutes *gin.Engine) {

	incommingRoutes.POST("/cart", middleware.Authentication(), controller.AddItemToCart())
	incommingRoutes.GET("/carts", middleware.Authentication(), controller.GetCartItems())
	incommingRoutes.PATCH("/cart", middleware.Authentication(), controller.RemoveCartItem())

}
