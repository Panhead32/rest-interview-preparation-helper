package service

import (
	"interview-project/internal/repository"
	authService "interview-project/internal/service/auth"
	questionsService "interview-project/internal/service/questions"
	usersService "interview-project/internal/service/users"
)

type Manager struct {
	Auth      authService.AuthService
	Users     usersService.UserService
	Questions questionsService.QuestionService
}

func New(repo *repository.Manager) *Manager {
	return &Manager{
		Auth:      authService.New(repo),
		Users:     usersService.New(repo),
		Questions: questionsService.New(repo),
	}
}
