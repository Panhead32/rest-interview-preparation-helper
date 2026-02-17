package health

import (
	"interview-project/pkg/utils"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{
		Status: "ok",
		Time:   time.Now().UTC().Format(time.RFC3339),
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", checkHealth)

	return mux
}
