// Package config handles application configuration loading and management.
package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name      string
		JwtSecret string `mapstructure:"jwt_secret"`
		TokenTTL  string `mapstructure:"token_ttl"`
	} `mapstructure:"app"`

	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		SSLMode  string
	} `mapstructure:"database"`

	Redis struct {
		Host     string
		Port     int
		Password string
		DB       int
	} `mapstructure:"redis"`
}

var AppConfig Config

func LoadConfig() {
	viper.SetConfigName("config") // config.yaml
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config") // config folder

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("❌ Failed to read config: %v", err)
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatalf("❌ Failed to unmarshal config: %v", err)
	}

	// Checking if JWT_SECRET is set
	// jwtKey := []byte(viper.GetString("app.jwt_secret"))
	// log.Println("✅ JWT Key:", string(jwtKey))

	log.Println("✅ Config loaded successfully")
}
