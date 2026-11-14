package tasks

import "github.com/hcivekhsim/as-govd/internal/config"

type RuntimeSettings struct {
	UploadDir        string
	MaxFilesCount    int
	MaxFileSizeBytes int64
	MaxBodySizeBytes int64
	AllowedOrigins   []string
	CaptchaSecret    string
}

var runtimeSettings = RuntimeSettings{
	UploadDir:        "uploads/appeals",
	MaxFilesCount:    5,
	MaxFileSizeBytes: 5 * 1024 * 1024,
	MaxBodySizeBytes: 30 * 1024 * 1024,
	AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"},
}

func ApplyRuntimeConfig(cfg *config.Config) {
	if cfg == nil {
		return
	}

	runtimeSettings = RuntimeSettings{
		UploadDir:        cfg.UploadDir,
		MaxFilesCount:    cfg.MaxFilesCount,
		MaxFileSizeBytes: cfg.MaxFileSizeBytes,
		MaxBodySizeBytes: cfg.MaxBodySizeBytes,
		AllowedOrigins:   cfg.CORSAllowedOrigins,
		CaptchaSecret:    cfg.HCaptchaSecret,
	}
}
