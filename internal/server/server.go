package server

import (
	"context"
	"database/sql"
	"fmt"
	"interview-project/internal/config"
	"interview-project/internal/middleware"
	"interview-project/internal/middleware/logger"
	pathHandler "interview-project/internal/middleware/path-handler"
	"interview-project/internal/repository"
	"interview-project/internal/routes"
	"interview-project/internal/service"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	config     *config.Config
}

// NewServer creates and initializes a new server
func NewServer(cfg *config.Config, db *sql.DB) *Server {

	stack := middleware.CreateStack(
		logger.LoggingMiddleware,
		pathHandler.PathHandlerMiddleware,
	)

	repositoryHandler := repository.New(db)
	serviceHandler := service.New(repositoryHandler)
	handlers := routes.New(serviceHandler)

	httpServer := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      stack(handlers),
		ReadTimeout:  time.Duration(cfg.Server.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.Timeout) * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		config:     cfg,
	}
}

func (s *Server) Start(ctx context.Context) error {
	fmt.Printf("Starting server on %s\n", s.httpServer.Addr)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	<-ctx.Done()

	return s.Shutdown(context.Background())
}

func (s *Server) Shutdown(ctx context.Context) error {
	fmt.Println("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}
