package global

import (
	"go-production/globalStructs"
	"log/slog"

	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

var (
	// Cfg chứa toàn bộ cấu hình ứng dụng được load từ config.yaml
	Cfg globalStructs.Config

	// DB là GORM database instance
	DB *gorm.DB

	// Logger là structured logger dùng slog (Go 1.21+)
	Logger *slog.Logger

	Redis *redis.Client
)
