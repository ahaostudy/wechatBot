package api

type RoleType string
type ModelType string

type WeeklyRequest struct {
	Prompt string `json:"prompt"`
}

type ChatResponse struct {
	Message    string `json:"message"`
	RawMessage string `json:"raw_message"`
	Status     string `json:"status"`
}

type GrantsResponse struct {
	TotalGranted   float64 `json:"total_granted"`
	TotalUsed      float64 `json:"total_used"`
	TotalAvailable float64 `json:"total_available"`
}

type ImageRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type ImageResponse struct {
	Created int `json:"created"`
	Data    []struct {
		Url string `json:"url"`
	} `json:"data"`
}

type Headers map[string]string
