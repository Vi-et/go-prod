package initialize

import (
	"fmt"
	"time"

	"go-production/global"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitPostgres() {
	cfg := global.Cfg.DB

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		// Hiện log query khi đang dev, tắt log khi production
		Logger: resolveGormLogger(),
	})
	if err != nil {
		panic(fmt.Errorf("không kết nối được PostgreSQL: %w", err))
	}

	// Lấy underlying *sql.DB để cấu hình connection pool
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("không lấy được sql.DB từ GORM: %w", err))
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.MaxIdleTime)
	if err != nil {
		panic(fmt.Errorf("max_idle_time không hợp lệ: %w", err))
	}
	sqlDB.SetConnMaxIdleTime(duration)

	// Kiểm tra kết nối thực sự
	if err = sqlDB.Ping(); err != nil {
		panic(fmt.Errorf("không ping được PostgreSQL: %w", err))
	}

	global.DB = db
	global.Logger.Info("PostgreSQL connected via GORM")
}

// resolveGormLogger chọn log level phù hợp theo môi trường
func resolveGormLogger() logger.Interface {
	if global.Cfg.Env == "production" {
		return logger.Default.LogMode(logger.Error) // chỉ log lỗi
	}
	return logger.Default.LogMode(logger.Info) // log toàn bộ query khi dev
}
