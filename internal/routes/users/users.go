package users

import (
	userService "interview-project/internal/service/users"
	"net/http"
)

type Handler struct {
	userSvc userService.UserService
}

func NewHandler(userSvc *userService.UserService) *Handler {
	return &Handler{userSvc: *userSvc}
}

func (h *Handler) getUserByID(res http.ResponseWriter, req *http.Request) {

}

func RegisterRoutes(userSvc *userService.UserService) *http.ServeMux {
	mux := http.NewServeMux()

	handler := NewHandler(userSvc)

	mux.HandleFunc("GET /", handler.getUserByID)

	return mux
}
