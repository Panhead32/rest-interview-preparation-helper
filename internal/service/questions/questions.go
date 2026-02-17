package questions

import (
	"context"
	"interview-project/internal/repository"
	"interview-project/pkg/openai"
)

type QuestionService interface {
	GetQuestions(ctx context.Context) ([]string, error)
}

type service struct {
	repo *repository.Manager
}

func New(repo *repository.Manager) *service {
	return &service{repo: repo}
}

func (s *service) GetQuestions(ctx context.Context) ([]string, error) {
	openaiClient := &openai.Client{}

	// TODO: Make prompt configurable
	prompt := "Generate 10 interview questions for a software engineer position"

	err, questions := openaiClient.GetResponse(prompt)
	if err != nil {
		return nil, err
	}
	return questions, nil
}
