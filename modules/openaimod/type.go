package openaimod

type ChatGPTResp struct {
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   struct {
		TotalTokens      int64 `json:"total_tokens"`
		PromptTokens     int64 `json:"prompt_tokens"`
		CompletionTokens int64 `json:"completion_tokens"`
	} `json:"usage"`
}

type Choice struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	FinishReason string `json:"finish_reason"`
	Index        int    `json:"index"`
}
