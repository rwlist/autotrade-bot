package app

const (
	prefix = "ovpn_"

	data = "data"
	tcp = "tcp"
	udp = "udp"
)

type Logic struct {}

func NewLogic() *Logic {
	return &Logic{}
}

func (l *Logic) CommandStatus() (string, error) {
	rate, err := binanceRateQuery()
	return rate, err
}