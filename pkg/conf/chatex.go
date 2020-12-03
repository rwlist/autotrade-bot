package conf

type Chatex struct {
	RefreshToken string `env:"CHATEX_REFRESH_TOKEN"`
	URL          string `env:"CHATEX_API" envDefault:"https://api.chatex.com/v1"`
}
