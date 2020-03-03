package app

import (
	"encoding/json"
	"net/http"
)

func binanceUrl() string {
	return "https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT"
}

const (
	errorValue = ""
)

func binanceRateQuery() (string, error) {
	resp, err := http.Get(binanceUrl())
	defer resp.Body.Close()
	if err != nil {
		return errorValue, err
	}
	return binanceAnswerParse(resp)
}

func binanceAnswerParse(resp *http.Response) (string, error) {
	type Ticker struct {
		Rate string `json:"price"`
	}
	var dec Ticker
	err := json.NewDecoder(resp.Body).Decode(&dec)
	if err != nil {
		return errorValue, err
	}
	return dec.Rate, nil
}
