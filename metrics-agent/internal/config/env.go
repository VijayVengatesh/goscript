package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Env         string
	APIEndpoint string
	InfluxURL   string
	InfluxToken string
	Org         string
	Bucket      string
}

func LoadConfig() *AppConfig {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// Load corresponding .env file
	envFile := fmt.Sprintf(".env.%s", env)
	if err := godotenv.Load(envFile); err != nil {
		fmt.Printf("⚠ Could not load %s file: %v\n", envFile, err)
	}

	return &AppConfig{
		Env:         env,
		APIEndpoint: os.Getenv("API_ENDPOINT"),
		InfluxURL:   os.Getenv("INFLUX_URL"),
		InfluxToken: os.Getenv("INFLUX_TOKEN"),
		Org:         os.Getenv("INFLUX_ORG"),
		Bucket:      os.Getenv("INFLUX_BUCKET"),
	}
}
