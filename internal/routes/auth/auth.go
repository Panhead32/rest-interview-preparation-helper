package auth

import (
	authMiddleware "interview-project/internal/middleware/auth"
	model "interview-project/internal/models/response"
	authService "interview-project/internal/service/auth"
	"interview-project/pkg/utils"
	"net/http"
)

type Handler struct {
	svc authService.AuthService
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type SignUpRequest struct {
	Name           string `json:"name" validate:"required,min=2,max=50"`
	Surname        string `json:"surname" validate:"required,min=2,max=50"`
	Nickname       string `json:"nickname" validate:"required,min=3,max=30,alphanum"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=6,max=100"`
	PasswordRepeat string `json:"password_repeat" validate:"required,eqfield=Password"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=6"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=100"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func NewHandler(authSvc *authService.AuthService) *Handler {
	return &Handler{svc: *authSvc}
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, model.NewErrorResponse("invalid request body", err.Error()))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, model.NewErrorResponse("validation error", err.Error()))
		return
	}

	token, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, model.NewErrorResponse("authentication failed", err.Error()))
		return
	}

	utils.WriteJSON(w, http.StatusOK, model.NewSuccessResponse("sign in successful", AuthResponse{Token: token}))
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, model.NewErrorResponse("invalid request body", err.Error()))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, model.NewErrorResponse("validation error", err.Error()))
		return
	}

	token, err := h.svc.Register(r.Context(), req.Name, req.Surname, req.Nickname, req.Email, req.Password)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, model.NewErrorResponse("registration failed", err.Error()))
		return
	}

	utils.WriteJSON(w, http.StatusCreated, model.NewSuccessResponse("sign up successful", AuthResponse{Token: token}))
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req ChangePasswordRequest
	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, model.NewErrorResponse("invalid request body", err.Error()))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, model.NewErrorResponse("validation error", err.Error()))
		return
	}

	// Extract userID from JWT token in middleware
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, model.NewErrorResponse("authorization required", "user not authenticated"))
		return
	}

	if err := h.svc.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, model.NewErrorResponse("password change failed", err.Error()))
		return
	}

	utils.WriteJSON(w, http.StatusOK, model.NewSuccessResponse("password changed successfully", nil))
}

func RegisterRoutes(authSvc *authService.AuthService) *http.ServeMux {
	mux := http.NewServeMux()

	handler := NewHandler(authSvc)

	mux.HandleFunc("POST /signin", handler.SignIn)
	mux.HandleFunc("POST /signup", handler.SignUp)

	// Protected routes that require authentication
	mux.Handle("POST /change-password", authMiddleware.AuthMiddleware(http.HandlerFunc(handler.ChangePassword)))

	return mux
}
