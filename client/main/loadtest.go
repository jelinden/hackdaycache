package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	hc "github.com/jelinden/hackdaycache"
	"github.com/jelinden/hackdaycache/client"
	"github.com/julienschmidt/httprouter"
)

func main() {
	url := "https://httpbin.org/ip"
	i := hc.CacheItem{
		Key:          url,
		Value:        client.DataFetch(url),
		Expire:       time.Now(),
		UpdateLength: 10 * time.Second,
		GetFunc:      client.DataFetch,
		InUse:        false,
	}
	hc.AddItem(i)

	router := httprouter.New()
	router.RedirectFixedPath = true
	router.RedirectTrailingSlash = true
	router.GET("/httpbin", Bin)
	go fetch()
	http.ListenAndServe(":8800", router)
}

func Bin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	item := hc.GetItem("https://httpbin.org/ip")
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(item)
}

func fetch() {
	t := time.Now()
	const howMany = 100000
	for i := 0; i < howMany; i++ {
		item := client.DataFetch("http://localhost:8800/httpbin")
		if !strings.Contains(string(item), "origin") {
			os.Exit(2)
		}
	}
	log.Println("fetched", howMany, "times from cache successfully in", time.Now().Sub(t))
	os.Exit(1)
}
