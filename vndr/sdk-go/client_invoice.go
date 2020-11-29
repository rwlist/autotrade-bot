package chatexsdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

const (
	InvoiceUnassigned InvoiceStatus = "UNASSIGNED"
	InvoiceActive     InvoiceStatus = "ACTIVE"
	InvoiceCompleted  InvoiceStatus = "COMPLETED"
	InvoiceCanceled   InvoiceStatus = "CANCELED"
)

type InvoiceStatus string

type Invoice struct {
	ID              string          `json:"id"`
	Status          InvoiceStatus   `json:"status"`
	CreatedAt       time.Time       `json:"created_at"`
	Amount          decimal.Decimal `json:"amount"`
	CurrencyName    string          `json:"currency"`
	FiatName        string          `json:"fiat"`
	PaymentSystemID uint64          `json:"payment_system_id"`
	LangID          string          `json:"lang_id"`
	CountryCode     string          `json:"country_code"`
	CallbackURL     string          `json:"callback_url"`
}

type InvoiceOptions struct {
	Amount          decimal.Decimal `json:"amount"`
	CountryAlpha3   string          `json:"country_code"`
	Coin            string          `json:"currency"`
	Fiat            string          `json:"fiat"`
	LanguageAlpha2  string          `json:"lang_id"`
	PaymentSystemID uint            `json:"payment_system_id"`
}

type InvoiceResponse struct {
	Invoice         Invoice `json:"invoice"`
	UserRedirectURL string  `json:"user_redirect_url"`
}

func (c *Client) GetInvoices(ctx context.Context) ([]Invoice, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/invoices", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	var result []Invoice
	if err := c.sendRequest(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) GetInvoice(ctx context.Context, id string) (*Invoice, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/invoices/%s", c.baseURL, id), nil)
	if err != nil {
		return nil, err
	}

	var result Invoice
	if err := c.sendRequest(ctx, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) CreateInvoice(ctx context.Context, invoice InvoiceOptions) (*InvoiceResponse, error) {
	body, err := json.Marshal(invoice)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/invoices", c.baseURL), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	var result InvoiceResponse
	if err := c.sendRequest(ctx, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
