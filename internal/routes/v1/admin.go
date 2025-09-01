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
		auth.GET("/admin-create", admin_controller.AdminCreate)
	}

	// After login
	admin := rg.Group("/").Use(middleware.AdminAuthMiddleware())
	{
		admin.GET("/dashboard", admin_controller.AdminDashboard)
		admin.GET("/logout", admin_controller.AdminLogout)

		// Job Type Routes
		admin.GET("/job-type-list", admin_controller.AdminJobTypeList)
		admin.GET("/job-type-create", admin_controller.AdminJobTypeCreate)
		admin.POST("/job-type-store", admin_controller.AdminJobTypeStore)
		admin.POST("/job-type-status/:id", admin_controller.AdminToggleJobTypeStatus)
		admin.DELETE("/job-type-delete/:id", admin_controller.AdminDeleteJobType)
		admin.GET("/job-type-edit/:id", admin_controller.AdminEditJobType)
		admin.POST("/job-type-update/:id", admin_controller.AdminUpdateJobType)

		// Category Routes
		admin.GET("/category-list", admin_controller.AdminCategoryList)
		admin.GET("/category-create", admin_controller.AdminCategoryCreate)
		admin.POST("/category-store", admin_controller.AdminCategoryStore)
		admin.GET("/category-edit/:id", admin_controller.AdminEditCategory)
		admin.POST("/category-update/:id", admin_controller.AdminUpdateCategory)
		admin.DELETE("/category-delete/:id", admin_controller.AdminDeleteCategory)
		admin.POST("/category-status/:id", admin_controller.AdminToggleCategoryStatus)

		// Sub category routes
		admin.GET("/subcategory-list", admin_controller.AdminSubCategoryList)
		admin.GET("/subcategory-create", admin_controller.AdminSubCategoryCreate)
		admin.POST("/subcategory-store", admin_controller.AdminSubCategoryStore)
		admin.GET("/subcategory-edit/:id", admin_controller.AdminEditSubCategory)
		admin.POST("/subcategory-update/:id", admin_controller.AdminUpdateSubCategory)
		admin.DELETE("/subcategory-delete/:id", admin_controller.AdminDeleteSubCategory)
		admin.POST("/subcategory-status/:id", admin_controller.AdminToggleSubCategoryStatus)
	}

	// Refresh token
	rg.POST("/refresh", admin_controller.AdminRefreshToken)

}
