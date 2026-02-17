package users

import (
	"context"
	modelUsers "interview-project/internal/models/users"
	"interview-project/internal/repository"
)

type UserService interface {
	Create(ctx context.Context, user *modelUsers.UserInput) error
	GetByID(ctx context.Context, id int64) (*modelUsers.User, error)
	GetByIDFullInfo(ctx context.Context, id int64) (*modelUsers.UserFullInfo, error)
	GetByEmail(ctx context.Context, email string) (*modelUsers.User, error)
	Update(ctx context.Context, user *modelUsers.User) error
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo *repository.Manager
}

func New(repo *repository.Manager) *service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, user *modelUsers.UserInput) error {
	return s.repo.Users.Create(ctx, user)
}

func (s *service) GetByID(ctx context.Context, id int64) (*modelUsers.User, error) {
	return s.repo.Users.GetByID(ctx, id)
}

func (s *service) GetByIDFullInfo(ctx context.Context, id int64) (*modelUsers.UserFullInfo, error) {
	return s.repo.Users.GetByIDFullInfo(ctx, id)
}

func (s *service) GetByEmail(ctx context.Context, email string) (*modelUsers.User, error) {
	return s.repo.Users.GetByEmail(ctx, email)
}

func (s *service) Update(ctx context.Context, user *modelUsers.User) error {
	return s.repo.Users.Update(ctx, user)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.repo.Users.Delete(ctx, id)
}
