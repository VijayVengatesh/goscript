// package config

// import (
// 	"fmt"
// 	"os"

// 	"github.com/joho/godotenv"
// )

// type AppConfig struct {
// 	Env         string
// 	APIEndpoint string
// 	InfluxURL   string
// 	InfluxToken string
// 	Org         string
// 	Bucket      string
// }

// func LoadConfig() *AppConfig {
// 	env := os.Getenv("APP_ENV")
// 	if env == "" {
// 		env = "production" // Default to production if not set
// 	}

// 	// Load corresponding .env file
// 	envFile := fmt.Sprintf(".env.%s", env)
// 	if err := godotenv.Load(envFile); err != nil {
// 		fmt.Printf("⚠ Could not load %s file: %v\n", envFile, err)
// 	}

//		return &AppConfig{
//			Env:         env,
//			APIEndpoint: os.Getenv("API_ENDPOINT"),
//			InfluxURL:   os.Getenv("INFLUX_URL"),
//			InfluxToken: os.Getenv("INFLUX_TOKEN"),
//			Org:         os.Getenv("INFLUX_ORG"),
//			Bucket:      os.Getenv("INFLUX_BUCKET"),
//		}
//	}
package config

import (
	"fmt"
	"os"
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
		env = "production" // default to production
	}

	var cfg AppConfig
	cfg.Env = env

	switch env {
	case "development":
		cfg.APIEndpoint = "http://localhost:3000/api"
		cfg.InfluxURL = "http://localhost:8086"
		cfg.InfluxToken = "dev-token"
		cfg.Org = "dev-org"
		cfg.Bucket = "dev-bucket"

	case "production":
		cfg.APIEndpoint = "http://localhost:5000"
		cfg.InfluxURL = "http://10.1.1.152:8086/"
		cfg.InfluxToken = "BQqNr9GGBFtAoFCkAxMwb3LdhsRIEhBZmyoI22ua1I38xU5nSIBcYfZI2-xyImoQHYv9fPESBdGKX542uncptA=="
		cfg.Org = "idevopz"
		cfg.Bucket = "metrics"

	case "test":
		cfg.APIEndpoint = "http://localhost:3001/test-api"
		cfg.InfluxURL = "http://localhost:8087"
		cfg.InfluxToken = "test-token"
		cfg.Org = "test-org"
		cfg.Bucket = "test-bucket"

	default:
		fmt.Printf("⚠ Unknown APP_ENV: %s, defaulting to production config\n", env)
		cfg.APIEndpoint = "https://api.yourdomain.com"
		cfg.InfluxURL = "https://influxdb.yourdomain.com"
		cfg.InfluxToken = "prod-token"
		cfg.Org = "prod-org"
		cfg.Bucket = "prod-bucket"
	}

	return &cfg
}
