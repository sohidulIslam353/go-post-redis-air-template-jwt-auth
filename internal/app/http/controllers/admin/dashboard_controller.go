package controllers

import (
	"gin-app/config"
	"gin-app/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminDashboard(c *gin.Context) {
	admin, _, err := utils.GetLoggedInAdmin(c)
	if err != nil {
		c.Redirect(303, "/admin/login")
		return
	}
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"admin": admin,
		"title": "Dashboard",
	})

}

// Admin logout
func AdminLogout(c *gin.Context) {
	accessToken, _ := c.Cookie("admin_access")
	refreshToken, _ := c.Cookie("admin_refresh")

	if accessToken != "" {
		config.RedisClient.Del(config.Ctx, "admin_access:"+accessToken)
	}
	if refreshToken != "" {
		config.RedisClient.Del(config.Ctx, "admin_refresh:"+refreshToken)
	}

	// Clear cookies
	c.SetCookie("admin_access", "", -1, "/", "", false, true)
	c.SetCookie("admin_refresh", "", -1, "/", "", false, true)

	c.Redirect(http.StatusSeeOther, "/admin/login")
}
