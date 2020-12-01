package redisdb

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type List struct {
	name string
	cli  *redis.Client
}

func NewList(name string, cli *redis.Client) *List {
	return &List{
		name: name,
		cli:  cli,
	}
}

func (l *List) Len() (int64, error) {
	return l.cli.LLen(context.Background(), l.name).Result()
}

func (l *List) LPush(obj interface{}) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	return l.cli.LPush(context.Background(), l.name, b).Err()
}

func (l *List) Left(obj interface{}) error {
	b, err := l.cli.LIndex(context.Background(), l.name, 0).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(b, obj)
}
