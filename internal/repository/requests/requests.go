package requests

import (
	"context"
	"database/sql"
	modelRequests "interview-project/internal/models/requests"
	"interview-project/pkg/utils"
)

type RequestsRepository interface {
	Create(ctx context.Context, request *modelRequests.RequestInput) (*modelRequests.Request, error)
	GetByID(ctx context.Context, id int64) (*modelRequests.Request, error)
	GetByUserID(ctx context.Context, userID int64) ([]modelRequests.Request, error)
	Update(ctx context.Context, request *modelRequests.Request) error
	Delete(ctx context.Context, id int64) error
}

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, request *modelRequests.RequestInput) (*modelRequests.Request, error) {
	const query = `
		INSERT INTO requests (
			request_text,
			user_id,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?)
	`

	now := utils.NowTimestamp()
	result, err := r.db.ExecContext(
		ctx,
		query,
		request.RequestText,
		request.UserID,
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

	return &modelRequests.Request{
		ID:          id,
		RequestText: request.RequestText,
		UserID:      request.UserID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*modelRequests.Request, error) {
	const query = `
		SELECT id, request_text, user_id, created_at, updated_at
		FROM requests
		WHERE id = ?
	`

	request := &modelRequests.Request{}
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&request.ID,
		&request.RequestText,
		&request.UserID,
		&request.CreatedAt,
		&request.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return request, nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID int64) ([]modelRequests.Request, error) {
	const query = `
		SELECT id, request_text, user_id, created_at, updated_at
		FROM requests
		WHERE user_id = ?
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []modelRequests.Request
	for rows.Next() {
		request := modelRequests.Request{}
		if err := rows.Scan(
			&request.ID,
			&request.RequestText,
			&request.UserID,
			&request.CreatedAt,
			&request.UpdatedAt,
		); err != nil {
			return nil, err
		}
		results = append(results, request)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *Repository) Update(ctx context.Context, request *modelRequests.Request) error {
	const query = `
		UPDATE requests
		SET request_text = ?, updated_at = ?
		WHERE id = ?
	`

	request.UpdatedAt = utils.NowTimestamp()
	_, err := r.db.ExecContext(
		ctx,
		query,
		request.RequestText,
		request.UpdatedAt,
		request.ID,
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM requests WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
