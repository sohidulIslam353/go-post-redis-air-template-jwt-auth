package middleware

import (
	"gin-app/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("admin_access")
		if err != nil || token == "" {
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.Abort()
			return
		}

		// Check Redis
		_, err = config.RedisClient.Get(config.Ctx, "admin_access:"+token).Result()
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/admin/login")
			c.Abort()
			return
		}

		c.Next()
	}
}
