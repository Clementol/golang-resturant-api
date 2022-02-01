package routes

import (
	controller "github.com/Clementol/restur-manag/controllers"
	"github.com/gin-gonic/gin"
)


func VendorRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.POST("/vendor/create", controller.CreateVendor())
	incommingRoutes.GET("/vendors", controller.GetVendors())
	incommingRoutes.GET("/vendor/:vendor_id", controller.GetVendor())
	incommingRoutes.PATCH("/vendor/update/:vendor_id", controller.UpdateVendor())
}