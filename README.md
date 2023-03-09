## 微信机器人（使用Openai最新接口: gpt-3.5-turbo & images 模型）
本项目修改自 @poorjobless 的 [wechatbot](https://github.com/djun/wechatbot)

### 目前实现了以下功能
+ 群聊@或.c进行聊天
+ 支持连续对话(自动清理会10分钟无更新的对话)
+ .img指令获取图片(回复时使用链接)
+ 支持群友的聊天请求频率设置
+ 私聊自动回复
+ 自动通过好友申请

## 安装使用

### 获取项目
`git clone https://github.com/ahaostudy/wechatBot.git`

### 进入项目目录
`cd wechatBot`

### 复制配置文件
`cp config.dev.json config.json`

### 修改config.json中的内容
```json
{
  "api": "反向代理IP，国外服务器可以用官方API的IP",
  "api_key": "api_key",
  "system_content": "给机器人提供初始信息",
  "auto_pass": "是否自动添加好友",
  "image": "是否支持图片",
  "duration": "时间",
  "count": "在duration时间内最大请求数"
}
```

### 启动项目
`go run main.go`

# 鸣谢
+ @poorjobless
+ @djun
+ @eatmoreapple
