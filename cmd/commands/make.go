package main

import (
	"fmt"
	"gin-app/internal/utils"
	"os"
	"strings"
	"text/template"
)

// This template is for creating a basic controller in a Gin web application
const controllerTemplate = `package controllers

import "github.com/gin-gonic/gin"

func {{.ControllerName}}(c *gin.Context) {
	c.JSON(200, gin.H{"message": "{{.ControllerName}} is working!"})
}
`

// This template is for creating a basic model with Bun ORM
const modelTemplate = `package models

	import (
		"time"

		"github.com/uptrace/bun"
	)

	type {{.ModelName}} struct {
		bun.BaseModel ` + "`bun:\"table:{{.TableName}}s\"`" + `
		ID        int64     ` + "`bun:\"id,pk,autoincrement\"`" + `
		Name      string    ` + "`bun:\"name,notnull\"`" + `
		CreatedAt time.Time ` + "`bun:\"created_at,default:now()\"`" + `
		UpdatedAt time.Time ` + "`bun:\"updated_at,default:now()\"`" + `
	}
`

// Middleware for creating a controller or model in a Go web application
const middlewareTemplate = `package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func {{.MiddlewareName}}() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: implement middleware logic here

		// Example logic (remove or customize this):
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Next()
	}
}
`

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run cmd/commands/make.go controller ControllerName")
		fmt.Println("  go run cmd/commands/make.go model ModelName")
		fmt.Println("  go run cmd/commands/make.go middleware MiddlewareName")
		return
	}

	command := os.Args[1]
	name := os.Args[2]

	switch command {
	case "controller":
		createController(name)
	case "model":
		createModel(name)
	case "middleware":
		createMiddleware(name)
	default:
		fmt.Println("Unknown command:", command)
	}
}

func createController(name string) {
	controllerName := "Index" // will be method name
	fileName := strings.ToLower(name) + ".go"
	filePath := "internal/app/http/controllers/" + fileName

	_ = os.MkdirAll("internal/app/http/controllers", os.ModePerm)

	if _, err := os.Stat(filePath); err == nil {
		fmt.Println("⚠️ Controller already exists:", filePath)
		return
	}

	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println("❌ Error creating controller:", err)
		return
	}
	defer f.Close()

	tmpl, _ := template.New("controller").Parse(controllerTemplate)
	tmpl.Execute(f, map[string]string{
		"ControllerName": controllerName,
	})

	fmt.Println("✅ Controller created at:", filePath)
}

func createModel(name string) {
	modelName := utils.ToTitleCase(name)
	tableName := strings.ToLower(name)
	fileName := strings.ToLower(name) + ".go"
	filePath := "internal/models/" + fileName

	_ = os.MkdirAll("internal/models", os.ModePerm)

	if _, err := os.Stat(filePath); err == nil {
		fmt.Println("⚠️ Model already exists:", filePath)
		return
	}

	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println("❌ Error creating model:", err)
		return
	}
	defer f.Close()

	tmpl, _ := template.New("model").Parse(modelTemplate)
	tmpl.Execute(f, map[string]string{
		"ModelName": modelName,
		"TableName": tableName,
	})

	fmt.Println("✅ Model created at:", filePath)
}

func createMiddleware(name string) {
	middlewareName := utils.ToTitleCase(name)
	fileName := strings.ToLower(name) + ".go"
	filePath := "internal/app/http/middleware/" + fileName

	_ = os.MkdirAll("internal/app/http/middleware", os.ModePerm)

	if _, err := os.Stat(filePath); err == nil {
		fmt.Println("⚠️ Middleware already exists:", filePath)
		return
	}

	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println("❌ Error creating middleware:", err)
		return
	}
	defer f.Close()

	tmpl, _ := template.New("middleware").Parse(middlewareTemplate)
	tmpl.Execute(f, map[string]string{
		"MiddlewareName": middlewareName,
	})

	fmt.Println("✅ Middleware created at:", filePath)
}
