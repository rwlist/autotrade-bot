package binance

import (
	"math"
	"strings"
	"time"

	"github.com/rwlist/autotrade-bot/pkg/trade"

	"github.com/rwlist/autotrade-bot/pkg/conf"
	"github.com/rwlist/autotrade-bot/pkg/draw"
	"github.com/rwlist/autotrade-bot/pkg/tostr"

	goBinance "github.com/adshao/go-binance"
)

type Binance struct {
	client Client
}

/*
	Создаёт новый Binance
*/
func NewBinance(cfg conf.Binance, debug bool) Binance {
	var cli Client
	if debug {
		cli = NewClientLog(cfg.APIKey, cfg.Secret)
	} else {
		cli = NewClientDefault(cfg.APIKey, cfg.Secret)
	}
	return Binance{cli}
}

/*
	Возвращает информацию по балансу пользователя
	Вернулась ли информация по конкретной валюте непонятно от чего зависит
	Возможно возвращается для когда-либо использованных пользователем валют
*/
func (b *Binance) AccountBalance() ([]trade.Balance, error) {
	info, err := b.client.AccountBalance()
	if err != nil {
		return nil, err
	}
	return convertBalanceSlice(info), err
}

/*
	Возвращает баланс для конкретной валюты
*/
func (b *Binance) AccountSymbolBalance(symbol string) (float64, error) {
	info, err := b.client.AccountBalance()
	if err != nil {
		return 0, err
	}
	for _, bal := range info {
		if bal.Asset == symbol {
			return sum(bal.Free, bal.Locked), nil
		}
	}
	return 0, nil
}

/*
	Получает баланс какой-то валюты, смотрит на курс валюты к USDT и возвращает баланс в USDT
*/
func (b *Binance) BalanceToUSD(bal *trade.Balance) (float64, error) {
	haveFree := tostr.StrToFloat64(bal.Free)
	haveLocked := tostr.StrToFloat64(bal.Locked)

	if bal.Asset == "USDT" {
		return haveFree + haveLocked, nil
	}

	symbolPrice, err := b.client.ListPrices(bal.Asset + "USDT")
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
	symbolPrice, err := b.client.ListPrices(symbol[0])
	if err != nil {
		return "", err
	}
	return symbolPrice[0].Price, nil
}

/*
	Закупается symbol[0] (default BTC) на все symbol[1] (default USDT)
	Возвращает nil, true, nil если закуплено на все деньги
*/
func (b *Binance) BuyAll(symbol ...string) *trade.Status {
	if len(symbol) == 0 {
		symbol = append(symbol, "BTC", "USDT")
	}
	price, err := b.GetRate()
	if err != nil {
		return &trade.Status{
			Order: nil,
			Done:  false,
			Err:   err,
		}
	}
	usdt, err := b.AccountSymbolBalance(symbol[1])
	if err != nil {
		return &trade.Status{
			Order: nil,
			Done:  false,
			Err:   err,
		}
	}
	quantity := usdt / tostr.StrToFloat64(price)

	req := &orderReq{
		Symbol:   symbol[0] + symbol[1],
		Side:     goBinance.SideTypeBuy,
		Type:     goBinance.OrderTypeLimit,
		Tif:      goBinance.TimeInForceTypeGTC,
		Price:    price,
		Quantity: tostr.Float64ToStr(quantity, 6),
	}
	order, err := b.client.CreateOrder(req)

	if err != nil {
		if strings.Contains(err.Error(), "code=-1013") {
			return &trade.Status{
				Order: nil,
				Done:  true,
				Err:   nil,
			}
		}
		return &trade.Status{
			Order: nil,
			Done:  false,
			Err:   err,
		}
	}
	return &trade.Status{
		Order: convertCreateOrderResponseToOrder(order),
		Done:  false,
		Err:   nil,
	}
}

/*
	Продаёт все symbol[0] (default BTC) за symbol[1] (default USDT)
	Возвращает nil, true, nil если всё продано
*/
func (b *Binance) SellAll(symbol ...string) *trade.Status {
	if len(symbol) == 0 {
		symbol = append(symbol, "BTC", "USDT")
	}
	price, err := b.GetRate()
	if err != nil {
		return &trade.Status{
			Order: nil,
			Done:  false,
			Err:   err,
		}
	}
	quantity, err := b.AccountSymbolBalance(symbol[0])
	if err != nil {
		return &trade.Status{
			Order: nil,
			Done:  false,
			Err:   err,
		}
	}

	req := &orderReq{
		Symbol:   symbol[0] + symbol[1],
		Side:     goBinance.SideTypeSell,
		Type:     goBinance.OrderTypeLimit,
		Tif:      goBinance.TimeInForceTypeGTC,
		Price:    price,
		Quantity: tostr.Float64ToStr(quantity, 6),
	}
	order, err := b.client.CreateOrder(req)

	if err != nil {
		if strings.Contains(err.Error(), "code=-1013") {
			return &trade.Status{
				Order: nil,
				Done:  true,
				Err:   nil,
			}
		}
		return &trade.Status{
			Order: nil,
			Done:  false,
			Err:   err,
		}
	}
	return &trade.Status{
		Order: convertCreateOrderResponseToOrder(order),
		Done:  false,
		Err:   nil,
	}
}

/*
	Получает информацию по ордеру для пары symbol[0] ("BTCUSDT") с данным id
*/
func (b *Binance) GetOrder(id int64, symbol ...string) (*trade.Order, error) {
	if len(symbol) == 0 {
		symbol = append(symbol, "BTCUSDT")
	}

	req := &orderID{
		Symbol: symbol[0],
		ID:     id,
	}
	order, err := b.client.GetOrder(req)
	return convertOrderToOrder(order), err
}

/*
	Закрывает ордер с данным id для пары symbol[0] ("BTCUSDT")
*/
func (b *Binance) CancelOrder(id int64, symbol ...string) error {
	if len(symbol) == 0 {
		symbol = append(symbol, "BTCUSDT")
	}

	req := &orderID{
		Symbol: symbol[0],
		ID:     id,
	}
	_, err := b.client.CancelOrder(req)
	if err != nil {
		if !strings.Contains(err.Error(), "code=-2011") {
			return err
		}
	}
	return nil
}

const timeShift = 1000
const hday = 24

/*
	Получает информацию по свечам symbol[0] (default "BTCUSDT")
*/
func (b *Binance) GetKlines(symbol ...string) (draw.Klines, error) {
	if len(symbol) == 0 {
		symbol = append(symbol, "BTCUSDT")
	}

	req := &klinesReq{
		Symbol:    symbol[0],
		T:         "15m",
		StartTime: int64(timeShift) * (time.Now().Add(-time.Hour * hday).Unix()),
	}
	klines, err := b.client.GetKlines(req)
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
