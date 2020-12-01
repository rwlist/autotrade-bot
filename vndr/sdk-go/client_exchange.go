package chatexsdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

const (
	Active    OrderStatus = "ACTIVE"
	Inactive  OrderStatus = "INACTIVE"
	Canceled  OrderStatus = "CANCELED"
	Completed OrderStatus = "COMPLETED"
)

type OrderStatus string

type Order struct {
	ID            uint64          `json:"id"`
	Pair          string          `json:"pair"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	Rate          decimal.Decimal `json:"rate"`
	Amount        decimal.Decimal `json:"amount"`
	IsOwner       bool            `json:"is_owner,omitempty" `
	InitialAmount decimal.Decimal `json:"initial_amount,omitempty"`
	Status        OrderStatus     `json:"status"`
}

type CreateOrderRequest struct {
	Amount decimal.Decimal `json:"amount"`
	Pair   string          `json:"pair"`
	Rate   decimal.Decimal `json:"rate"`
}

type MyOrdersRequest struct {
	Pair   string
	Status OrderStatus
	Offset uint
	Limit  uint
}

type UpdateOrderRequest struct {
	Amount decimal.Decimal `json:"amount,omitempty"`
	Rate   decimal.Decimal `json:"rate,omitempty"`
}

type Trade struct {
	ID             uint64    `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Amount         string    `json:"amount"`
	ReceivedAmount string    `json:"received_amount"`
	Fee            string    `json:"fee"`
	Order          *Order    `json:"order"`
}

type TradeRequest struct {
	Amount decimal.Decimal `json:"amount"`
	Rate   decimal.Decimal `json:"rate"`
}

type TradesRequest struct {
	OrderID uint
	Offset  uint
	Limit   uint
}

func (c *Client) CreateOrder(ctx context.Context, req CreateOrderRequest) (*Order, error) {
	body, _ := json.Marshal(req)

	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/exchange/orders", c.baseURL), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var result Order
	if err := c.sendRequest(ctx, r, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetOrder(ctx context.Context, id uint) (*Order, error) {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/exchange/orders/%d", c.baseURL, id), nil)
	if err != nil {
		return nil, err
	}

	var result Order
	if err := c.sendRequest(ctx, r, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) ActivateOrder(ctx context.Context, id uint) (*Order, error) {
	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/exchange/orders/%d/activate", c.baseURL, id), nil)
	if err != nil {
		return nil, err
	}

	var result Order
	if err := c.sendRequest(ctx, r, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) DeactivateOrder(ctx context.Context, id uint) (*Order, error) {
	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/exchange/orders/%d/deactivate", c.baseURL, id), nil)
	if err != nil {
		return nil, err
	}

	var result Order
	if err := c.sendRequest(ctx, r, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) CancelOrder(ctx context.Context, id uint) (*Order, error) {
	r, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/exchange/orders/%d", c.baseURL, id), nil)
	if err != nil {
		return nil, err
	}

	var result Order
	if err := c.sendRequest(ctx, r, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetOrders(ctx context.Context, pair string, offset, limit uint) ([]Order, error) {
	uri, _ := url.Parse(fmt.Sprintf("%s/exchange/orders", c.baseURL))

	data := url.Values{}
	data.Set("pair", pair)
	if offset > 0 {
		data.Set("offset", strconv.Itoa(int(offset)))
	}
	if limit > 0 {
		data.Set("limit", strconv.Itoa(int(limit)))
	}

	uri.RawQuery = data.Encode()

	r, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, err
	}

	var result []Order
	if err := c.sendRequest(ctx, r, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) GetMyOrders(ctx context.Context, req MyOrdersRequest) ([]Order, error) {
	uri, _ := url.Parse(fmt.Sprintf("%s/exchange/orders/my", c.baseURL))

	data := url.Values{}
	if req.Pair != "" {
		data.Set("pair", req.Pair)
	}
	if req.Status != "" {
		data.Set("status", string(req.Status))
	}
	if req.Offset > 0 {
		data.Set("offset", strconv.Itoa(int(req.Offset)))
	}
	if req.Limit > 0 {
		data.Set("limit", strconv.Itoa(int(req.Limit)))
	}

	uri.RawQuery = data.Encode()

	r, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, err
	}

	var result []Order
	if err := c.sendRequest(ctx, r, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) UpdateOrder(ctx context.Context, id uint, req UpdateOrderRequest) (*Order, error) {
	b, _ := json.Marshal(req)

	r, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/exchange/orders/%d", c.baseURL, id), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	var result Order
	if err := c.sendRequest(ctx, r, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) CreateTrade(ctx context.Context, orderID uint, req TradeRequest) (*Trade, error) {
	b, _ := json.Marshal(req)

	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/exchange/orders/%d/trades", c.baseURL, orderID), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	var result Trade
	if err := c.sendRequest(ctx, r, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetTrade(ctx context.Context, tradeID uint) (*Trade, error) {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/exchange/orders/trades/%d", c.baseURL, tradeID), nil)
	if err != nil {
		return nil, err
	}

	var result Trade
	if err := c.sendRequest(ctx, r, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetTrades(ctx context.Context, req TradesRequest) ([]Trade, error) {
	uri, _ := url.Parse(fmt.Sprintf("%s/exchange/orders/trades", c.baseURL))

	data := url.Values{}
	if req.OrderID > 0 {
		data.Set("order_id", strconv.Itoa(int(req.OrderID)))
	}
	if req.Offset > 0 {
		data.Set("offset", strconv.Itoa(int(req.Offset)))
	}
	if req.Limit > 0 {
		data.Set("limit", strconv.Itoa(int(req.Limit)))
	}

	uri.RawQuery = data.Encode()

	r, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, err
	}

	var result []Trade
	if err := c.sendRequest(ctx, r, &result); err != nil {
		return nil, err
	}

	return result, nil
}
