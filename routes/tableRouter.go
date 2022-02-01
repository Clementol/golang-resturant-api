package routes

import (
	controller "github.com/Clementol/restur-manag/controllers"
	"github.com/gin-gonic/gin"
)

func TableRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/tables", controller.GetTables())
	incommingRoutes.GET("/tables/:table_id", controller.GetTable())
	incommingRoutes.POST("/tables", controller.CreateTable())
	incommingRoutes.PATCH("/tables/:table_id", controller.UpdateTable())
}
