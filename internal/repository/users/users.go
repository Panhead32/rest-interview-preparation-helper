package users

import (
	"context"
	"database/sql"
	modelUsers "interview-project/internal/models/users"
	"interview-project/pkg/utils"
)

type UsersRepository interface {
	Create(ctx context.Context, user *modelUsers.UserInput) error
	GetByID(ctx context.Context, id int64) (*modelUsers.User, error)
	GetByIDFullInfo(ctx context.Context, id int64) (*modelUsers.UserFullInfo, error)
	GetByEmail(ctx context.Context, email string) (*modelUsers.User, error)
	Update(ctx context.Context, user *modelUsers.User) error
	Delete(ctx context.Context, id int64) error
}

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, user *modelUsers.UserInput) error {
	const query = `
		INSERT INTO users (
			name,
			surname,
			nickname,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?)
	`

	now := utils.NowTimestamp()
	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.Surname,
		user.NickName,
		now,
		now,
	)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*modelUsers.User, error) {
	const query = `
		SELECT id, name, surname, nickname, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	user := &modelUsers.User{}
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Surname,
		&user.NickName,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetByIDFullInfo(ctx context.Context, id int64) (*modelUsers.UserFullInfo, error) {
	const query = `
		SELECT u.id, u.name, u.surname, u.nickname, u.created_at, u.updated_at, a.email
		FROM users u
		JOIN auth a ON u.id = a.user_id
		WHERE u.id = ?
	`

	user := &modelUsers.UserFullInfo{}
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Surname,
		&user.NickName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*modelUsers.User, error) {
	const query = `
		SELECT u.id, u.name, u.surname, u.nickname, u.created_at, u.updated_at
		FROM users u
		JOIN auth a ON u.id = a.user_id
		WHERE a.email = ?
	`

	user := &modelUsers.User{}
	if err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Surname,
		&user.NickName,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) Update(ctx context.Context, user *modelUsers.User) error {
	const query = `
		UPDATE users
		SET name = ?, surname = ?, nickname = ?, updated_at = ?
		WHERE id = ?
	`

	user.UpdatedAt = utils.NowTimestamp()
	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.Surname,
		user.NickName,
		user.UpdatedAt,
		user.ID,
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM users WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
