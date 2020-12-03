package redisdb

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	h := NewHash("test:hash", redisCli(t))

	assert.Nil(t, h.Set("abc", "123"))

	res, err := h.Read()
	assert.Nil(t, err)

	assert.Equal(t, "123", res.Get("abc"))

	spew.Dump(res)
}
