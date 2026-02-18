package openai

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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

func (c *Client) GetResponse(ctx context.Context, prompt string) (error, string) {

	hashedPrompt := sha256.New()
	hashedPrompt.Write([]byte(prompt))
	hash := hex.EncodeToString(hashedPrompt.Sum(nil))

	resp, err := c.getClient().Responses.New(ctx, responses.ResponseNewParams{
		Model:          openai.ChatModelGPT5Nano,
		Input:          responses.ResponseNewParamsInputUnion{OfString: openai.String(prompt)},
		PromptCacheKey: openai.String(hash),
	})

	if err != nil {
		fmt.Printf("Error creating response: %v\n", err)
		return err, ""
	}

	return nil, resp.OutputText()
}
