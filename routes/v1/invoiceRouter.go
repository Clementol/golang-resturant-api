package routes

import (
	controller "github.com/Clementol/restur-manag/controllers/v1/invoice"
	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(incommingRoutes *gin.RouterGroup) {

	incommingRoutes.GET("/invoices", controller.GetInvoices())
	// incommingRoutes.GET("/invoices/:invoice_id", controller.GetInvoice())
	incommingRoutes.POST("/invoices", controller.CreateInvoice())
	incommingRoutes.PATCH("/invoices/:invoice_id", controller.UpdateInvoice())

}
