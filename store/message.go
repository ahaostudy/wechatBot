package store

import (
	"fmt"
	"github.com/poorjobless/wechatbot/config"
	"time"
)

type Messages struct {
	Msgs  []map[string]string // 记录有效收发消息
	Times []int64             // 只记录发消息时的时间
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

func (ms MessageStore) Append(id, role, content string) (msgs *Messages) {
	if m, ok := ms[id]; ok {
		msgs = m
	} else {
		msgs = new(Messages)
		msgs.Msgs = NewMsgs()
	}
	msgs.Append(role, content)
	ms[id] = msgs
	fmt.Printf("%v\n", msgs.Times)
	return
}

// Append 添加消息和时间
func (messages *Messages) Append(role, content string) {
	if messages.Check() {
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

// Check O(1)
func (messages *Messages) Check() bool {
	if len(messages.Times) > config.Config.Count &&
		time.Now().UnixMicro()-messages.Times[len(messages.Times)-config.Config.Count] < config.Config.Duration*60000000 {
		return false
	}
	return true
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
