package initialize

import (
	"log/slog"
	"os"

	"go-production/global"
)

func InitLogger() {
	var handler slog.Handler

	// Nếu là production thì dùng JSON, development thì dùng Text để dễ đọc
	if global.Cfg.Env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	global.Logger = slog.New(handler)

	// Đặt làm logger mặc định của cả app
	slog.SetDefault(global.Logger)
}
