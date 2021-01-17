package config

import (
	"os"
)

// Config defines the shape of the server's configuration
type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	ServerPort string
}

// Get returns a pointer to the configuration
func Get() *Config {
	dbUser := getEnvWithDefault("POSTGRES_USER", "user")
	dbPassword := getEnvWithDefault("POSTGRES_PASSWORD", "password")
	dbName := getEnvWithDefault("POSTGRES_DB", "chef")
	dbHost := getEnvWithDefault("POSTGRES_HOST", "database")
	dbPort := getEnvWithDefault("POSTGRES_PORT", "5432")
	serverPort := getEnvWithDefault("SERVER_PORT", "8080")

	return &Config{
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		DBHost:     dbHost,
		DBPort:     dbPort,
		ServerPort: serverPort,
	}
}

func getEnvWithDefault(key, defaultVal string) string {
	envVar := os.Getenv(key)
	if envVar == "" {
		return defaultVal
	}
	return envVar
}
