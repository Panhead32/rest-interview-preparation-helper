package questions

import (
	model "interview-project/internal/models/response"
	service "interview-project/internal/service/questions"
	"interview-project/pkg/utils"
	"net/http"
)

type Handler struct {
	svc service.QuestionService
}

type parseQuestionsRequest struct {
	Count int `json:"count" validate:"gte=0,lte=1000"`
}

type parseQuestionsResponse struct {
	Questions []string `json:"questions"`
}

func NewHandler(svc *service.QuestionService) *Handler {
	return &Handler{svc: *svc}
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

	questions, err := h.svc.GetQuestions(req.Context())
	if err != nil {
		utils.WriteJSON(res, http.StatusInternalServerError, model.NewErrorResponse("failed to fetch questions", err.Error()))
		return
	}

	if payload.Count > 0 && payload.Count < len(questions) {
		questions = questions[:payload.Count]
	}

	utils.WriteJSON(res, http.StatusOK, model.NewSuccessResponse("ok", parseQuestionsResponse{Questions: questions}))
}

func RegisterRoutes(questionSvc *service.QuestionService) *http.ServeMux {
	mux := http.NewServeMux()
	handler := NewHandler(questionSvc)

	mux.HandleFunc("POST /", handler.parseQuestions)

	return mux
}
