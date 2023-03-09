package handlers

import (
	"github.com/poorjobless/wechatbot/chatgpt"
	"github.com/poorjobless/wechatbot/store"
	"log"
	"strings"
	"sync"

	"github.com/eatmoreapple/openwechat"
	"github.com/poorjobless/wechatbot/config"
)

// Handler 全局处理入口
func Handler(msg *openwechat.Message) {
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

	if ns, f := IsPrefix(msg.Content, ".img"); config.Config.Image && f {
		// 图片请求
		msg.Content = ns
		var name string
		if msg.IsSendByGroup() {
			user, _ := msg.SenderInGroup()
			name = user.NickName
		} else {
			user, _ := msg.Sender()
			name = user.NickName
		}
		println(name + ": 请求图片\"" + msg.Content + "\"")
		ReplyImage(msg)
	} else if ns, f := IsPrefix(msg.Content, ".c"); f || !msg.IsSendByGroup() ||
		(msg.IsSendByGroup() && msg.IsAt()) {
		// 聊天请求
		msg.Content = ns
		wg := new(sync.WaitGroup)
		wg.Add(1)
		uid := store.CS.Append(msg)
		msgs := store.MS.Append(uid, "user", msg.Content)
		go store.Check(wg)
		// 判断是否为特定指令，优先处理并阻止请求
		if !IsCommand(msg, msgs) {
			ReplyText(msg, msgs)
		}
		wg.Wait()
	}

}

func IsPrefix(str, pre string) (string, bool) {
	if len(str) > len(pre) && str[:len(pre)] == pre {
		return str[len(pre):], true
	}
	return str, false
}

func ReplyImage(msg *openwechat.Message) {
	url, err := chatgpt.ImageCompletions(msg.Content)
	if err != nil {
		println("image completions error:", err)
	}
	//defer func() {
	//	err = image.Close()
	//	if err != nil {
	//		println("image close error:", err.Error())
	//		return
	//	}
	//}()
	//_, err = msg.ReplyImage(image)
	//if err != nil {
	//	println("reply image error", err)
	//	return
	//}
	Reply(msg, msg.Content+"已经生成完毕啦，请点击查看哦~\n"+url)
}

// ReplyText 发送文本消息到群
func ReplyText(msg *openwechat.Message, msgs *store.Messages) {
	var reply string
	var err error
	if !msg.IsSendByGroup() || msgs.Check() {
		// 发送请求
		reply, err = chatgpt.ChatCompletions(msgs.Msgs)
		msgs.AppendMessages("assistant", reply)
	} else {
		if msgs.Count() > config.Config.Count+1 {
			return
		}
		reply = "请求太频繁了！"
	}
	if err != nil {
		println("chatgpt request error:", err)
	}

	Reply(msg, reply)
}

// IsCommand 判断是否为特定指令，优先处理
func IsCommand(msg *openwechat.Message, msgs *store.Messages) (flag bool) {
	flag = true
	var reply string

	if msg.Content == "清空" || msg.Content == "清空记录" {
		msgs.Clear()
		reply = "会话记录已经清空啦！"
	} else {
		flag = false
	}

	if flag {
		Reply(msg, reply)
	}
	return
}

// Reply 发送消息
func Reply(msg *openwechat.Message, reply string) {
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
