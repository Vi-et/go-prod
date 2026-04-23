package initialize

import (
	"fmt"
	"time"

	"go-production/global"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

func InitPostgres() {
	cfg := global.Cfg.DB

	// 1. Kết nối vào Master (Primary)
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: resolveGormLogger(),
	})
	if err != nil {
		panic(fmt.Errorf("không kết nối được Master DB: %w", err))
	}

	// 2. Cấu hình DBResolver cho nhiều Replicas (Read-only)
	var replicas []gorm.Dialector
	for _, dsn := range cfg.Replicas {
		replicas = append(replicas, postgres.Open(dsn))
	}

	duration, _ := time.ParseDuration(cfg.MaxIdleTime)
	err = db.Use(dbresolver.Register(dbresolver.Config{
		Replicas: replicas,
		Policy:   dbresolver.RandomPolicy{}, // Tự động load balance ngẫu nhiên giữa các Slave
	}).
		SetMaxOpenConns(cfg.MaxOpenConns).
		SetMaxIdleConns(cfg.MaxIdleConns).
		SetConnMaxIdleTime(duration))

	if err != nil {
		panic(fmt.Errorf("không cấu hình được DBResolver: %w", err))
	}

	global.DB = db
	global.Logger.Info("PostgreSQL Read/Write Splitting enabled")
}

// resolveGormLogger chọn log level phù hợp theo môi trường
func resolveGormLogger() logger.Interface {
	if global.Cfg.Env == "production" {
		return logger.Default.LogMode(logger.Error) // chỉ log lỗi
	}
	return logger.Default.LogMode(logger.Info) // log toàn bộ query khi dev
}
