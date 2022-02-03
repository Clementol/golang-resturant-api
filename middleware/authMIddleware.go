package middleware

import (
	"net/http"
	"strings"

	helper "github.com/Clementol/restur-manag/helpers"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// helpers

		header := c.Request.Header.Get("Authorization")

		clientToken := strings.Split(header, " ")[1]
		if clientToken == "" {
			msg := "Cannot proceed to Authentication"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			c.Abort()
			return
		}
		claims, err := helper.ValidateToken(clientToken)

		if err != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_Name)
		c.Set("last_name", claims.Last_Name)
		c.Set("user_id", claims.Uid)
		c.Next()
	}
}
