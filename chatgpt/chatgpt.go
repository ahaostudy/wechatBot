package chatgpt

import (
	"github.com/poorjobless/wechatbot/config"
)

func ChatCompletions(messages []map[string]string) (string, error) {

	c := NewChatClient(config.Config.ApiKey)

	req := &Request{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}
	resp, err := c.GetChat(req)

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, err
	}
	println("content error...")
	return "哎呀！出问题啦，请尝试重新提问吧~", err
}

func ImageCompletions(prompt string) (string, error) {
	// 生成
	c := NewImageClient(config.Config.ApiKey)
	req := &ImageRequest{
		Prompt: prompt,
		N:      1,
		Size:   "512x512",
	}
	resp, _ := c.GetImage(req)
	if len(resp.Data) == 0 {
		println("image error...")
		return "", nil
	}
	return resp.Data[0].Url, nil

	// 获取
	//res, err := http.Get(url)
	//if err != nil {
	//	fmt.Println("get image error:", err)
	//	return nil, err
	//}
	//defer res.Body.Close()
	//reader := bufio.NewReaderSize(res.Body, 512*512)
	//image, err := os.CreateTemp("", "temp.png")
	//if err != nil {
	//	println("create temp file error:", err)
	//	return nil, err
	//}
	//defer os.Remove(image.Name())
	//_, err = io.Copy(image, reader)
	//if err != nil {
	//	println("io copy error:", err)
	//	return nil, err
	//}
	//return image, nil
}
