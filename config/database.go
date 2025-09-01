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

// GetDSN builds a PostgreSQL DSN string from config
func GetDSN() string {
	if AppConfig.Database.Host == "" {
		LoadConfig()
	}

	dbConf := AppConfig.Database
	sslmode := dbConf.SSLMode
	if sslmode == "" {
		sslmode = "disable" // default fallback
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbConf.User,
		dbConf.Password,
		dbConf.Host,
		dbConf.Port,
		dbConf.Name,
		sslmode,
	)
}

// InitDB initializes Bun with PostgreSQL
func InitDB() {
	dsn := GetDSN()

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
