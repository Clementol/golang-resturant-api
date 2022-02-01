package middleware

import (
	"net/http"

	helper "github.com/Clementol/restur-manag/helpers"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// helpers
		
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			msg := "Cannot proceed to Authentication"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			c.Abort()
			return
		}
		claims, err := helper.ValidateToken(clientToken)
		
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_Name)
		c.Set("last_name", claims.Last_Name)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
