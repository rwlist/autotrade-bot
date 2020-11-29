package chatexsdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
)

type PaymentSystem struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type PaymentSystemEstimation struct {
	PaymentSystem PaymentSystem   `json:"payment_system"`
	FiatAmount    decimal.Decimal `json:"estimated_fiat_amount"`
}

type FiatEstimation struct {
	Fiat    Fiat                      `json:"fiat"`
	Systems []PaymentSystemEstimation `json:"estimates"`
}

type Currency struct {
	Code     string `json:"name"`
	Name     string `json:"full_name"`
	Decimals uint64 `json:"decimals"`
}

type Fiat Currency

type Coin Currency

func (c *Client) GetCoins(ctx context.Context) ([]Coin, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/coins", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	var coins []Coin
	if err := c.sendRequest(ctx, req, &coins); err != nil {
		return nil, err
	}

	return coins, nil
}

func (c *Client) GetCoin(ctx context.Context, code string) (*Coin, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/coins/%s", c.baseURL, code), nil)
	if err != nil {
		return nil, err
	}

	var coin Coin
	if err := c.sendRequest(ctx, req, &coin); err != nil {
		return nil, err
	}

	return &coin, nil
}

func (c *Client) GetPaymentSystem(ctx context.Context, id uint) (*PaymentSystem, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/payment-systems/%d", c.baseURL, id), nil)
	if err != nil {
		return nil, err
	}

	var result PaymentSystem
	if err := c.sendRequest(ctx, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) Estimate(ctx context.Context, coin string, amount decimal.Decimal) ([]FiatEstimation, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/payment-systems/estimated?currency=%s&amount=%s", c.baseURL, coin, amount.String()), nil)
	if err != nil {
		return nil, err
	}

	var result []FiatEstimation
	if err := c.sendRequest(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
