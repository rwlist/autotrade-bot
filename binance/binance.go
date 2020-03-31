package binance

import (
	"context"
	"github.com/rwlist/autotrade-bot/draw"
	"math"
	"strings"
	"time"

	"github.com/rwlist/autotrade-bot/conf"

	"github.com/adshao/go-binance"
	"github.com/rwlist/autotrade-bot/to_str"
)

type MyBinance struct {
	client *binance.Client
}

func NewMyBinance(cfg conf.Binance, debug bool) *MyBinance {
	cli := binance.NewClient(cfg.APIKey, cfg.Secret)
	cli.Debug = debug
	return &MyBinance{cli}
}

func (b *MyBinance) AccountBalance() ([]binance.Balance, error) {
	info, err := b.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, err
	}
	return info.Balances, err
}

func (b *MyBinance) AccountSymbolBalance(symbol string) (float64, error) {
	info, err := b.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return 0, err
	}
	for _, bal := range info.Balances {
		if bal.Asset == symbol {
			return sum(bal.Free, bal.Locked), nil
		}
	}
	return 0, nil
}

func (b *MyBinance) BalanceToUSD(bal *binance.Balance) (float64, error) {
	haveFree := to_str.StrToFloat64(bal.Free)
	haveLocked := to_str.StrToFloat64(bal.Locked)
	if bal.Asset == "USDT" {
		return haveFree + haveLocked, nil
	}

	symbolPrice, err := b.client.NewListPricesService().Symbol(bal.Asset + "USDT").Do(context.Background())
	if err != nil {
		return 0, err
	}
	price := to_str.StrToFloat64(symbolPrice[0].Price)
	haveFree *= price
	haveLocked *= price
	return haveFree + haveLocked, nil
}

func (b *MyBinance) GetRate() (string, error) {
	symbolPrice, err := b.client.NewListPricesService().Symbol("BTCUSDT").Do(context.Background())
	if err != nil {
		return "", err
	}
	return symbolPrice[0].Price, nil
}

func (b *MyBinance) BuyAll() (Order, error) {
	price, err := b.GetRate()
	if err != nil {
		return nil, err
	}
	usdt, err := b.AccountSymbolBalance("USDT")
	if err != nil {
		return nil, err
	}
	quantity := usdt / to_str.StrToFloat64(price)
	order, err := b.client.NewCreateOrderService().Symbol("BTCUSDT").
		Side(binance.SideTypeBuy).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Price(price).Quantity(to_str.Float64ToStr(quantity, 6)).Do(context.Background())
	return &OrderNew{order}, err
}

func (b *MyBinance) SellAll() (Order, error) {
	price, err := b.GetRate()
	if err != nil {
		return nil, err
	}
	quantity, err := b.AccountSymbolBalance("BTC")
	if err != nil {
		return nil, err
	}
	order, err := b.client.NewCreateOrderService().Symbol("BTCUSDT").
		Side(binance.SideTypeSell).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Price(price).Quantity(to_str.Float64ToStr(quantity, 6)).Do(context.Background())
	return &OrderNew{order}, err
}

func (b *MyBinance) GetOrder(id int64) (Order, error) {
	order, err := b.client.NewGetOrderService().Symbol("BTCUSDT").
		OrderID(id).Do(context.Background())
	return &OrderExist{order}, err
}

func (b *MyBinance) CancelOrder(id int64) error {
	_, err := b.client.NewCancelOrderService().Symbol("BTCUSDT").
		OrderID(id).Do(context.Background())
	return err
}

func (b *MyBinance) GetKlines() (draw.Klines, error) {
	klines, err := b.client.
		NewKlinesService().Symbol("BTCUSDT").
		Interval("15m").
		StartTime(int64(1000) * (time.Now().Add(-time.Hour * 24).Unix())).
		Do(context.Background())
	if err != nil {
		return draw.Klines{}, err
	}

	var result draw.Klines

	// Extracting data from response
	result.Min = 1000000000.
	result.Max = -1.
	for _, val := range klines {
		result.Klines = append(result.Klines, draw.KlineTOHLCV{
			T: val.CloseTime / 1000,
			O: to_str.StrToFloat64(val.Open),
			H: to_str.StrToFloat64(val.High),
			L: to_str.StrToFloat64(val.Low),
			C: to_str.StrToFloat64(val.Close),
			V: to_str.StrToFloat64(val.Volume),
		})
		result.Min = math.Min(result.Min, to_str.StrToFloat64(val.Low))
		result.Max = math.Max(result.Max, to_str.StrToFloat64(val.High))
	}
	result.Last = to_str.StrToFloat64(klines[len(klines)-1].Close)
	result.StartTime = float64(klines[0].OpenTime / 1000)
	return result, nil
}

func sum(str1, str2 string) float64 {
	return to_str.StrToFloat64(str1) + to_str.StrToFloat64(str2)
}

func IsEmptyBalance(str string) bool {
	return strings.Trim(str, ".0") == ""
}
