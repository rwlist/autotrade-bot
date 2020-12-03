package redisdb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestAbc(t *testing.T) {
	rdb := redisCli(t)

	ts := time.Now()

	set(rdb, "k", 0)

	for i := 0; i < 20000; i++ {
		val := get(rdb, "k")
		now := time.Now()

		if now.Sub(ts) >= time.Second {
			fmt.Printf("Passed %s, iter %d\n", now.Sub(ts), i)
			ts = now
		}

		set(rdb, "k", val+1)
	}
}

func set(rdb *redis.Client, key string, val int) {
	res, err := rdb.Set(context.Background(), key, val, redis.KeepTTL).Result()
	if err != nil {
		panic(err)
	}
	if res != "OK" {
		panic(res)
	}
}

func get(rdb *redis.Client, key string) int {
	res, err := rdb.Get(context.Background(), key).Int()
	if err != nil {
		panic(err)
	}
	return res
}
