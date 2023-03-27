package main

import "encoding/json"

type ChatGPT3Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewChatGPT3Message(content string) ChatGPT3Message {
	return ChatGPT3Message{
		Role:    "user",
		Content: content,
	}
}

type RequestChatGPT3 struct {
	Model       string            `json:"model"`
	Messages    []ChatGPT3Message `json:"messages"`
	Temperature float64           `json:"temperature"`
}

func NewSimpleChatGPTRequest(systemPrompt string, contents ...string) ([]byte, error) {
	messages := make([]ChatGPT3Message, 0)
	for _, content := range contents {
		messages = append(messages, NewChatGPT3Message(content))
	}

	systemMessage := ChatGPT3Message{
		Role:    "system",
		Content: systemPrompt,
	}
	messages = append([]ChatGPT3Message{systemMessage}, messages...)

	req := &RequestChatGPT3{
		Model:       "gpt-3.5-turbo",
		Messages:    messages,
		Temperature: 0.7,
	}
	return json.Marshal(req)
}

type ResponseChatGPT3 struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}

func ContentsFromResponseChatGPT3(res []byte) ([]string, error) {
	resp := &ResponseChatGPT3{}
	if err := json.Unmarshal(res, resp); err != nil {
		return nil, err
	}
	contents := make([]string, 0)
	for _, choice := range resp.Choices {
		contents = append(contents, choice.Message.Content)
	}
	return contents, nil
}
