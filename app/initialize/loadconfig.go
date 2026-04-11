package initialize

import (
	"fmt"

	"go-production/global"

	"strings"

	"github.com/spf13/viper"
)

func LoadConfig() {
	v := viper.New()
	v.AddConfigPath(".")      // Tìm config.yaml ở thư mục gốc (nơi chạy binary)
	v.SetConfigName("config") // Tên file: config.yaml
	v.SetConfigType("yaml")

	// Cho phép đọc biến môi trường (Ví dụ: APP_PORT=5000)
	v.AutomaticEnv()
	// Map dấu chấm (.) trong struct sang gạch dưới (_) trong Env (db.dsn -> DB_DSN)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("không đọc được file config: %w", err))
	}

	err = v.Unmarshal(&global.Cfg)
	if err != nil {
		panic(fmt.Errorf("không parse được config: %w", err))
	}

}
