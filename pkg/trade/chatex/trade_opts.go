package chatex

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/rwlist/autotrade-bot/pkg/store/redisdb"
)

type TradeOpts struct {
	hash *redisdb.Hash
}

func NewTradeOpts(hash *redisdb.Hash) *TradeOpts {
	return &TradeOpts{
		hash: hash,
	}
}

func (o *TradeOpts) GetAll() (redisdb.HashValue, error) {
	return o.hash.Read()
}

func (o *TradeOpts) GetSingle(key string) (string, error) {
	val, err := o.GetAll()
	if err != nil {
		return "", err
	}

	return val.Get(key), nil
}

func (o *TradeOpts) SetOption(key, val string) error {
	return o.hash.Set(key, val)
}

func (o *TradeOpts) FetchCollectorPeriod() time.Duration {
	const (
		def = time.Minute
		min = time.Second
	)

	str, err := o.GetSingle("chatex.collector.period")
	if err != nil {
		log.WithError(err).Error("failed to fetch collector period")
		return def
	}

	if str == "" {
		return def
	}

	dur, err := time.ParseDuration(str)
	if err != nil {
		log.WithError(err).Error("failed to parse collector period")
		return def
	}

	if min > dur {
		return min
	}

	return dur
}
