package messages

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/poorjobless/wechatbot/config"
	"time"
)

type Messages struct {
	Msgs  []map[string]string // 记录有效收发消息
	Times []int64             // 只记录用户发送消息时的时间
}

type MessageStore map[string]*Messages

var MS MessageStore

func init() {
	MS = make(MessageStore)
}

func NewMsgs() (msgs []map[string]string) {
	if len(config.Config.SystemContent) > 0 {
		msgs = append(msgs, map[string]string{
			"role":    "system",
			"content": config.Config.SystemContent,
		})
	} else {
		msgs = []map[string]string{}
	}
	return
}

func (ms MessageStore) Append(msg *openwechat.Message) (msgs *Messages) {
	var id string
	sender, _ := msg.Sender()
	if msg.IsSendByGroup() {
		user, _ := msg.SenderInGroup()
		id = "g" + sender.NickName + user.NickName
	} else {
		id = "u" + sender.NickName
	}
	// 加入到队列中
	if m, ok := ms[id]; ok {
		msgs = m
	} else {
		msgs = new(Messages)
		msgs.Msgs = NewMsgs()
	}
	msgs.Append("user", msg.Content)
	ms[id] = msgs
	// 输出日志
	println("count:", len(msgs.Times))
	return
}

// Append 添加消息和时间
func (messages *Messages) Append(role, content string) {
	if messages.RemainingTime() == 0 {
		messages.AppendMessages(role, content)
	}
	messages.Times = append(messages.Times, time.Now().UnixMicro())
}

// AppendMessages 只添加消息
func (messages *Messages) AppendMessages(role, content string) {
	messages.Msgs = append(messages.Msgs, map[string]string{
		"role":    role,
		"content": content,
	})
}

// RemainingTime 返回下次请求需要的剩余用时，不需要时返回0
func (messages *Messages) RemainingTime() int {
	if len(messages.Times) > config.Config.Count {
		diff := time.Now().UnixMicro() - messages.Times[len(messages.Times)-config.Config.Count]
		if diff < config.Config.Duration*60000000 {
			return int(config.Config.Duration*60000000-diff) / 1000000
		}
	}
	return 0
}

// Count O(n)
func (messages *Messages) Count() (count int) {
	now := time.Now().UnixMicro()
	for i := len(messages.Times) - 1; i >= 0; i-- {
		if now-messages.Times[i] < config.Config.Duration*60000000 {
			count++
		} else {
			return count
		}
	}
	return count
}

func (messages *Messages) Clear() {
	messages.Msgs = NewMsgs()
}
