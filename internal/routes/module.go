package routes

import (
	v1 "gin-app/internal/routes/v1"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	v1.RegisterAdminRoutes(router.Group("/admin"))
	v1.RegisterApiRoutes(router.Group("/api/v1"))
}
