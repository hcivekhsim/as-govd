// Package config
package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DBPassword         string
	DBPort             string
	UploadDir          string
	MaxFilesCount      int
	MaxFileSizeBytes   int64
	MaxBodySizeBytes   int64
	CORSAllowedOrigins []string
	HCaptchaSecret     string
}

func Load() *Config {
	cfg := &Config{
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBPort:             os.Getenv("DB_PORT"),
		UploadDir:          defaultString(os.Getenv("UPLOAD_DIR"), "uploads/appeals"),
		MaxFilesCount:      defaultInt(os.Getenv("MAX_FILES_COUNT"), 5),
		MaxFileSizeBytes:   defaultInt64(os.Getenv("MAX_FILE_SIZE_BYTES"), 5*1024*1024),
		MaxBodySizeBytes:   defaultInt64(os.Getenv("MAX_BODY_SIZE_BYTES"), 30*1024*1024),
		CORSAllowedOrigins: defaultOrigins(os.Getenv("CORS_ALLOWED_ORIGINS")),
		HCaptchaSecret:     os.Getenv("HCAPTCHA_SECRET"),
	}

	if cfg.DBPassword == "" {
		return nil
	}

	if cfg.DBPort == "" {
		return nil
	}

	return cfg
}

func defaultString(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func defaultInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func defaultInt64(value string, fallback int64) int64 {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func defaultOrigins(value string) []string {
	if value == "" {
		return []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	if len(result) == 0 {
		return []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	}
	return result
}
