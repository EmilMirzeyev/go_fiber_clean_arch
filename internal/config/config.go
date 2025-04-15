package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseDSN string
	ServerPort  int
	ServerHost  string
}

func NewConfig() *Config {
	config := &Config{
		DatabaseDSN: getEnv("DATABASE_DSN", "user_crud.db"),
		ServerPort:  getEnvAsInt("SERVER_PORT", 8080),
		ServerHost:  getEnv("SERVER_HOST", "0.0.0.0"),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
