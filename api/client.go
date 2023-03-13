package api

import (
	"bytes"
	"io"
	"net/http"
)

var (
	ChatUrl   string
	WeeklyUrl string
	ImageUrl  string
	GrantsUrl string
)

func init() {
	// chat
	ChatUrl = "https://jiu.aizhineng.cc/message.php"
	// weekly
	WeeklyUrl = "https://weeklyreport.avemaria.fun/api/generate"
	// image
	ImageUrl = "https://ai.mdzx.me/v1/images/generations"
	// 余额
	GrantsUrl = "https://ai.mdzx.me/dashboard/billing/credit_grants"
}

type Client struct {
	transport *http.Client
	url       string
}

func NewClient(url string) *Client {
	return &Client{
		transport: http.DefaultClient,
		url:       url,
	}
}

func (c *Client) GET(data []byte, headers Headers) ([]byte, error) {
	req, err := http.NewRequest("GET", c.url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		println("client do err:", err)
		return nil, err
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()

	d, _ := io.ReadAll(httpResp.Body)

	return d, nil
}

func (c *Client) POST(data []byte, headers Headers) ([]byte, error) {
	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		println("client do err:", err)
		return nil, err
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()

	d, _ := io.ReadAll(httpResp.Body)

	return d, nil
}
