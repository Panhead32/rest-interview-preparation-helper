package questions

import (
	"context"
	"database/sql"
	modelQuestions "interview-project/internal/models/questions"
	"interview-project/pkg/utils"
)

type QuestionsRepository interface {
	Create(ctx context.Context, question *modelQuestions.QuestionInput) error
	CreateBatch(ctx context.Context, questions *[]modelQuestions.QuestionInput) ([]modelQuestions.Question, error)
	GetByRequestID(ctx context.Context, requestID int64) ([]modelQuestions.Question, error)
	GetByID(ctx context.Context, id int64) (modelQuestions.Question, error)
	GetByUserID(ctx context.Context, userID int64) ([]modelQuestions.Question, error)
	Update(ctx context.Context, question *modelQuestions.Question) error
	Delete(ctx context.Context, id int64) error
}

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, question *modelQuestions.QuestionInput) error {
	const query = `
		INSERT INTO questions (
			request_id,
			question,
			answer,
			user_id,
			level,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	now := utils.NowTimestamp()
	_, err := r.db.ExecContext(
		ctx,
		query,
		question.RequestID,
		question.Question,
		question.Answer,
		question.UserID,
		question.Level,
		now,
		now,
	)
	return err
}

func (r *Repository) CreateBatch(ctx context.Context, questions *[]modelQuestions.QuestionInput) ([]modelQuestions.Question, error) {
	const query = `
		INSERT INTO questions (
			request_id,
			question,
			answer,
			user_id,
			level,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := utils.NowTimestamp()
	var createdQuestions []modelQuestions.Question
	for _, question := range *questions {
		result, err := tx.ExecContext(
			ctx,
			query,
			question.RequestID,
			question.Question,
			question.Answer,
			question.UserID,
			question.Level,
			now,
			now,
		)
		if err != nil {
			return nil, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		createdQuestions = append(createdQuestions, modelQuestions.Question{
			ID:        id,
			RequestID: question.RequestID,
			Question:  question.Question,
			Answer:    question.Answer,
			UserID:    question.UserID,
			Level:     question.Level,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return createdQuestions, nil
}

func (r *Repository) GetByRequestID(ctx context.Context, requestID int64) ([]modelQuestions.Question, error) {
	const query = `
		SELECT id, request_id, question, answer, user_id, level, created_at, updated_at
		FROM questions
		WHERE request_id = ?
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []modelQuestions.Question
	for rows.Next() {
		item := modelQuestions.Question{}
		if err := rows.Scan(
			&item.ID,
			&item.RequestID,
			&item.Question,
			&item.Answer,
			&item.UserID,
			&item.Level,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (modelQuestions.Question, error) {
	const query = `
		SELECT id, request_id, question, answer, user_id, level, created_at, updated_at
		FROM questions
		WHERE id = ?
	`

	item := modelQuestions.Question{}
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.RequestID,
		&item.Question,
		&item.Answer,
		&item.UserID,
		&item.Level,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return modelQuestions.Question{}, err
	}

	return item, nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID int64) ([]modelQuestions.Question, error) {
	const query = `
		SELECT id, request_id, question, answer, user_id, level, created_at, updated_at
		FROM questions
		WHERE user_id = ?
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []modelQuestions.Question
	for rows.Next() {
		item := modelQuestions.Question{}
		if err := rows.Scan(
			&item.ID,
			&item.RequestID,
			&item.Question,
			&item.Answer,
			&item.UserID,
			&item.Level,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *Repository) Update(ctx context.Context, question *modelQuestions.Question) error {
	const query = `
		UPDATE questions
		SET request_id = ?, question = ?, answer = ?, user_id = ?, level = ?, updated_at = ?
		WHERE id = ?
	`

	question.UpdatedAt = utils.NowTimestamp()
	_, err := r.db.ExecContext(
		ctx,
		query,
		question.RequestID,
		question.Question,
		question.Answer,
		question.UserID,
		question.Level,
		question.UpdatedAt,
		question.ID,
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM questions WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
