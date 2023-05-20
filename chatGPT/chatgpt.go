package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TextCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []TextCompletionChoice `json:"choices"`
	Usage   TextCompletionUsage    `json:"usage"`
}

type TextCompletionChoice struct {
	Text         string      `json:"text"`
	Index        int         `json:"index"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

type TextCompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func Ask(apiKey, prompt string) (TextCompletionResponse, error) {
	// Set up the request body
	requestBody, err := json.Marshal(map[string]interface{}{
		"model":       "text-davinci-003",
		"prompt":      prompt,
		"n":           1,
		"stop":        ".",
		"temperature": 0.1,
	})
	if err != nil {
		return TextCompletionResponse{}, err
	}

	// Set up the HTTP request
	url := "https://api.openai.com/v1/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return TextCompletionResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// Send the HTTP request and parse the response
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return TextCompletionResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return TextCompletionResponse{}, fmt.Errorf("non 200 returned: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return TextCompletionResponse{}, err
	}

	textCompletionResponse := TextCompletionResponse{}
	err = json.Unmarshal(body, &textCompletionResponse)
	if err != nil {
		return TextCompletionResponse{}, err
	}

	return textCompletionResponse, nil
}
