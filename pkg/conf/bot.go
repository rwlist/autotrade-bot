package conf

type Bot struct {
	AdminID     int    `env:"ADMIN_TELEGRAM_ID,required"`
	Token       string `env:"BOT_TOKEN,required"`
	PrettyPrint bool   `env:"PRETTY_LOGS" envDefault:"false"`
}
