package binance

import (
	"context"
	"math"
	"time"

	"github.com/rwlist/autotrade-bot/pkg/conf"
	"github.com/rwlist/autotrade-bot/pkg/draw"
	"github.com/rwlist/autotrade-bot/pkg/tostr"

	goBinance "github.com/adshao/go-binance"
)

type Binance struct {
	client *goBinance.Client
}

/*
	Создаёт новый Binance
*/
func NewBinance(cfg conf.Binance, debug bool) *Binance {
	cli := goBinance.NewClient(cfg.APIKey, cfg.Secret)
	cli.Debug = debug
	return &Binance{cli}
}

/*
	Возвращает информацию по балансу пользователя
	Вернулась ли информация по конкретной валюте непонятно от чего зависит
	Возможно возвращается для когда-либо использованных пользователем валют
*/
func (b *Binance) AccountBalance() ([]goBinance.Balance, error) {
	info, err := b.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, err
	}
	return info.Balances, err
}

/*
	Возвращает баланс для конкретной валюты
*/
func (b *Binance) AccountSymbolBalance(symbol string) (float64, error) {
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

/*
	Получает баланс какой-то валюты, смотрит на курс валюты к USDT и возвращает баланс в USDT
*/
func (b *Binance) BalanceToUSD(bal *goBinance.Balance) (float64, error) {
	haveFree := tostr.StrToFloat64(bal.Free)
	haveLocked := tostr.StrToFloat64(bal.Locked)
	if bal.Asset == "USDT" {
		return haveFree + haveLocked, nil
	}

	symbolPrice, err := b.client.NewListPricesService().Symbol(bal.Asset + "USDT").Do(context.Background())
	if err != nil {
		return 0, err
	}
	price := tostr.StrToFloat64(symbolPrice[0].Price)
	haveFree *= price
	haveLocked *= price
	return haveFree + haveLocked, nil
}

/*
	Возвращает текущий курс symbol[0] (default BTCUSDT)
*/
func (b *Binance) GetRate(symbol ...string) (string, error) {
	if len(symbol) == 0 {
		symbol = append(symbol, "BTCUSDT")
	}
	symbolPrice, err := b.client.NewListPricesService().Symbol(symbol[0]).Do(context.Background())
	if err != nil {
		return "", err
	}
	return symbolPrice[0].Price, nil
}

/*
	Закупается symbol[0] (default BTC) на все symbol[1] (default USDT)
*/
func (b *Binance) BuyAll(symbol ...string) (*Order, error) {
	if len(symbol) == 0 {
		symbol = append(symbol, "BTC", "USDT")
	}
	price, err := b.GetRate()
	if err != nil {
		return nil, err
	}
	usdt, err := b.AccountSymbolBalance(symbol[1])
	if err != nil {
		return nil, err
	}
	quantity := usdt / tostr.StrToFloat64(price)
	order, err := b.client.NewCreateOrderService().Symbol(symbol[0] + symbol[1]).
		Side(goBinance.SideTypeBuy).Type(goBinance.OrderTypeLimit).
		TimeInForce(goBinance.TimeInForceTypeGTC).Price(price).Quantity(tostr.Float64ToStr(quantity, 6)).Do(context.Background())
	return convertCreateOrderResponseToOrder(order), err
}

/*
	Продаёт все symbol[0] (default BTC) за symbol[1] (default USDT)
*/
func (b *Binance) SellAll(symbol ...string) (*Order, error) {
	if len(symbol) == 0 {
		symbol = append(symbol, "BTC", "USDT")
	}
	price, err := b.GetRate()
	if err != nil {
		return nil, err
	}
	quantity, err := b.AccountSymbolBalance(symbol[0])
	if err != nil {
		return nil, err
	}
	order, err := b.client.NewCreateOrderService().Symbol(symbol[0] + symbol[1]).
		Side(goBinance.SideTypeSell).Type(goBinance.OrderTypeLimit).
		TimeInForce(goBinance.TimeInForceTypeGTC).Price(price).Quantity(tostr.Float64ToStr(quantity, 6)).Do(context.Background())
	return convertCreateOrderResponseToOrder(order), err
}

/*
	Получает информацию по ордеру с данным id
*/
func (b *Binance) GetOrder(id int64) (*Order, error) {
	order, err := b.client.NewGetOrderService().Symbol("BTCUSDT").
		OrderID(id).Do(context.Background())
	return convertOrderToOrder(order), err
}

/*
	Закрывает ордер
*/
func (b *Binance) CancelOrder(id int64) error {
	_, err := b.client.NewCancelOrderService().Symbol("BTCUSDT").
		OrderID(id).Do(context.Background())
	return err
}

const timeShift = 1000
const hday = 24

/*
	Получает информацию по свечам symbol[0] (default BTCUSDT)
*/
func (b *Binance) GetKlines(symbol ...string) (draw.Klines, error) {
	if len(symbol) == 0 {
		symbol = append(symbol, "BTCUSDT")
	}
	klines, err := b.client.
		NewKlinesService().Symbol(symbol[0]).
		Interval("15m").
		StartTime(int64(timeShift) * (time.Now().Add(-time.Hour * hday).Unix())).
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
			T: val.CloseTime / timeShift,
			O: tostr.StrToFloat64(val.Open),
			H: tostr.StrToFloat64(val.High),
			L: tostr.StrToFloat64(val.Low),
			C: tostr.StrToFloat64(val.Close),
			V: tostr.StrToFloat64(val.Volume),
		})
		result.Min = math.Min(result.Min, tostr.StrToFloat64(val.Low))
		result.Max = math.Max(result.Max, tostr.StrToFloat64(val.High))
	}
	result.Last = tostr.StrToFloat64(klines[len(klines)-1].Close)
	result.StartTime = float64(klines[0].OpenTime / timeShift)
	return result, nil
}

func sum(str1, str2 string) float64 {
	return tostr.StrToFloat64(str1) + tostr.StrToFloat64(str2)
}
