package controllers

import (
	"gin-app/config"
	"gin-app/internal/models"
	"gin-app/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func AdminLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login Page",
	})
}

func AdminLoginAction(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	remember := c.PostForm("remember") // "on" if checked

	if email == "" || password == "" {
		c.HTML(http.StatusOK, "login.html", gin.H{"error": "Email and password required"})
		return
	}

	var admin models.User
	err := config.DB.NewSelect().Model(&admin).Where("email = ?", email).Scan(c.Request.Context())
	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{"error": "Invalid email or password"})
		return
	}

	// Set TTL based on "remember me"
	var accessTTL, refreshTTL time.Duration
	if remember == "on" {
		accessTTL, _ = time.ParseDuration("24h")  // 1 day access token
		refreshTTL, _ = time.ParseDuration("30d") // 30 days refresh token
	} else {
		accessTTL, _ = time.ParseDuration("15m") // default 15 min
		refreshTTL, _ = time.ParseDuration("7d") // default 7 days
	}

	// Generate tokens
	accessToken, _ := utils.GenerateToken(admin.ID, config.AppConfig.App.JwtSecret, accessTTL)
	refreshToken, _ := utils.GenerateToken(admin.ID, config.AppConfig.App.JwtSecret, refreshTTL)

	// Store token as Redis key (simplified)
	config.RedisClient.Set(config.Ctx, "admin_access:"+accessToken, admin.ID, accessTTL)
	config.RedisClient.Set(config.Ctx, "admin_refresh:"+refreshToken, admin.ID, refreshTTL)

	// Set cookies
	c.SetCookie("admin_access", accessToken, int(accessTTL.Seconds()), "/", "", false, true)
	c.SetCookie("admin_refresh", refreshToken, int(refreshTTL.Seconds()), "/", "", false, true)

	c.Redirect(http.StatusSeeOther, "/admin/dashboard")
}

func AdminForgetPassword(c *gin.Context) {
	c.HTML(http.StatusOK, "forget-password.html", gin.H{
		"title": "Forget Password",
	})
}

func AdminRefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("admin_refresh")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing refresh token"})
		return
	}

	// Check Redis
	val, err := config.RedisClient.Get(config.Ctx, "admin_refresh:"+refreshToken).Result()
	if err != nil || val == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed successfully"})
}
