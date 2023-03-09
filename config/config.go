package config

import (
	"encoding/json"
	"log"
	"os"
)

// Configuration 项目配置
type Configuration struct {
	// 代理API
	Api string `json:"api"`
	// gtp apikey
	ApiKey string `json:"api_key"`
	// 默认提示信息
	SystemContent string `json:"system_content"`
	// 自动通过好友
	AutoPass bool `json:"auto_pass"`
	// 图片模式
	Image bool `json:"image"`
	// 时间（分钟）
	Duration int64 `json:"duration"`
	// Duration 时间内限制的请求数
	Count int `json:"count"`
}

var Config *Configuration

func init() {
	// 从文件中读取
	Config = &Configuration{}
	f, err := os.Open("Config.json")
	if err != nil {
		log.Fatalf("open Config err: %v", err)
		return
	}
	defer f.Close()
	encoder := json.NewDecoder(f)
	err = encoder.Decode(Config)
	if err != nil {
		log.Fatalf("decode Config err: %v", err)
		return
	}
}
