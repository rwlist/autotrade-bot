package formula

type Formula interface {
	Calc(now float64) float64
	Start() float64
	Rate() float64
}
