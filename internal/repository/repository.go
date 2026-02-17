package repository

import (
	"database/sql"

	authRepo "interview-project/internal/repository/auth"
	questionsRepo "interview-project/internal/repository/questions"
	requestsRepo "interview-project/internal/repository/requests"
	usersRepo "interview-project/internal/repository/users"
)

type Manager struct {
	Users     usersRepo.UsersRepository
	Auth      authRepo.AuthRepository
	Requests  requestsRepo.RequestsRepository
	Questions questionsRepo.QuestionsRepository
}

func New(db *sql.DB) *Manager {
	return &Manager{
		Users:     usersRepo.New(db),
		Auth:      authRepo.New(db),
		Requests:  requestsRepo.New(db),
		Questions: questionsRepo.New(db),
	}
}
