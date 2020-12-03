package redisdb

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Hash struct {
	key string
	cli *redis.Client
}

func NewHash(key string, cli *redis.Client) *Hash {
	return &Hash{
		key: key,
		cli: cli,
	}
}

func (h *Hash) Read() (HashValue, error) {
	res, err := h.cli.HGetAll(context.Background(), h.key).Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *Hash) Set(key, value string) error {
	return h.cli.HSet(context.Background(), h.key, key, value).Err()
}
