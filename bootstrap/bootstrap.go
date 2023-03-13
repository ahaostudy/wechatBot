package bootstrap

import (
	"log"

	"github.com/eatmoreapple/openwechat"
	"github.com/poorjobless/wechatbot/handlers"
)

var BotID string

func Run() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式

	// 注册消息处理函数
	bot.MessageHandler = handlers.Handler
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")
	// 执行热登录
	err := bot.HotLogin(reloadStorage)
	if err != nil {
		if err = bot.Login(); err != nil {
			log.Printf("bot.login error: %v \n", err)
			return
		}
	}
	// 设置机器人账号ID
	self, err := bot.GetCurrentUser()
	BotID = self.ID()
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	err = bot.Block()
	if err != nil {
		log.Printf("bot.block error: %v \n", err)
		return
	}
}
