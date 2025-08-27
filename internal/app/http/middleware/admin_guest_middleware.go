package middleware

import (
	"gin-app/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminGuestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie("admin_access")
		if token != "" {
			_, err := config.RedisClient.Get(config.Ctx, "admin_access:"+token).Result()
			if err == nil {
				c.Redirect(http.StatusSeeOther, "/admin/dashboard")
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
