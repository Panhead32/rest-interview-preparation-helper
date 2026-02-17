package auth

import (
	"context"
	"database/sql"
	modelAuth "interview-project/internal/models/auth"
	"interview-project/pkg/utils"
)

type AuthRepository interface {
	Create(ctx context.Context, userID int64, auth *modelAuth.AuthInput, passwordHash string) error
	GetByUserID(ctx context.Context, userID int64) (*modelAuth.Auth, error)
	GetByEmail(ctx context.Context, email string) (*modelAuth.Auth, error)
	Update(ctx context.Context, userID int64, auth *modelAuth.AuthInput, passwordHash string) error
	Delete(ctx context.Context, userID int64) error
}

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, userID int64, auth *modelAuth.AuthInput, passwordHash string) error {
	const query = `
		INSERT INTO auth (
			user_id,
			email,
			password_hash,
			created_at
		) VALUES (?, ?, ?, ?)
	`

	now := utils.NowTimestamp()
	_, err := r.db.ExecContext(
		ctx,
		query,
		userID,
		auth.Email,
		passwordHash,
		now,
	)
	return err
}

func (r *Repository) GetByUserID(ctx context.Context, userID int64) (*modelAuth.Auth, error) {
	const query = `
		SELECT user_id, email, password_hash, created_at
		FROM auth
		WHERE user_id = ?
	`

	authRec := &modelAuth.Auth{}
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&authRec.UserID,
		&authRec.Email,
		&authRec.PasswordHash,
		&authRec.CreatedAt,
	); err != nil {
		return nil, err
	}

	return authRec, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*modelAuth.Auth, error) {
	const query = `
		SELECT user_id, email, password_hash, created_at
		FROM auth
		WHERE email = ?
	`

	authRec := &modelAuth.Auth{}
	if err := r.db.QueryRowContext(ctx, query, email).Scan(
		&authRec.UserID,
		&authRec.Email,
		&authRec.PasswordHash,
		&authRec.CreatedAt,
	); err != nil {
		return nil, err
	}

	return authRec, nil
}

func (r *Repository) Update(ctx context.Context, userID int64, auth *modelAuth.AuthInput, passwordHash string) error {
	const query = `
		UPDATE auth
		SET email = ?, password_hash = ?
		WHERE user_id = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		auth.Email,
		passwordHash,
		userID,
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, userID int64) error {
	const query = `DELETE FROM auth WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
