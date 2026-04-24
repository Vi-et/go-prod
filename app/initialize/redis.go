package initialize

import (
	"context"
	"fmt"
	"go-production/global"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis() {
	cfg := global.Cfg.Redis

	dialTimeout, err := time.ParseDuration(cfg.DialTimeout)
	if err != nil {
		dialTimeout = 3 * time.Second
	}
	readTimeout, err := time.ParseDuration(cfg.ReadTimeout)
	if err != nil {
		readTimeout = 1 * time.Second
	}
	writeTimeout, err := time.ParseDuration(cfg.WriteTimeout)
	if err != nil {
		writeTimeout = 1 * time.Second
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(fmt.Errorf("không kết nối được Redis: %w", err))
	}
	global.Redis = rdb
	global.Logger.Info("Redis connected", "addr", cfg.Addr)

}
