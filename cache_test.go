package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var client *http.Client

func init() {
	client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 64,
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 10 * time.Second,
			}).Dial,
			ResponseHeaderTimeout: 5 * time.Second,
		},
	}
}

func TestSetItem(t *testing.T) {
	url := "https://httpbin.org/ip"
	i := CacheItem{
		Key:          url,
		Value:        nil,
		Expire:       time.Now(),
		UpdateLength: 1 * time.Second,
		GetFunc:      dataFetch,
		InUse:        false,
	}
	AddItem(i)
	time.Sleep(2 * time.Second)
	assert.True(t, strings.Contains(string(GetItem(url)), "origin"))
	cache.Remove(url)
	assert.False(t, strings.Contains(string(GetItem(url)), "origin"))
	assert.True(t, string(GetItem(url)) == "")
}

func httpGet(url string) ([]byte, error) {
	t := time.Now()
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Println("GET", resp.StatusCode, url, time.Now().Sub(t), time.Now().Format("15:04:05.000"))
	return body, nil
}

func dataFetch(url string) []byte {
	data, err := httpGet(url)
	if err != nil {
		log.Println("getting", url, "failed", err.Error())
	}
	return data
}
