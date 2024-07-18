package llm

import (
	"encoding/json"
	"fmt"
)

func init() {
	RegisterProvider("groq", NewGroqProvider)
}

// GroqProvider implements the Provider interface for Groq
type GroqProvider struct {
	apiKey string
	model  string
}

func NewGroqProvider(apiKey, model string) Provider {
	return &GroqProvider{
		apiKey: apiKey,
		model:  model,
	}
}

func (p *GroqProvider) Name() string {
	return "groq"
}

func (p *GroqProvider) Endpoint() string {
	return "https://api.groq.com/openai/v1/chat/completions"
}

func (p *GroqProvider) Headers() map[string]string {
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + p.apiKey,
	}
}

func (p *GroqProvider) PrepareRequest(prompt string, options map[string]interface{}) ([]byte, error) {
	requestBody := map[string]interface{}{
		"model": p.model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	for k, v := range options {
		requestBody[k] = v
	}

	return json.Marshal(requestBody)
}

func (p *GroqProvider) ParseResponse(body []byte) (string, error) {
	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	err := json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}

	if len(response.Choices) == 0 || response.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("empty response from API")
	}

	return response.Choices[0].Message.Content, nil
}

