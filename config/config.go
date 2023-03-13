package config

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"
)

// Configuration 项目配置
type Configuration struct {
	// 默认提示信息
	SystemContent string `json:"system_content"`
	// Api_Key
	ApiKey []string `json:"api_key"`
	// 邮件提醒
	Emails    []string `json:"emails"`
	IsWarning bool     `json:"is_warning"`
	// 管理员名称
	AdminId string `json:"admin_id"`
	// 白名单
	Whitelist []string `json:"whitelist"`
	// 自动通过好友
	AutoPass bool `json:"auto_pass"`
	// 时间（分钟）
	Duration int64 `json:"duration"`
	// Duration 时间内限制的请求数
	Count int `json:"count"`
}

var Config *Configuration

func init() {
	// 从文件中读取
	Config = &Configuration{}
	f, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("open config err: %v", err)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("file.close error: %v\n", err)
		}
	}(f)
	encoder := json.NewDecoder(f)
	err = encoder.Decode(Config)
	if err != nil {
		log.Fatalf("decode Config err: %v", err)
		return
	}
}

func GetAPIKey() string {
	rand.Seed(time.Now().UnixMicro())
	return Config.ApiKey[rand.Intn(len(Config.ApiKey))]
}
