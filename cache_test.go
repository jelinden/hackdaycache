package hackdaycache

import (
	"strings"
	"testing"
	"time"

	"github.com/jelinden/hackdaycache/client"
	"github.com/stretchr/testify/assert"
)

func TestSetItem(t *testing.T) {
	url := "https://httpbin.org/ip"
	i := CacheItem{
		Key:          url,
		Value:        nil,
		Expire:       time.Now(),
		UpdateLength: 1 * time.Second,
		GetFunc:      client.DataFetch,
	}
	AddItem(i)
	time.Sleep(2 * time.Second)
	assert.True(t, strings.Contains(string(GetItem(url)), "origin"))
	cache.Remove(url)
	assert.False(t, strings.Contains(string(GetItem(url)), "origin"))
	assert.True(t, string(GetItem(url)) == "")
}

func TestSetItemWithParams(t *testing.T) {
	i := CacheItem{
		Key:          p1,
		Value:        nil,
		Expire:       time.Now(),
		UpdateLength: 1 * time.Second,
		GetFunc:      withParams,
		FuncParams:   []string{p1, p2},
	}
	AddItem(i)
	time.Sleep(2 * time.Second)
	assert.True(t, strings.Contains(string(GetItem(p1)), "test1"))
	cache.Remove(p1)
	assert.False(t, strings.Contains(string(GetItem(p1)), "test1"))
	assert.True(t, string(GetItem(p1)) == "")
}

var p1 = "test"
var p2 = "1"

func withParams(key string, params ...string) []byte {
	return append([]byte(p1), []byte(p2)...)
}
