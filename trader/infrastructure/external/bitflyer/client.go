package bitflyer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
	baseUrl, err := url.Parse(baseURL)
	if err != nil {
		return
	}
	apiURL, err := url.Parse(path)
	if err != nil {
		return
	}
	endpoint := baseUrl.ResolveReference(apiURL).String()
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
