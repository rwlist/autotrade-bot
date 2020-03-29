package binance

import (
	"context"
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

type KlineTOHLCV struct {
	T int64
	O float64
	H float64
	L float64
	C float64
	V float64
}

type TOHLCVs []KlineTOHLCV

func (TOHLCV TOHLCVs) Len() int {
	return len(TOHLCV)
}

func (TOHLCV *TOHLCVs) Last() KlineTOHLCV {
	return (*TOHLCV)[TOHLCV.Len()-1]
}

func (TOHLCV TOHLCVs) TOHLCV(i int) (float64, float64, float64, float64, float64, float64) {
	return float64(TOHLCV[i].T), TOHLCV[i].O, TOHLCV[i].H, TOHLCV[i].L, TOHLCV[i].C, TOHLCV[i].V
}

func (b *MyBinance) GetKlines() (TOHLCVs, float64, float64, float64, float64, error) {
	klines, err := b.client.
		NewKlinesService().Symbol("BTCUSDT").
		Interval("15m").
		StartTime(int64(1000) * (time.Now().Add(-time.Hour * 24).Unix())).
		Do(context.Background())
	if err != nil {
		return nil, .0, .0, .0, .0, err
	}

	var result TOHLCVs

	// Extracting data from response
	min := 1000000000.
	max := -1.
	for _, val := range klines {
		result = append(result, KlineTOHLCV{
			T: val.CloseTime / 1000,
			O: to_str.StrToFloat64(val.Open),
			H: to_str.StrToFloat64(val.High),
			L: to_str.StrToFloat64(val.Low),
			C: to_str.StrToFloat64(val.Close),
			V: to_str.StrToFloat64(val.Volume),
		})
		min = math.Min(min, to_str.StrToFloat64(val.Low))
		max = math.Max(max, to_str.StrToFloat64(val.High))
	}
	return result, result.Last().C, min, max, float64(klines[0].OpenTime / 1000), nil
}

func sum(str1, str2 string) float64 {
	return to_str.StrToFloat64(str1) + to_str.StrToFloat64(str2)
}

func IsEmptyBalance(str string) bool {
	return strings.Trim(str, ".0") == ""
}
