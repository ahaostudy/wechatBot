package handlers

import (
	"fmt"
	"github.com/poorjobless/wechatbot/api"
	"github.com/poorjobless/wechatbot/messages"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/eatmoreapple/openwechat"
	"github.com/poorjobless/wechatbot/config"
)

// Handler 全局处理入口
func Handler(msg *openwechat.Message) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 好友申请
		if msg.IsFriendAdd() {
			if config.Config.AutoPass {
				_, err := msg.Agree("你好我是基于ChatGPT开发的WeChat Robot，你可以向我提问任何问题。")
				if err != nil {
					log.Fatalf("add friend agree error : %v", err)
					return
				}
			}
		}
		// 非文本消息不处理
		if !msg.IsText() {
			return
		}
		// 处理文本 去除@
		sender, _ := msg.Sender()
		msg.Content = strings.TrimSpace(strings.ReplaceAll(msg.Content, "@"+sender.Self.NickName, ""))
		// 处理业务
		handle(msg)(msg)
	}()
	wg.Wait()
}

// 获取当前用户
func getUser(msg *openwechat.Message) *openwechat.User {
	var user *openwechat.User
	if msg.IsSendByGroup() {
		user, _ = msg.SenderInGroup()
	} else {
		user, _ = msg.Sender()
	}
	println("userid", user.ID())
	return user
}

// 命令处理
func command(msg *openwechat.Message) (flag bool) {
	if msg.IsSendByGroup() && !msg.IsAt() {
		return false
	}
	var reply string

	// 检测
	if flag = isCommand(msg, "清空", "清空记录"); flag {
		messages.MS.Append(msg).Clear()
		reply = "会话记录已经清空啦！"
	} else if flag = isCommand(msg, "余额"); flag {
		reply, _ = api.Grants()
	}

	// 响应
	if flag {
		replyText(msg, reply)
	}
	return
}

// 消息处理
func handle(msg *openwechat.Message) func(message *openwechat.Message) {
	println(getUser(msg).NickName+":", msg.Content)
	if command(msg) {
		return func(message *openwechat.Message) {}
	}
	if isPrefix(msg, "//set") {
		return setConfig
	}
	if isPrefix(msg, "//id") {
		return getId
	}
	if isPrefix(msg, "/weekly") {
		return weekly
	}
	if isPrefix(msg, "/img") {
		return image
	}
	if isPrefix(msg, "/c", "/chat") || !msg.IsSendByGroup() || msg.IsAt() {
		return chat
	}
	return func(message *openwechat.Message) {}
}

// 检测msg.Content是否以pre为前缀
func isPrefix(msg *openwechat.Message, pres ...string) bool {
	for _, pre := range pres {
		if len(msg.Content) >= len(pre) && msg.Content[:len(pre)] == pre {
			println(pre, "request")
			msg.Content = strings.TrimSpace(msg.Content[len(pre):])
			return true
		}
	}
	return false
}

// 检测msg.Content是否等于str
func isCommand(msg *openwechat.Message, strs ...string) bool {
	for _, str := range strs {
		if msg.Content == str {
			return true
		}
	}
	return false
}

// 判断是否为白名单
func isWhiteList(msg *openwechat.Message) bool {
	if !msg.IsSendByGroup() {
		return true
	}
	sender, _ := msg.Sender()
	for _, w := range config.Config.Whitelist {
		if sender.NickName == w {
			return true
		}
	}
	return false
}

// 聊天功能
func chat(msg *openwechat.Message) {
	// 添加消息记录
	msgs := messages.MS.Append(msg)
	// 异步锁
	wg := new(sync.WaitGroup)

	// 1. 清除过期记录
	wg.Add(1)
	go messages.Check(wg)

	// 2. 处理消息
	var reply string
	// 私聊或群聊请求限制内
	if isWhiteList(msg) || msgs.RemainingTime() == 0 {
		// 发送请求
		var err error
		reply, err = api.Chat(msgs.Msgs)
		if err != nil {
			println("api chat request error:", err)
		} else {
			msgs.AppendMessages("assistant", reply)
		}
	} else {
		if msgs.Count() > config.Config.Count+1 {
			return
		}
		// 第一次请求频繁时提示
		reply = "你的访问频率太高了，请在 " + strconv.Itoa(msgs.RemainingTime()) + "s 后在提问吧。"
	}
	// 返回响应
	replyText(msg, reply)
	// 阻塞等待检测完毕
	wg.Wait()
}

// 周报功能
func weekly(msg *openwechat.Message) {
	reply, err := api.Weekly(msg.Content)
	if err != nil {
		println("api weekly request error:", err)
		return
	}
	replyText(msg, reply)
}

// 图片功能
func image(msg *openwechat.Message) {
	url, err := api.Image(msg.Content)
	if err != nil {
		println("api image request error:", err)
		return
	}
	replyImage(msg, url)
}

// 设置配置
func setConfig(msg *openwechat.Message) {
	if getUser(msg).ID() != config.Config.AdminId {
		return
	}
	params := strings.Split(msg.Content, ":")
	if len(params) != 2 {
		println("set request: params error")
		return
	}
	d, _ := strconv.ParseInt(params[0], 10, 64)
	c, _ := strconv.Atoi(params[1])
	if d == 0 || c == 0 {
		reply := "不支持数据为0"
		println("ChatGPT:", reply)
		replyText(msg, reply)
		return
	}
	// 设置限制
	config.Config.Duration = d
	config.Config.Count = c
	reply := fmt.Sprintf("已设置限制为 %d 分钟最多访问 %d 次。", config.Config.Duration, config.Config.Count)
	println("ChatGPT:", reply)
	replyText(msg, reply)
}

// 获取用户ID（只回复私聊获取）
func getId(msg *openwechat.Message) {
	if !msg.IsSendByGroup() {
		replyText(msg, getUser(msg).ID())
	}
}

// replyText 发送消息
func replyText(msg *openwechat.Message, reply string) {
	// 消息为空
	if len(reply) == 0 {
		log.Printf("response text length ==  0")
		reply = "哎呀！出问题啦，请尝试重新提问吧~"
	}

	// 消息修剪
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")

	println("ChatGPT:", reply)

	// 群聊回复时加上@
	if msg.IsSendByGroup() {
		groupSender, err := msg.SenderInGroup()
		if err != nil {
			log.Printf("get sender in group error :%v \n", err)
		}
		atText := "@" + groupSender.NickName
		reply = atText + " " + reply
	}

	// 发送
	_, err := msg.ReplyText(reply)
	if err != nil {
		log.Printf("reply text error: %v \n", err)
	}
}

// replyImage 发送图片
func replyImage(msg *openwechat.Message, url string) {
	resp, err := http.Get(url)
	if err != nil {
		println("http get image error:", err.Error())
		replyText(msg, "哎呀，图片生成失败啦。")
		return
	}
	println("生成图片:", url)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	// 临时保存
	tempFileName := path.Base(url) + ".jpg"
	data, _ := io.ReadAll(resp.Body)
	err = os.WriteFile(tempFileName, data, 0644)
	if err != nil {
		replyText(msg, msg.Content+"已经生成好啦，请点击进行查看吧~\n"+url)
		return
	}
	file, _ := os.Open(tempFileName)
	defer func(file *os.File) {
		_ = file.Close()
		_ = os.Remove(tempFileName)
	}(file)
	_, err = msg.ReplyImage(file)
	if err != nil {
		replyText(msg, msg.Content+"已经生成好啦，请点击进行查看吧~\n"+url)
	}
}
