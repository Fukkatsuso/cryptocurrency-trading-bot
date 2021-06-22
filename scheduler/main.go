package main

import (
	"log"
	"net/http"

	"github.com/robfig/cron/v3"
)

func traderFetchTicker() {
	url := "http://trading_trader:8080/fetch-ticker"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	log.Println("[cron]", resp.StatusCode, resp.Request.URL)
}

func traderTrade() {
	url := "http://trading_trader:8080/trade"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	log.Println("[cron]", resp.StatusCode, resp.Request.URL)
}

func main() {
	c := cron.New()
	c.AddFunc("*/5 * * * *", traderFetchTicker)
	c.AddFunc("*/10 * * * *", traderTrade)
	c.Start()

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
