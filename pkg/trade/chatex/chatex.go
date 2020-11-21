package chatex

import (
	"context"

	chatexsdk "github.com/chatex-com/sdk-go"
	"github.com/shopspring/decimal"

	"github.com/rwlist/autotrade-bot/pkg/draw"
	"github.com/rwlist/autotrade-bot/pkg/trade"
)

type Chatex struct {
	cli *chatexsdk.Client
}

func NewChatex(cli *chatexsdk.Client) *Chatex {
	return &Chatex{
		cli: cli,
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
	// TODO: implement
	return decimal.Zero, nil
}

func (c *Chatex) GetRate(symbol ...string) (decimal.Decimal, error) {
	// TODO: implement
	return decimal.Zero, nil
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
