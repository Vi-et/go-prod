package initialize

import (
	"fmt"

	"go-production/global"
)

func InitializeApp() {
	// 1. Load config → global.Cfg
	LoadConfig()

	// 2. Init Logger → global.Logger
	InitLogger()
	global.Logger.Info("Logger initialized")

	// 3. Init Database → global.DB
	InitPostgres()

	// 4. Init Router và chạy server
	r := InitRouter()

	addr := fmt.Sprintf(":%d", global.Cfg.Port)
	global.Logger.Info("Server starting", "addr", addr, "env", global.Cfg.Env)
	r.Run(addr)
}
