package chatgpt

type RoleType string
type ModelType string

type Request struct {
	Model    ModelType           `json:"model"`
	Messages []map[string]string `json:"messages"`
}

type Response struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int64     `json:"created"`
	Choices []*Choice `json:"choices"`
	Usage   *Usage    `json:"usage"`
	Error   *Error    `json:"error,omitempty"`
}

type Message struct {
	Role    RoleType `json:"role,omitempty"`
	Content string   `json:"content"`
}

type Choice struct {
	Index        int      `json:"index"`
	Message      *Message `json:"message"`
	FinishReason string   `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ImageRequest struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type ImageResponse struct {
	Created int         `json:"created"`
	Data    []ImageData `json:"data"`
}

type ImageData struct {
	Url string `json:"url"`
}

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}
