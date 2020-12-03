package redisdb

import (
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

type HashValue map[string]string

func (h HashValue) Get(key string) string {
	return h[key]
}

func (h HashValue) Decimal(key string) decimal.Decimal {
	val := h.Get(key)
	if val == "" {
		return decimal.Zero
	}

	dec, err := decimal.NewFromString(val)
	if err != nil {
		log.WithField("dec", dec).WithError(err).Error("failed to parse decimal")
		return decimal.Zero
	}

	return dec
}
