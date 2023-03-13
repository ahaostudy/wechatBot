package api

import (
	"encoding/json"
	"fmt"
	"github.com/poorjobless/wechatbot/config"
	"net/url"
)

func Chat(messages []map[string]string) (string, error) {
	c := NewClient(ChatUrl)
	// form-data 类型数据
	data := url.Values{}
	context, _ := json.Marshal(messages)
	data.Set("context", string(context))
	data.Set("key", config.GetAPIKey())
	data.Set("id", "1")
	headers := Headers{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	// 获取结果并处理
	resByte, err := c.POST([]byte(data.Encode()), headers)
	var res ChatResponse
	err = json.Unmarshal(resByte, &res)
	if err != nil {
		return "", err
	}
	return res.Message, err
}

func Weekly(prompt string) (string, error) {
	c := NewClient(WeeklyUrl)
	req := &WeeklyRequest{
		Prompt: prompt,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	resByte, err := c.POST(data, nil)
	return string(resByte), err
}

func Grants() (string, error) {
	c := NewClient(GrantsUrl)
	var totalGranted, totalUsed, totalAvailable float64
	var res GrantsResponse
	for _, key := range config.Config.ApiKey {
		headers := Headers{
			"Authorization": "Bearer " + key,
		}
		resByte, _ := c.GET(nil, headers)
		if json.Unmarshal(resByte, &res) == nil {
			totalGranted += res.TotalGranted
			totalUsed += res.TotalUsed
			totalAvailable += res.TotalAvailable
		}
	}
	reply := fmt.Sprintf("余额查询\n全部配额: %.4f$\n已使用配额: %.4f$\n剩余配额: %.4f$", totalGranted, totalUsed, totalAvailable)
	return reply, nil
}

func Image(prompt string) (string, error) {
	c := NewClient(ImageUrl)
	req := &ImageRequest{
		Model:  "image-alpha-001",
		Prompt: prompt,
		N:      1,
		Size:   "256x256",
	}
	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	headers := Headers{
		"Authorization": "Bearer " + config.GetAPIKey(),
		"Content-Type":  "application/json",
	}
	resByte, err := c.POST(data, headers)
	var res ImageResponse
	err = json.Unmarshal(resByte, &res)
	if err != nil {
		return "", err
	}
	return res.Data[0].Url, err
}
