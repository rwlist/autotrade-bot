package chatex

import (
	"context"
	"errors"
	"strings"

	chatexsdk "github.com/chatex-com/sdk-go"
	"github.com/rwlist/autotrade-bot/pkg/money"

	"github.com/shopspring/decimal"

	"github.com/rwlist/autotrade-bot/pkg/draw"
	"github.com/rwlist/autotrade-bot/pkg/trade"
)

type Chatex struct {
	cli       *chatexsdk.Client
	collector *OrdersCollector
}

func NewChatex(cli *chatexsdk.Client, collector *OrdersCollector) *Chatex {
	return &Chatex{
		cli:       cli,
		collector: collector,
	}
}

func (c *Chatex) AccountBalance() ([]trade.Balance, error) {
	resp, err := c.cli.GetMyBalance(context.Background())
	if err != nil {
		return nil, err
	}

	return convertBalanceSlice(resp), nil
}

func (c *Chatex) AccountSymbolBalance(symbol string) (decimal.Decimal, error) {
	panic("implement me")
}

func (c *Chatex) BalanceToUSD(bal *trade.Balance) (decimal.Decimal, error) {
	amount := decimal.Sum(bal.Free, bal.Locked)

	if strings.EqualFold(bal.Asset, "usdt") {
		return amount, nil
	}

	rate, err := c.GetRate(bal.Asset + "usdt")
	if err != nil {
		return decimal.Zero, err
	}

	return amount.Mul(rate), nil
}

func (c *Chatex) GetRate(symbols ...string) (decimal.Decimal, error) {
	if symbols == nil {
		symbols = []string{"BTCUSDT"}
	}

	if len(symbols) != 1 {
		return decimal.Zero, errors.New("invalid arguments")
	}

	last, err := c.collector.Last()
	if err != nil {
		return decimal.Zero, err
	}

	coins := last.Coins
	for _, c1 := range coins {
		for _, c2 := range coins {
			if !strings.EqualFold(c1.Code+c2.Code, symbols[0]) {
				continue
			}

			orders := last.Fetched[c2.Code+"/"+c1.Code].Orders
			if len(orders) == 0 {
				continue
			}

			return money.One.DivRound(orders[0].Rate, money.Precision), nil
		}
	}

	return decimal.Zero, errors.New("not found")
}

func (c *Chatex) BuyAll(symbol ...string) *trade.Status {
	panic("implement me")
}

func (c *Chatex) SellAll(symbol ...string) *trade.Status {
	panic("implement me")
}

func (c *Chatex) GetOrder(id int64, symbol ...string) (*trade.Order, error) {
	panic("implement me")
}

func (c *Chatex) CancelOrder(id int64, symbol ...string) error {
	panic("implement me")
}

func (c *Chatex) GetKlines(opts ...draw.KlinesOpts) (*draw.Klines, error) {
	panic("implement me")
}

func (c *Chatex) SetScale(scale string) {
	panic("implement me")
}

func (c *Chatex) SetSymbol(symbol string) {
	panic("implement me")
}
