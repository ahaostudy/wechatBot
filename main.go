package main

import (
	"github.com/poorjobless/wechatbot/bootstrap"
	"github.com/poorjobless/wechatbot/config"
	"github.com/poorjobless/wechatbot/utils"
)

func main() {
	if config.Config.IsWarning {
		defer utils.Send(config.Config.Emails...)
	}
	bootstrap.Run()
}
