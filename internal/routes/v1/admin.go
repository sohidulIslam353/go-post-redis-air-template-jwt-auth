package v1

import (
	admin_controller "gin-app/internal/app/http/controllers/admin"
	"gin-app/internal/app/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(rg *gin.RouterGroup) {

	// Admin routes go here
	auth := rg.Group("/").Use(middleware.AdminGuestMiddleware())
	{
		auth.GET("/login", admin_controller.AdminLogin)
		auth.POST("/login", admin_controller.AdminLoginAction)
		auth.GET("/forget-password", admin_controller.AdminForgetPassword)
	}

	// After login
	admin := rg.Group("/").Use(middleware.AdminAuthMiddleware())
	{
		admin.GET("/dashboard", admin_controller.AdminDashboard)
		admin.GET("/logout", admin_controller.AdminLogout)

	}
	// Refresh token
	rg.POST("/refresh", admin_controller.AdminRefreshToken)

}
