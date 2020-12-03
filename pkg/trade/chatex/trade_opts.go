package chatex

import "github.com/rwlist/autotrade-bot/pkg/store/redisdb"

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
