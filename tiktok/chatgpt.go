package tiktok

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rashad-j/tiktokuploader/configs"
)

type payload struct {
	Query string `json:"query"`
}

type response struct {
	ConversationID string `json:"conversationId"`
	Response       string `json:"response"`
}

func AskGPT(cfg configs.Config, query string) (string, error) {
	url := "https://chatgpt-ai-chat-bot.p.rapidapi.com/ask"
	payload := payload{Query: query}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		err = fmt.Errorf("error marshaling payload: %w", err)
		return "", err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		err = fmt.Errorf("error creating request: %w", err)
		return "", err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-RapidAPI-Key", "596fd604f9msh9205a55cc027e9ap1bdb1ejsn79730d1f6b79")
	req.Header.Add("X-RapidAPI-Host", "chatgpt-ai-chat-bot.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("error sending request: %w", err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("error reading response body: %w", err)
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("request failed with status code: %d", res.StatusCode)
		return "", err
	}

	var response response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(string(body), response)
		err = fmt.Errorf("error unmarshaling response: %w", err)
		return "", err
	}

	return response.Response, nil
}
