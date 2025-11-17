package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppName          string
	AppEnv           string
	AppPort          string
	AllowedOrigins   []string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBSSLMode        string
	DBMaxOpenConns   int
	DBMaxIdleConns   int
	DBConnMaxLifeSec int
	DBConnMaxIdleSec int
	JWTSecret        string
	JWTExpiresMinute int
	AdminEmail       string
	AdminPassword    string
	LogLevel         string
	Redis            RedisConfig
}

type RedisConfig struct {
	Addr       string
	Password   string
	DB         int
	TTLSeconds int
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func Load() *Config {
	expStr := getEnv("JWT_EXPIRES_MINUTES", "1440")
	expInt, err := strconv.Atoi(expStr)
	if err != nil {
		log.Printf("invalid JWT_EXPIRES_MINUTES, fallback 1440")
		expInt = 1440
	}

	rdDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	rdTTL, _ := strconv.Atoi(getEnv("REDIS_TTL_SECONDS", "300"))
	dbMaxOpen, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	dbMaxIdle, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "10"))
	dbLife, _ := strconv.Atoi(getEnv("DB_CONN_MAX_LIFETIME_SEC", "600"))
	dbIdle, _ := strconv.Atoi(getEnv("DB_CONN_MAX_IDLE_TIME_SEC", "300"))

	return &Config{
		AppName:          getEnv("APP_NAME", "qrmenu"),
		AppEnv:           getEnv("APP_ENV", "development"),
		AppPort:          getEnv("APP_PORT", "8080"),
		AllowedOrigins:   strings.Split(getEnv("APP_ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "qrmenu"),
		DBPassword:       getEnv("DB_PASSWORD", "qrmenu"),
		DBName:           getEnv("DB_NAME", "qrmenu_dev"),
		DBSSLMode:        getEnv("DB_SSLMODE", "disable"),
		DBMaxOpenConns:   dbMaxOpen,
		DBMaxIdleConns:   dbMaxIdle,
		DBConnMaxLifeSec: dbLife,
		DBConnMaxIdleSec: dbIdle,
		JWTSecret:        getEnv("JWT_SECRET", "dev_secret"),
		JWTExpiresMinute: expInt,
		AdminEmail:       getEnv("ADMIN_EMAIL", "admin@qrmenu.local"),
		AdminPassword:    getEnv("ADMIN_PASSWORD", "admin123"),
		LogLevel:         getEnv("LOG_LEVEL", "debug"), // dev=debug, prod=info
		Redis: RedisConfig{
			Addr:       getEnv("REDIS_ADDR", "127.0.0.1:6379"),
			Password:   getEnv("REDIS_PASSWORD", ""),
			DB:         rdDB,
			TTLSeconds: rdTTL,
		},
	}
}

func (c *Config) IsProd() bool { return c.AppEnv == "production" }
