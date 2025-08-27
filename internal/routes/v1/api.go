package v1

import (
	"github.com/gin-gonic/gin"
)

func RegisterApiRoutes(rg *gin.RouterGroup) {
	// Api all routes
	rg.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to the API")
	})
}
