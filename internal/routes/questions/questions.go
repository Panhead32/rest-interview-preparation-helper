package questions

import (
	authMiddleware "interview-project/internal/middleware/auth"
	"interview-project/internal/models/questions"
	model "interview-project/internal/models/response"
	service "interview-project/internal/service/questions"
	"interview-project/pkg/utils"
	"net/http"
	"strconv"
)

type Handler struct {
	svc service.QuestionService
}

func NewHandler(svc *service.QuestionService) *Handler {
	return &Handler{svc: *svc}
}

type parseQuestionsRequest struct {
	Link string `json:"link" validate:"required,url"`
}

type parseQuestionsResponse struct {
	Questions []questions.Question `json:"questions"`
}

func (h *Handler) parseQuestions(res http.ResponseWriter, req *http.Request) {
	var payload parseQuestionsRequest
	if err := utils.DecodeJSONBody(res, req, &payload); err != nil {
		utils.WriteJSON(res, http.StatusBadRequest, model.NewErrorResponse("invalid request body", err.Error()))
		return
	}

	if err := utils.ValidateStruct(&payload); err != nil {
		utils.WriteJSON(res, http.StatusBadRequest, model.NewErrorResponse("validation error", err.Error()))
		return
	}

	questions, err := h.svc.ParseQuestions(req.Context(), payload.Link)
	if err != nil {
		utils.WriteJSON(res, http.StatusInternalServerError, model.NewErrorResponse("failed to fetch questions", err.Error()))
		return
	}

	utils.WriteJSON(res, http.StatusOK, model.NewSuccessResponse("ok", parseQuestionsResponse{Questions: questions}))
}

type getQuestionsRequest struct {
}

type getQuestionsResponse struct {
	Questions []questions.Question `json:"questions"`
}

func (h *Handler) getQuestions(res http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value(utils.UserIDKey).(int64)

	questions, err := h.svc.GetQuestions(req.Context(), userID)
	if err != nil {
		utils.WriteJSON(res, http.StatusInternalServerError, model.NewErrorResponse("failed to fetch questions", err.Error()))
		return
	}

	utils.WriteJSON(res, http.StatusOK, model.NewSuccessResponse("ok", getQuestionsResponse{Questions: questions}))
}

func (h *Handler) getQuestionExplanation(res http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("id")
	if idStr == "" {
		utils.WriteJSON(res, http.StatusBadRequest, model.NewErrorResponse("invalid request", "missing question id"))
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteJSON(res, http.StatusBadRequest, model.NewErrorResponse("invalid request", "question id must be a number"))
		return
	}

	question, err := h.svc.GetQuestionByID(req.Context(), id)
	if err != nil {
		utils.WriteJSON(res, http.StatusInternalServerError, model.NewErrorResponse("failed to fetch question", err.Error()))
		return
	}

	utils.WriteJSON(res, http.StatusOK, model.NewSuccessResponse("ok", &question))
}

func (h *Handler) explainQuestion(res http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("id")
	if idStr == "" {
		utils.WriteJSON(res, http.StatusBadRequest, model.NewErrorResponse("invalid request", "missing question id"))
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteJSON(res, http.StatusBadRequest, model.NewErrorResponse("invalid request", "question id must be a number"))
		return
	}

	err = h.svc.ExplainQuestion(req.Context(), id)
	if err != nil {
		utils.WriteJSON(res, http.StatusInternalServerError, model.NewErrorResponse("failed to explain question", err.Error()))
		return
	}

	question, err := h.svc.GetQuestionByID(req.Context(), id)
	if err != nil {
		utils.WriteJSON(res, http.StatusInternalServerError, model.NewErrorResponse("failed to fetch question", err.Error()))
		return
	}

	utils.WriteJSON(res, http.StatusOK, model.NewSuccessResponse("ok", &question))
}

func RegisterRoutes(questionSvc *service.QuestionService) *http.ServeMux {
	router := http.NewServeMux()
	handler := NewHandler(questionSvc)

	router.Handle("POST /", authMiddleware.AuthMiddleware(http.HandlerFunc(handler.parseQuestions)))
	router.Handle("GET /", authMiddleware.AuthMiddleware(http.HandlerFunc(handler.getQuestions)))
	router.Handle("GET /{id}", authMiddleware.AuthMiddleware(http.HandlerFunc(handler.getQuestionExplanation)))
	router.Handle("POST /{id}/explain", authMiddleware.AuthMiddleware(http.HandlerFunc(handler.explainQuestion)))

	return router
}
