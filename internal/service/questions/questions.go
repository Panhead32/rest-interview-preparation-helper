package questions

import (
	"context"
	"encoding/json"
	"interview-project/internal/models/questions"
	"interview-project/internal/models/requests"
	"interview-project/internal/repository"
	"interview-project/pkg/chromedp"
	"interview-project/pkg/openai"
	"interview-project/pkg/utils"
)

type QuestionService interface {
	ParseQuestions(ctx context.Context, link string) ([]questions.Question, error)
	GetQuestions(ctx context.Context, userID int64) ([]questions.Question, error)
	GetQuestionByID(ctx context.Context, questionID int64) (questions.Question, error)
	ExplainQuestion(ctx context.Context, questionID int64) error
}

type service struct {
	repo *repository.Manager
}

func New(repo *repository.Manager) *service {
	return &service{repo: repo}
}

func (s *service) ParseQuestions(ctx context.Context, link string) ([]questions.Question, error) {
	text, err := chromedp.ScrapeArticleText(link)
	if err != nil {
		return nil, err
	}

	request, err := s.repo.Requests.Create(ctx, &requests.RequestInput{
		RequestText: text,
		UserID:      ctx.Value(utils.UserIDKey).(int64),
	})

	if err != nil {
		return nil, err
	}

	openaiClient := &openai.Client{}

	prompt := "Generate json file from the article text with list of questions that can be asked in an interview based on the article text. The json file should have the following format: {\"questions\": [{\"text\": \"question text\", \"topic\": \"topic name\", \"level\": \"junior|middle|senior\"}, ...]}. Article text: " + text

	err, responseStr := openaiClient.GetResponse(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var response questions.QuestionsResponse
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return nil, err
	}

	processedQuestions := make([]questions.QuestionInput, 0, len(response.Questions))
	for _, question := range response.Questions {
		processedQuestions = append(processedQuestions, questions.QuestionInput{
			RequestID: request.ID,
			Question:  question.Text,
			UserID:    ctx.Value(utils.UserIDKey).(int64),
			Level:     question.Level,
		})
	}

	insertedQuestions, err := s.repo.Questions.CreateBatch(ctx, &processedQuestions)
	if err != nil {
		return nil, err
	}

	return insertedQuestions, nil
}

func (s *service) GetQuestions(ctx context.Context, userID int64) ([]questions.Question, error) {
	questionsList, err := s.repo.Questions.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return questionsList, nil
}

func (s *service) GetQuestionByID(ctx context.Context, questionID int64) (questions.Question, error) {
	question, err := s.repo.Questions.GetByID(ctx, questionID)
	if err != nil {
		return questions.Question{}, err
	}

	return question, nil
}

func (s *service) ExplainQuestion(ctx context.Context, questionID int64) error {
	question, err := s.repo.Questions.GetByID(ctx, questionID)
	if err != nil {
		return err
	}

	openaiClient := &openai.Client{}

	prompt := "Provide a detailed explanation for the following interview question: " + question.Question

	err, responseStr := openaiClient.GetResponse(ctx, prompt)
	if err != nil {
		return err
	}

	question.Answer = responseStr
	question.UpdatedAt = utils.NowTimestamp()

	err = s.repo.Questions.Update(ctx, &question)
	if err != nil {
		return err
	}

	return nil
}
