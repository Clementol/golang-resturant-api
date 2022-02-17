package routes

import (
	controller "github.com/Clementol/restur-manag/controllers/v1"
	"github.com/gin-gonic/gin"
)

func FoodRoutes(incommingRoutes *gin.RouterGroup) {
	
	incommingRoutes.GET("/foods", controller.GetFoods())
	incommingRoutes.GET("/foods/:food_id", controller.GetFood())
	incommingRoutes.POST("/foods", controller.CreateFood())
	incommingRoutes.PATCH("/foods/:food_id", controller.UpdateFood())
}
