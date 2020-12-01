package chatexsdk

import (
	"context"
	"fmt"
	"net/http"
)

const (
	LevelNone     VerificationLevel = "NONE"
	LevelPhone    VerificationLevel = "VERIFIED_PHONE_NUMBER"
	LevelPassport VerificationLevel = "VERIFIED_PASSPORT"
	LevelAddress  VerificationLevel = "VERIFIED_ADDRESS"
)

type Me struct {
	ID           uint64        `json:"id"`
	Profile      Profile       `json:"profile"`
	MerchantInfo *MerchantInfo `json:"merchant_info,omitempty"`
}

type MerchantInfo struct {
	Name           string `json:"name"`
	MaxAmountLimit string `json:"usd_amount_max_limit"`
}

type VerificationLevel = string

type Verification struct {
	CurrentLevel VerificationLevel `json:"current_level"` // Maximum passed verification level.
}

type AML5Limits struct {
	CurrentTurnover string `json:"current_turnover"` // During the last month.
	CurrentWithdraw string `json:"current_withdraw"` // During the last month.

	TurnoverLimit      string `json:"turnover_limit"`       // Monthly limit.
	WithdrawLimit      string `json:"withdraw_limit"`       // Monthly limit.
	WithdrawLimitDaily string `json:"withdraw_limit_daily"` // Daily limit.
}

type Profile struct {
	Username string `json:"username"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email,omitempty"`

	LangID      string `json:"lang_id"`      // Alpha2 code.
	CountryCode string `json:"country_code"` // Alpha3 code.

	IsFinanceBlocked bool `json:"is_finance_blocked"`

	Verification Verification `json:"verification"`
	AML5Limits   AML5Limits   `json:"limits"`
}

type Balance struct {
	Coin   string `json:"coin"`   // Coin name
	Amount string `json:"amount"` // Available amount
	Held   string `json:"held"`   // Held amount
}

func (c *Client) GetMyProfile(ctx context.Context) (*Me, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/me", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	var me Me
	if err := c.sendRequest(ctx, req, &me); err != nil {
		return nil, err
	}

	return &me, nil
}

func (c *Client) GetMyBalance(ctx context.Context) ([]Balance, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/me/balance", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	var result []Balance
	if err := c.sendRequest(ctx, req, &result); err != nil {
		return nil, err
	}

	return result, nil
}
