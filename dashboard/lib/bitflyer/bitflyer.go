package bitflyer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const baseURL = "https://api.bitflyer.com/v1/"

type Client struct {
	key        string
	secret     string
	httpClient *http.Client
}

func NewClient(key, secret string) *Client {
	c := &Client{
		key:        key,
		secret:     secret,
		httpClient: &http.Client{},
	}
	return c
}

func (c *Client) header(method, path string, body []byte) map[string]string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	text := timestamp + method + path + string(body)

	mac := hmac.New(sha256.New, []byte(c.secret))
	mac.Write([]byte(text))
	sign := hex.EncodeToString(mac.Sum(nil))

	return map[string]string{
		"ACCESS-KEY":       c.key,
		"ACCESS-TIMESTAMP": timestamp,
		"ACCESS-SIGN":      sign,
		"Content-Type":     "application/json",
	}
}

func (c *Client) doRequest(method, path string, query map[string]string, body []byte) (respBody []byte, err error) {
	baseURL, err := url.Parse(baseURL)
	if err != nil {
		return
	}
	apiURL, err := url.Parse(path)
	if err != nil {
		return
	}
	endpoint := baseURL.ResolveReference(apiURL).String()
	fmt.Printf("[doRequest] %s %s\n", method, endpoint)

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return
	}
	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	for key, value := range c.header(method, req.URL.RequestURI(), body) {
		req.Header.Add(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

type Ticker struct {
	ProductCode     string  `json:"product_code"`
	State           string  `json:"state"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	MarketBidSize   float64 `json:"market_bid_size"`
	MarketAskSize   float64 `json:"market_ask_size"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

func (c *Client) GetTicker(productCode string) (*Ticker, error) {
	path := "ticker"
	query := map[string]string{"product_code": productCode}
	resp, err := c.doRequest("GET", path, query, nil)
	if err != nil {
		return nil, err
	}

	var ticker Ticker
	err = json.Unmarshal(resp, &ticker)
	if err != nil {
		return nil, err
	}
	if ticker.State != "RUNNING" {
		return nil, errors.New("bitflyer is not running")
	}
	return &ticker, nil
}

func (t *Ticker) GetMidPrice() float64 {
	return (t.BestBid + t.BestAsk) / 2
}

type Balance struct {
	CurrencyCode string  `json:"currency_code"`
	Amount       float64 `json:"amount"`
	Available    float64 `json:"available"`
}

func (c *Client) GetBalance() ([]Balance, error) {
	path := "me/getbalance"
	resp, err := c.doRequest("GET", path, map[string]string{}, nil)
	if err != nil {
		return nil, err
	}

	var balance []Balance
	err = json.Unmarshal(resp, &balance)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// 利用可能な通貨量を返す
// 仮想通貨, 現金の順
func (c *Client) GetAvailableBalance(CoinCode, CurrencyCode string) (float64, float64) {
	balances, err := c.GetBalance()
	if err != nil {
		return 0.0, 0.0
	}

	availableCoin, availableCurrency := 0.0, 0.0
	for _, balance := range balances {
		if balance.CurrencyCode == CoinCode {
			availableCoin = balance.Available
		} else if balance.CurrencyCode == CurrencyCode {
			availableCurrency = balance.Available
		}
	}
	return availableCoin, availableCurrency
}

// 新規注文リクエスト
// 注文一覧レスポンス
type Order struct {
	ID                     int     `json:"id"`
	ProductCode            string  `json:"product_code"`
	ChildOrderType         string  `json:"child_order_type"`
	Side                   string  `json:"side"`
	Price                  float64 `json:"price"`
	AveragePrice           int     `json:"average_price"`
	Size                   float64 `json:"size"`
	MinuteToExpires        int     `json:"minute_to_expire"`
	TimeInForce            string  `json:"time_in_force"`
	ChildOrderID           string  `json:"child_order_id"`
	ChildOrderState        string  `json:"child_order_state"`
	ExpireDate             string  `json:"expire_date"`
	ChildOrderDate         string  `json:"child_order_date"`
	ChildOrderAcceptanceID string  `json:"child_order_acceptance_id"`
	OutstandingSize        float64 `json:"outstanding_size"`
	CancelSize             float64 `json:"cancel_size"`
	ExecutedSize           float64 `json:"executed_size"`
	TotalCommission        float64 `json:"total_commission"`
}

type ResponseSendChildOrder struct {
	ChildOrderAcceptanceID string `json:"child_order_acceptance_id"`
}

// 新規注文を出す
func (c *Client) SendOrder(order *Order) (*ResponseSendChildOrder, error) {
	// json化
	data, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	// リクエスト送信
	url := "me/sendchildorder"
	resp, err := c.doRequest("POST", url, map[string]string{}, data)
	if err != nil {
		return nil, err
	}

	// レスポンス処理
	var response ResponseSendChildOrder
	if err = json.Unmarshal(resp, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) ListOrder(query map[string]string) ([]Order, error) {
	// リクエスト送信
	resp, err := c.doRequest("GET", "me/getchildorders", query, nil)
	if err != nil {
		return nil, err
	}

	// レスポンス処理
	var responseListOrder []Order
	if err = json.Unmarshal(resp, &responseListOrder); err != nil {
		return nil, err
	}
	return responseListOrder, nil
}
