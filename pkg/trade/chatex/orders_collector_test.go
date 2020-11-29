package chatex

import (
	"context"
	"testing"

	chatexsdk "github.com/chatex-com/sdk-go"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"

	"github.com/rwlist/autotrade-bot/pkg/conf"
	"github.com/rwlist/autotrade-bot/pkg/store/redisdb"
)

func TestOrdersCollector_CollectAll(t *testing.T) {
	cfg, err := conf.ParseEnv()
	assert.Nil(t, err)

	cli := chatexsdk.NewClient("https://api.chatex.com/v1", cfg.Chatex.RefreshToken)

	col := NewOrdersCollector(cli, nil)
	res, err := col.CollectAll()
	assert.Nil(t, err)

	spew.Dump(res)
}

func TestOrdersCollector_CollectInf(t *testing.T) {
	cfg, err := conf.ParseEnv()
	assert.Nil(t, err)

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
	})

	cli := chatexsdk.NewClient("https://api.chatex.com/v1", cfg.Chatex.RefreshToken)

	snapshotList := redisdb.NewList("chatex_order_snapshots", rdb)

	col := NewOrdersCollector(cli, snapshotList)

	_ = col.CollectInf(context.Background())
}
