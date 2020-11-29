module github.com/rwlist/autotrade-bot

go 1.13

require (
	github.com/adshao/go-binance v0.0.0-20200302052924-65a935d32ae9
	github.com/caarlos0/env/v6 v6.1.0
	github.com/chatex-com/sdk-go v0.0.0
	github.com/davecgh/go-spew v1.1.1
	github.com/go-redis/redis/v8 v8.4.0
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/petuhovskiy/telegram v0.0.0-20200207220211-6250cca70f10
	github.com/pplcc/plotext v0.0.0-20180221170324-68ab3c6e05c3
	github.com/prometheus/client_golang v1.8.0
	github.com/shopspring/decimal v1.2.0
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/image v0.0.0-20200119044424-58c23975cae1 // indirect
	gonum.org/v1/netlib v0.0.0-20200317120129-c5a04cffd98a // indirect
	gonum.org/v1/plot v0.7.0
)

replace github.com/chatex-com/sdk-go => ../sdk-go
