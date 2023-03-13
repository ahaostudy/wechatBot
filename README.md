## 微信机器人（使用Openai最新接口: gpt-3.5-turbo & images 模型）
基于openWechat与openai开发的微信机器人，支持自动回复，生成图片、生成周报，自动同意好友请求，自定义访问频率等功能。

### 目前已实现以下功能
+ 群聊@或/c进行聊天
+ 支持连续对话(自动清理会10分钟无更新的对话)
+ /img指令生成图片
+ /weekly指令生成周报
+ 支持限制请求频率
+ 私聊自动回复
+ 自动通过好友申请
+ 机器人退出时邮箱提醒
+ //set线上设置聊天频率

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
  // 机器人系统设置
  "system_content": "你是微信机器人",
  // api_key列表，每次调用时随机使用一个key
  "api_key": [  
    "sk-xxx"
  ],
  // 收件邮箱
  "emails": [  
    "xxx@email.com"
  ],
  // 管理员ID，请在开启机器人后向其发送`//id`命令获取您的ID（9位数）
  // 如何你和机器人使用的是同一个账号，可以不用获取，机器人默认为管理员
  "admin_id": "123456789",
  // 白名单群聊
  // 机器人对于群聊的访问频率有限制，在次添加白名单可以将群设置为无限制
  "whitelist": ["WechatRobot测试群"],
  // 是否开启机器人退出时邮箱提醒
  "is_warning": true,
  // 是否自动同意添加好友
  "auto_pass": true,
  // 限制频率，duration分钟内只可访问count次
  "duration": 5,
  "count": 8
}

```

### 启动项目
`go run main.go`

## 鸣谢
+ @poorjobless
+ @djun
+ @eatmoreapple
