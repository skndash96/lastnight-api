package config

import (
	"net/http"
	"os"
	"strconv"
	"time"
)

type AppConfig struct {
	IsProd bool
	Port   int
	DbURL  string
	Auth   AuthConfig
	Minio  MinioConfig
}

type AuthConfig struct {
	Session SessionConfig
	Cookie  CookieConfig
}

type SessionConfig struct {
	Expiry time.Duration
}

type CookieConfig struct {
	Name     string
	Secure   bool
	SameSite http.SameSite
}

type MinioConfig struct {
	Endpoint   string
	Username   string
	Password   string
	UseSSL     bool
	BucketName string
	Expiration time.Duration
	MaxSize    int64
}

func New() *AppConfig {
	port, err := strconv.Atoi(GetEnv("PORT", "1323"))
	if err != nil {
		port = 1323
	}

	isProd, err := strconv.ParseBool(GetEnv("IS_PROD", "true"))
	if err != nil {
		isProd = false
	}

	appCfg := &AppConfig{
		IsProd: isProd,
		Port:   port,
		DbURL:  GetEnv("GOOSE_DBSTRING", ""),

		Auth: AuthConfig{
			Session: SessionConfig{
				Expiry: time.Duration(14*24) * time.Hour,
			},
			Cookie: CookieConfig{
				Name:     "lastnight_token",
				Secure:   isProd,
				SameSite: http.SameSiteLaxMode,
			},
		},

		Minio: MinioConfig{
			Endpoint:   GetEnv("MINIO_ENDPOINT", ""),
			Username:  GetEnv("MINIO_USERNAME", ""),
			Password:  GetEnv("MINIO_PASSWORD", ""),
			BucketName: GetEnv("MINIO_BUCKET_NAME", ""),
			Expiration: time.Duration(15 * time.Minute),
			UseSSL:     isProd,
			MaxSize:    200 * 1024 * 1024, // 200MB
		},
	}

	return appCfg
}

func GetEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
