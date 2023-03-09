package chatgpt

import (
	"bytes"
	"encoding/json"
	"github.com/poorjobless/wechatbot/config"
	"net/http"
)

var (
	ChatUrl  string
	ImageUrl string
)

func init() {
	ChatUrl = config.Config.Api + "/v1/chat/completions"
	ImageUrl = config.Config.Api + "/v1/images/generations"
}

type Client struct {
	transport *http.Client
	apiKey    string
	url       string
}

func NewChatClient(apiKey string) *Client {
	return &Client{
		transport: http.DefaultClient,
		apiKey:    apiKey,
		url:       ChatUrl,
	}
}

func NewImageClient(apiKey string) *Client {
	return &Client{
		transport: http.DefaultClient,
		apiKey:    apiKey,
		url:       ImageUrl,
	}
}

func (c *Client) GetChat(r *Request) (*Response, error) {
	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	client := &http.Client{}
	httpResp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()

	var resp Response
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) GetImage(r *ImageRequest) (*ImageResponse, error) {
	jsonData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	client := &http.Client{}
	httpResp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()

	var resp ImageResponse
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
