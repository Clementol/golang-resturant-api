package main

import (
	"os"

	"github.com/Clementol/restur-manag/routes"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func main() {
	godotenv.Load()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	// router.Static("/public", "./images")
	routes.UserRoutes(router)
	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.CartRoutes(router)
	routes.OrderRoutes(router)
	routes.InvoiceRoutes(router)
	routes.VendorRoutes(router)
	router.SetTrustedProxies(nil)

	router.Run(":" + PORT)
}
