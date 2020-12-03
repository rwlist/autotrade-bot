package redisdb

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type Store struct {
	key string
	cli *redis.Client
}

func NewStore(key string, cli *redis.Client) *Store {
	return &Store{
		key: key,
		cli: cli,
	}
}

func (s *Store) Set(obj interface{}) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	return s.cli.Set(context.Background(), s.key, b, redis.KeepTTL).Err()
}

func (s *Store) GetOrNop(obj interface{}) error {
	b, err := s.cli.Get(context.Background(), s.key).Bytes()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}

	return json.Unmarshal(b, obj)
}
