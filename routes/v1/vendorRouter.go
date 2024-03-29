package routes

import (
	controller "github.com/Clementol/restur-manag/controllers/v1/vendor"
	"github.com/gin-gonic/gin"
)


func VendorRoutes(incommingRoutes *gin.RouterGroup) {
	incommingRoutes.POST("/vendor/create", controller.CreateVendor())
	incommingRoutes.PATCH("/vendor/update/:vendor_id", /*vendorAuth*/ controller.UpdateVendor())
	incommingRoutes.GET("/vendor/:vendor_id", controller.GetVendor())
	incommingRoutes.GET("/vendor/:vendor_id/orders", /*vendorAuth*/ controller.GetVendorOrders())

	incommingRoutes.PATCH("/vendor/:vendor_id/orders/update", /*vendorAuth*/ controller.UpdateVendorOrder())

	// incommingRoutes.GET("/vendors", controller.GetVendors())
}