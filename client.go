// llm_client.go
package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type LLMClient interface {
	StreamCompletion(model, prompt string) error
}

type OpenAIClient struct {
	apiKey string
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{apiKey: apiKey}
}

func (c *OpenAIClient) StreamCompletion(model, prompt string) error {
	url := "https://api.openai.com/v1/chat/completions"
	payload := fmt.Sprintf(`{
		"model": "%s",
		"messages": [{"role": "user", "content": "%s"}],
		"stream": true
	}`, model, prompt)

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data != "[DONE]" {
				fmt.Printf(strings.ReplaceAll(extractContent(data), "\\n", "\n"))
			}
		}
	}

	return nil
}

func extractContent(data string) string {
	// This is a simplified extraction. In a real implementation,
	// you would want to properly parse the JSON.
	if strings.Contains(data, "content") {
		parts := strings.Split(data, "content")
		if len(parts) > 1 {
			content := strings.Split(parts[1], "\"")[2]
			return content
		}
	}
	return ""
}

type AnthropicClient struct {
	apiKey string
}

func NewAnthropicClient(apiKey string) *AnthropicClient {
	return &AnthropicClient{apiKey: apiKey}
}

func (c *AnthropicClient) StreamCompletion(model, prompt string) error {
	url := "https://api.anthropic.com/v1/completions"
	payload := fmt.Sprintf(`{
		"model": "%s",
		"prompt": "\n\nHuman: %s\n\nAssistant:",
		"max_tokens_to_sample": 300,
		"stream": true
	}`, model, prompt)

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", c.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data != "[DONE]" {
				fmt.Printf(strings.ReplaceAll(extractAnthropicContent(data), "\\n", "\n"))
			}
		}
	}

	return nil
}

func extractAnthropicContent(data string) string {
	// This is a simplified extraction. In a real implementation,
	// you would want to properly parse the JSON.
	if strings.Contains(data, "completion") {
		parts := strings.Split(data, "completion")
		if len(parts) > 1 {
			content := strings.Split(parts[1], "\"")[2]
			return content
		}
	}
	return ""
}
