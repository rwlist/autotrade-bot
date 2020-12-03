package redisdb

type HashValue map[string]string

func (h HashValue) Get(key string) string {
	return h[key]
}
