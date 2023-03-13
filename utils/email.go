package utils

import (
	"fmt"
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
)

func Send(toEmails ...string) {
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = "阿浩 <1993584108@qq.com>"
	// 设置接收方的邮箱
	e.To = toEmails
	//设置主题
	e.Subject = "WechatRobot退出提醒"
	//设置文件发送的内容
	e.HTML = []byte(`<h1>WechatRobot退出提醒</h1>
    <p style="font-size: 20px;">检测到您的微信机器人已经结束运行，若不是您正常关闭，请注意检测原因哦。</p>`)
	//设置服务器相关的配置
	err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", "1993584108@qq.com", "bmxqzyaxrmjceiad", "smtp.qq.com"))
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("邮件已送达")
	}
}

func main() {
	Send("ahao_study@163.com")
}
