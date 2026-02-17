package auth

import (
	"context"
	"errors"
	"interview-project/internal/models/auth"
	modelUsers "interview-project/internal/models/users"
	"interview-project/internal/repository"
	"interview-project/pkg/utils"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, name, surname, nickname, email, password string) (string, error)
	ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error
}

type service struct {
	repo *repository.Manager
}

func New(repo *repository.Manager) *service {
	return &service{repo: repo}
}

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	authData, err := s.repo.Auth.GetByEmail(ctx, email)
	if err != nil || authData == nil {
		return "", errors.New("invalid email or password")
	}

	if err := utils.VerifyPassword(authData.PasswordHash, password); err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(authData.UserID, authData.Email)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (s *service) Register(ctx context.Context, name, surname, nickname, email, password string) (string, error) {
	authData, err := s.repo.Auth.GetByEmail(ctx, email)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return "", err
	}

	if authData != nil {
		return "", errors.New("email already in use")
	}

	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	userInput := &modelUsers.UserInput{
		Name:     name,
		Surname:  surname,
		NickName: nickname,
	}

	if err := s.repo.Users.Create(ctx, userInput); err != nil {
		return "", errors.New("failed to create user")
	}

	user, err := s.repo.Users.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("failed to retrieve created user")
	}

	authInput := &auth.AuthInput{
		Email:    email,
		Password: password,
	}

	if err := s.repo.Auth.Create(ctx, user.ID, authInput, passwordHash); err != nil {
		return "", errors.New("failed to create auth record")
	}

	token, err := utils.GenerateToken(user.ID, email)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (s *service) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	authData, err := s.repo.Auth.GetByUserID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := utils.VerifyPassword(authData.PasswordHash, oldPassword); err != nil {
		return errors.New("invalid current password")
	}

	newPasswordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	authInput := &auth.AuthInput{
		Email:    authData.Email,
		Password: newPassword,
	}

	if err := s.repo.Auth.Update(ctx, userID, authInput, newPasswordHash); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}
