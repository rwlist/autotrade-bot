package redisdb

import (
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"

	"github.com/rwlist/autotrade-bot/pkg/conf"
)

func redisCli(t *testing.T) *redis.Client {
	cfg, err := conf.ParseEnv()
	assert.Nil(t, err)

	return redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
	})
}
