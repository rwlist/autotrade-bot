package binance

import (
	"math"
	"strings"

	"github.com/rwlist/autotrade-bot/pkg/trade"

	"github.com/rwlist/autotrade-bot/pkg/draw"
	"github.com/rwlist/autotrade-bot/pkg/tostr"

	gobinance "github.com/adshao/go-binance"
)

type Binance struct {
	client Client
	opts   *trade.Opts
}

// Создаёт новый Binance
func NewBinance(cli Client) Binance {
	opts := trade.Opts{
		Symbol: "BTCUSDT", // default "BTCUSDT"
		Scale:  "15m",     // default "15m"
	}
	return Binance{cli, &opts}
}

//	Возвращает информацию по балансу пользователя
//	Вернулась ли информация по конкретной валюте непонятно от чего зависит
//	Возможно возвращается для когда-либо использованных пользователем валют
func (b *Binance) AccountBalance() ([]trade.Balance, error) {
	info, err := b.client.AccountBalance()
	if err != nil {
		return nil, err
	}
	return convertBalanceSlice(info), err
}

//	Возвращает баланс для конкретной валюты
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

//	Получает баланс какой-то валюты, смотрит на курс валюты к USDT и возвращает баланс в USDT
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

//	Возвращает текущий курс symbol[0] (default BTCUSDT)
func (b *Binance) GetRate(symbol ...string) (string, error) {
	if len(symbol) == 0 {
		symbol = append(symbol, b.opts.Symbol)
	}
	symbolPrice, err := b.client.ListPrices(symbol[0])
	if err != nil {
		return "", err
	}
	return symbolPrice[0].Price, nil
}

//	Закупается symbol[0] (default BTC) на все symbol[1] (default USDT)
//	Возвращает nil, true, nil если закуплено на все деньги
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
		Side:     gobinance.SideTypeBuy,
		Type:     gobinance.OrderTypeLimit,
		Tif:      gobinance.TimeInForceTypeGTC,
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

//	Продаёт все symbol[0] (default BTC) за symbol[1] (default USDT)
//	Возвращает nil, true, nil если всё продано
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
		Side:     gobinance.SideTypeSell,
		Type:     gobinance.OrderTypeLimit,
		Tif:      gobinance.TimeInForceTypeGTC,
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

//	Получает информацию по ордеру для пары symbol[0] ("BTCUSDT") с данным id
func (b *Binance) GetOrder(id int64, symbol ...string) (*trade.Order, error) {
	if len(symbol) == 0 {
		symbol = append(symbol, b.opts.Symbol)
	}

	req := &orderID{
		Symbol: symbol[0],
		ID:     id,
	}
	order, err := b.client.GetOrder(req)
	return convertOrderToOrder(order), err
}

//	Закрывает ордер с данным id для пары symbol[0] ("BTCUSDT")
func (b *Binance) CancelOrder(id int64, symbol ...string) error {
	if len(symbol) == 0 {
		symbol = append(symbol, b.opts.Symbol)
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

//	Получает информацию по свечам (default "BTCUSDT", "15m")
func (b *Binance) GetKlines(opts ...draw.KlinesOpts) (*draw.Klines, error) {
	if len(opts) == 0 {
		opts = append(opts, b.klinesOpts())
	}

	req := &klinesReq{
		Symbol:    opts[0].Symbol,
		T:         opts[0].T,
		StartTime: calcStartTime(opts[0].T),
	}
	klines, err := b.client.GetKlines(req)
	if err != nil {
		return nil, err
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
	return &result, nil
}

func (b *Binance) SetScale(scale string) {
	b.opts.Scale = scale
}

func (b *Binance) SetSymbol(symbol string) {
	b.opts.Symbol = symbol
}
