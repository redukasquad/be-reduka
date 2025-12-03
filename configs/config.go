package configs

import (
	"os"
)

type Config struct {
	Port   string
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string

	JWTSecret          string
	JWTExpiry          string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port:   port,
		DBUser: os.Getenv("MYSQL_USER"),
		DBPass: os.Getenv("MYSQL_PASSWORD"),
		DBHost: "mysql", // Default to docker service name, or localhost if running locally
		DBPort: "3306",
		DBName: os.Getenv("MYSQL_DATABASE"),

		JWTSecret:          os.Getenv("JWT_SECRET"),
		JWTExpiry:          os.Getenv("JWT_EXPIRY"),
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	}
}
