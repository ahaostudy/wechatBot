package store

import (
	"github.com/eatmoreapple/openwechat"
	"time"
)

// ConverseStore user.ID()为键 time.Now().UnixMicro()为值
type ConverseStore map[string]int64

var CS ConverseStore

func init() {
	CS = ConverseStore{}
}

func (cs ConverseStore) Append(msg *openwechat.Message) (uid string) {
	var name string

	if msg.IsSendByGroup() {
		user, _ := msg.SenderInGroup()
		uid, name = "g"+user.NickName, user.NickName
	} else {
		user, _ := msg.Sender()
		uid, name = "u"+user.NickName, user.NickName
	}
	println(name, ":", msg.Content)
	cs[uid] = time.Now().UnixMicro()
	return
}
