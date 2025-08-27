// config/database.go
package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var DB *bun.DB

// InitDB initializes Bun with PostgreSQL / CockroachDB
func InitDB() {
	// Load app config from viper or config.yaml
	LoadConfig()

	dbConf := AppConfig.Database
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbConf.User,
		dbConf.Password,
		dbConf.Host,
		dbConf.Port,
		dbConf.Name,
	)

	// Open underlying *sql.DB
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to open DB: %v", err)
	}

	// Test connectivity
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping DB: %v", err)
	}

	// Assign Bun DB globally
	DB = bun.NewDB(sqlDB, pgdialect.New())

	log.Println("✅ Connected to PostgreSQL (Bun)")
}
