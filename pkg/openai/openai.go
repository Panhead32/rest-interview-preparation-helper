package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"interview-project/internal/config"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type Client struct {
	client *openai.Client
}

func (c *Client) getClient() *openai.Client {
	if c.client == nil {
		c.client = &openai.Client{}

		cfg := config.LoadConfig()
		apiKey := cfg.OpenAI.APIKey

		*c.client = openai.NewClient(
			option.WithAPIKey(apiKey),
		)
	}
	return c.client
}

func (c *Client) GetResponse(prompt string) (error, []string) {

	resp, err := c.getClient().Responses.New(context.TODO(), responses.ResponseNewParams{
		Model: openai.ChatModelGPT5Nano,
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(prompt)},
		// PromptCacheKey: openai.String("test-prompt-cache-key"),
	})
	if err != nil {
		fmt.Printf("Error creating response: %v\n", err)
		return err, nil
	}

	var questions []string
	if err := json.Unmarshal([]byte(resp.OutputText()), &questions); err != nil {
		fmt.Printf("Error parsing JSON response: %v\n", err)
		return err, nil
	}
	return nil, questions
}
