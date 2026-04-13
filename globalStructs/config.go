package globalStructs

type Config struct {
	Port int    `mapstructure:"port"`
	Env  string `mapstructure:"env"`

	DB struct {
		DSN          string `mapstructure:"dsn"`
		MaxOpenConns int    `mapstructure:"max_open_conns"`
		MaxIdleConns int    `mapstructure:"max_idle_conns"`
		MaxIdleTime  string `mapstructure:"max_idle_time"`
	} `mapstructure:"db"`

	Limiter struct {
		RPS     float64 `mapstructure:"rps"`
		Burst   int     `mapstructure:"burst"`
		Enabled bool    `mapstructure:"enabled"`
	} `mapstructure:"limiter"`

	SMTP struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Sender   string `mapstructure:"sender"`
	} `mapstructure:"smtp"`

	CORS struct {
		TrustedOrigins []string `mapstructure:"trusted_origins"`
	} `mapstructure:"cors"`

	JWT struct {
		Secret     string `mapstructure:"secret"`
		Expiration string `mapstructure:"expiration"`
	} `mapstructure:"jwt"`
}
