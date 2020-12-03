package redisdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

type someObject struct {
	Time time.Time
	Map  map[string]interface{}
}

func TestList_LPush(t *testing.T) {
	rdb := redisCli(t)

	list := NewList("temp_test", rdb)

	length, err := list.Len()
	assert.Nil(t, err)

	fmt.Println("Len", length)

	obj := someObject{
		Time: time.Now(),
		Map: map[string]interface{}{
			"abc":  "qwe",
			"val1": 228,
			"len":  length,
		},
	}

	err = list.LPush(obj)
	assert.Nil(t, err)

	spew.Dump("inserted", obj)

	var obj2 someObject
	err = list.Left(&obj2)
	assert.Nil(t, err)

	spew.Dump("read", obj2)
}
