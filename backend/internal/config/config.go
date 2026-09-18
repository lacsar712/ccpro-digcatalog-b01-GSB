package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	Port       string
}

func Load() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3307"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "root"),
		DBName:     getEnv("DB_NAME", "digcatalog"),
		JWTSecret:  getEnv("JWT_SECRET", "digcatalog-jwt-secret-change-me"),
		Port:       getEnv("PORT", "8080"),
	}
}

func (c *Config) DSN() string {
	// loc 显式固定为 Asia/Shanghai：MySQL DATETIME 不含时区，
	// GORM 读写 created_at 等时间字段时一律按东八区解释，
	// 不依赖后端容器/宿主机的 TZ 设置（运行镜像需含 tzdata）。
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
