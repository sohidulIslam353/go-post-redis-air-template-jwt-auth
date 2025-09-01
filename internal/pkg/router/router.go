package router

import (
	"gin-app/internal/routes"
	"gin-app/internal/utils"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"asset": utils.Asset, // utils.Asset কে "asset" নামে template এ expose করলাম
		"add":   func(a, b int) int { return a + b },
		"formatDate": func(t time.Time) string {
			return t.Format("02 Jan 2006")
		},
	})

	// load html
	r.Static("/static", "./static") // load static files
	r.LoadHTMLGlob("templates/*")

	// ✅ Root route
	r.GET("/", func(c *gin.Context) {
		c.String(200, "সাবধান বিপদের দিকে এগোবেন না।")
	})
	// 404 page
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"title": "Page Not Found",
		})
	})

	// ✅ Register API routes
	routes.RegisterRoutes(r)
	return r
}
