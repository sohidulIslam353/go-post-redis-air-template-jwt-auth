package main

import (
	"database/sql"
	"fmt"
	"gin-app/config"
	"gin-app/internal/utils"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/pressly/goose"
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

const migrationTemplate = `
	-- +goose Up
	-- +goose StatementBegin
	CREATE TABLE {{.TableName}} (
		id BIGSERIAL PRIMARY KEY,
		category_id BIGINT NOT NULL,
		name VARCHAR(255) NOT NULL,
		slug VARCHAR(255) NOT NULL,
		status SMALLINT NOT NULL DEFAULT 1,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		-- CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
	);
	-- +goose StatementEnd

	-- +goose StatementBegin
	-- Index on name for faster search
	-- CREATE INDEX idx_{{.TableName}}_name ON {{.TableName}}(name);
	-- +goose StatementEnd

	-- +goose StatementBegin
	-- Unique index on slug to avoid duplicates
	-- CREATE UNIQUE INDEX idx_{{.TableName}}_slug ON {{.TableName}}(slug);
	-- +goose StatementEnd

	-- +goose StatementBegin
	-- Combined index on category_id + name (optional, useful for filtering by category)
	-- CREATE INDEX idx_{{.TableName}}_category_name ON {{.TableName}}(category_id, name);
	-- +goose StatementEnd

	-- +goose Down

	-- +goose StatementBegin
	DROP TABLE IF EXISTS {{.TableName}};
	-- +goose StatementEnd
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run cmd/commands/make.go make:controller ControllerName")
		fmt.Println("  go run cmd/commands/make.go make:model ModelName")
		fmt.Println("  go run cmd/commands/make.go make:middleware MiddlewareName")
		fmt.Println("  go run cmd/commands/make.go make:migration MigrationName")
		fmt.Println("  go run cmd/commands/make.go migrate:up")
		fmt.Println("  go run cmd/commands/make.go migrate:down")
		fmt.Println("  go run cmd/commands/make.go migrate:status")
		return
	}

	command := os.Args[1]

	// Only load name if the command requires it
	var name string
	if len(os.Args) >= 3 {
		name = os.Args[2]
	}

	// for migration task
	config.LoadConfig()
	dsn := config.GetDSN() // ✅ use same DSN everywhere
	dir := "./migrations"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// end

	switch command {
	case "make:controller":
		if name == "" {
			log.Fatal("❌ Please provide ControllerName")
		}
		createController(name)
	case "make:model":
		if name == "" {
			log.Fatal("❌ Please provide ModelName")
		}
		createModel(name)
	case "make:middleware":
		if name == "" {
			log.Fatal("❌ Please provide MiddlewareName")
		}
		createMiddleware(name)
	case "make:migration":
		if name == "" {
			log.Fatal("❌ Please provide MigrationName")
		}
		createMigration(name)
	case "migrate:up":
		_ = goose.Up(db, dir)
		fmt.Println("✅ Migration completed!")
	case "migrate:down":
		_ = goose.Down(db, dir)
		fmt.Println("✅ Migration rolled back!")
	case "migrate:status":
		_ = goose.Status(db, dir)
	default:
		fmt.Println("❌ Unknown command:", command)
	}
}

// Controller Create
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

// Model Create
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

// Middleware Creat
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

// Migration file generated
func createMigration(name string) {
	migrationName := strings.ToLower(name)
	fileName := time.Now().Format("20060102150405") + "_create_" + strings.ToLower(name) + "_table.sql"
	filePath := "migrations/" + fileName

	_ = os.MkdirAll("migrations", os.ModePerm)

	if _, err := os.Stat(filePath); err == nil {
		fmt.Println("⚠️ Migration already exists:", filePath)
		return
	}

	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println("❌ Error creating migration:", err)
		return
	}
	defer f.Close()

	tmpl, _ := template.New("migration").Parse(migrationTemplate)
	tmpl.Execute(f, map[string]string{
		"MigrationName": migrationName,
		"TableName":     strings.ToLower(name),
	})

	fmt.Println("✅ Migration created at:", filePath)
}

// Migration run command
