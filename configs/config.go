package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	User           string
	Password       string
	Name           string
	Host           string
	Port           string
	ConnectTimeout string
	Retries        int
}

type grpcServerConfig struct {
	Port              string
	ConnectionTimeout time.Duration
	MaxMessageSizeMB  int
	ReconnectDelay    time.Duration
}

type AnalyzerConfig struct {
	LogFrequency time.Duration
}

type Config struct {
	DB      DBConfig
	GRPC    grpcServerConfig
	Anomaly AnalyzerConfig
}

func Load() (*Config, error) {
	// Загружаем .env файл (если есть)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Создаем конфиг
	cfg := &Config{
		DB: DBConfig{
			User:           getEnv("POSTGRES_USER", ""),
			Password:       getEnv("POSTGRES_PASSWORD", ""),
			Name:           getEnv("POSTGRES_DB", ""),
			Host:           getEnv("POSTGRES_HOST", ""),
			Port:           getEnv("POSTGRES_PORT", ""),
			ConnectTimeout: getEnv("POSTGRES_CONNECT_TIMEOUT", "10"),
			Retries:        getEnvAsInt("POSTGRES_RETRIES", 1),
		},
		GRPC: grpcServerConfig{
			Port:              getEnv("GRPC_PORT", "50051"),
			ConnectionTimeout: getEnvAsDuration("GRPC_CONNECTION_TIMEOUT", 5*time.Second),
			MaxMessageSizeMB:  getEnvAsInt("GRPC_MAX_MESSAGE_SIZE", 8),
			ReconnectDelay:    getEnvAsDuration("GRPC_RECONNECT_DELAY", 5*time.Second),
		},
		Anomaly: AnalyzerConfig{
			LogFrequency: getEnvAsDuration("ANALYZER_LOG_FREQ", 5*time.Second),
		},
	}

	// Валидация
	if cfg.DB.User == "" || cfg.DB.Password == "" || cfg.DB.Name == "" ||
		cfg.DB.Host == "" || cfg.DB.Port == "" || cfg.DB.Retries <= 0 {
		return nil, fmt.Errorf("missing required database configuration")
	}

	return cfg, nil
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	strValue := getEnv(key, "")
	if strValue == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(strValue)
	if err != nil {
		log.Printf("Invalid value for %s, using default: %v", key, defaultValue)
		return defaultValue
	}
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	strValue := getEnv(key, "")
	if strValue == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(strValue)
	if err != nil {
		log.Printf("Invalid value for %s, using default: %v", key, defaultValue)
		return defaultValue
	}
	return value
}
