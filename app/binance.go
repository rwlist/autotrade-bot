package app

import (
	"context"
	"github.com/adshao/go-binance"
	"log"
	"strings"
)

type MyBinance struct {
	client *binance.Client
}

func NewMyBinance(c *binance.Client) *MyBinance {
	return &MyBinance{client : c}
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
	haveFree := strToFloat64(bal.Free)
	haveLocked := strToFloat64(bal.Locked)
	if bal.Asset == "USDT" {
		return haveFree + haveLocked, nil
	}

	symbolPrice, err := b.client.NewListPricesService().Symbol(bal.Asset + "USDT").Do(context.Background())
	if err != nil {
		return 0, err
	}
	price := strToFloat64(symbolPrice[0].Price)
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

func (b *MyBinance) BuyAll() (*binance.CreateOrderResponse, error) {
	price, err := b.GetRate()
	if err != nil {
		return nil, err
	}
	usdt, err := b.AccountSymbolBalance("USDT")
	if err != nil {
		return nil, err
	}
	quantity := usdt / strToFloat64(price)
	order, err := b.client.NewCreateOrderService().Symbol("BTCUSDT").
		Side(binance.SideTypeBuy).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Price(price).Quantity(float64ToStr(quantity, 6)).Do(context.Background())
	return order, err
}

func (b *MyBinance) SellAll() (*binance.CreateOrderResponse, error) {
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
		TimeInForce(binance.TimeInForceTypeGTC).Price(price).Quantity(float64ToStr(quantity, 6)).Do(context.Background())
	return order, err
}

func (b *MyBinance) GetOrder(id int64) (*binance.Order, error) {
	order, err := b.client.NewGetOrderService().Symbol("BTCUSDT").
		OrderID(id).Do(context.Background())
	return order, err
}

func (b *MyBinance) CancelOrder(id int64) error {
	_, err := b.client.NewCancelOrderService().Symbol("BTCUSDT").
		OrderID(id).Do(context.Background())
	return err
}

func sum(str1, str2 string) float64 {
	return strToFloat64(str1) + strToFloat64(str2)
}

func isEmptyBalance(str string) bool {
	return strings.Trim(str, ".0") == ""
}

//------------------------TEST_BUY_COMMAND------------------------------------------
func (b *MyBinance) TestBuyAll() error {
	price, err := b.GetRate()
	if err != nil {
		return err
	}
	usdt, err := b.AccountSymbolBalance("USDT")
	if err != nil {
		return err
	}
	quantity := usdt / strToFloat64(price)
	err = b.client.NewCreateOrderService().Symbol("BTCUSDT").
		Side(binance.SideTypeBuy).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Price(price).Quantity(float64ToStr(quantity, 6)).Test(context.Background())
	log.Println(err)
	return err
}
//---------------------------------------------------------------------------------------
