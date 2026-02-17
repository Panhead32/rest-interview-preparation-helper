package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	Server ServerConfig
	DB     DatabaseConfig
	Log    LogConfig
	OpenAI OpenAIConfig
	Auth   AuthConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port    string
	Timeout int // in seconds
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	DSN string
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level string
}

type OpenAIConfig struct {
	APIKey string
}

type AuthConfig struct {
	JWTSecret string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	var once sync.Once
	once.Do(func() {
		fmt.Println("Loading configuration...")
		if err := godotenv.Load(); err != nil {
			fmt.Printf("Error loading .env file: %v\n", err)
		}
	})

	return &Config{
		Server: ServerConfig{
			Port:    getStringEnv("SERVER_PORT", "8080"),
			Timeout: getIntEnv("SERVER_TIMEOUT", 30),
		},
		DB: DatabaseConfig{
			DSN: getStringEnv("DB_DSN", ""),
		},
		Log: LogConfig{
			Level: getStringEnv("LOG_LEVEL", "INFO"),
		},
		OpenAI: OpenAIConfig{
			APIKey: getStringEnv("OPENAI_API_KEY", "")},
		Auth: AuthConfig{
			JWTSecret: getStringEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		},
	}
}

func getEnv(key string, defaultValue any) any {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getIntEnv(key string, defaultValue int) int {
	valueStr := getEnv(key, strconv.Itoa(defaultValue)).(string)

	value, err := strconv.Atoi(valueStr)

	if err != nil {
		fmt.Printf("Error parsing %s: %v. Using default value %d\n", key, err, defaultValue)
		return defaultValue
	}

	return value
}

func getStringEnv(key string, defaultValue string) string {
	value := getEnv(key, defaultValue)
	return value.(string)
}
