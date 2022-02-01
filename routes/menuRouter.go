package routes

import (
	controller "github.com/Clementol/restur-manag/controllers"
	"github.com/gin-gonic/gin"
)

func MenuRoutes(incommingRoutes *gin.Engine) {

	incommingRoutes.GET("/menus", controller.GetMenus())
	incommingRoutes.GET("/menu/:menu_id", controller.GetMenu())
	incommingRoutes.POST("/menu/create", controller.CreateMenu())
	incommingRoutes.PATCH("/menu/update/:menu_id", controller.UpdateMenu())
	incommingRoutes.GET("/menus/foods/:vendor_id", controller.GetMenusWithFoods())
}
