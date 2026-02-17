package routes

import (
	"interview-project/internal/routes/auth"
	"interview-project/internal/routes/health"
	"interview-project/internal/routes/questions"
	"interview-project/internal/routes/users"
	"interview-project/internal/service"
	"net/http"
)

func New(svc *service.Manager) *http.ServeMux {
	router := http.NewServeMux()

	// Register routes
	router.Handle("/auth/", http.StripPrefix("/auth", auth.RegisterRoutes(&svc.Auth)))
	router.Handle("/users/", http.StripPrefix("/users", users.RegisterRoutes(&svc.Users)))
	router.Handle("/health/", http.StripPrefix("/health", health.RegisterRoutes()))
	router.Handle("/questions/", http.StripPrefix("/questions", questions.RegisterRoutes(&svc.Questions)))

	// Static pages
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/main-page.html")
	})
	router.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/login-page.html")
	})

	return router
}
