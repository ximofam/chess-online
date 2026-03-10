package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	MySQL  MySQLConfig
	JWT    JWTConfig
	Redis  RedisConfig
}

var Envs *Config

func init() {
	if err := godotenv.Load(".env.dev"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	Envs = &Config{
		Server: ServerConfig{
			Addr:   fmt.Sprintf(":%s", getEnv("SERVER_PORT", "8080")),
			APIKey: getEnv("SERVER_API_KEY", "MY_API_KEY"),
		},
		MySQL: MySQLConfig{
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "3306"),
			DBUser:     getEnv("DB_USER", "root"),
			DBPassword: getEnv("DB_PASSWORD", "123456"),
			DBName:     getEnv("DB_NAME", "quiz"),
		},
		JWT: JWTConfig{
			SecretKey:       getEnv("JWT_SECRET_KEY", "MY_SECRET_KEY"),
			AccessTokenTTL:  getEnvAsDuration("ACCESS_TOKEN_TTL", 2*time.Hour),
			RefreshTokenTTL: getEnvAsDuration("REFRESH_TOKEN_TTL", 72*time.Hour),
		},
		Redis: RedisConfig{
			Addr:     fmt.Sprintf("%s:%s", getEnv("REDIS_HOST", "localhost"), getEnv("REDIS_PORT", "6379")),
			Password: getEnv("REDIS_PASSWORD", "123456"),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
	}
}

type ServerConfig struct {
	Addr   string
	APIKey string
}

type MySQLConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

type JWTConfig struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}
	return d
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}

	return fallback
}
