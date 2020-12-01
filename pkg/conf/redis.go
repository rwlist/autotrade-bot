package conf

type Redis struct {
	Password string `env:"REDIS_PASSWORD"`
	Addr     string `env:"REDIS_ADDR"`
}
