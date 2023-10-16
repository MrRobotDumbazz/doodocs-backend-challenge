package server

import (
	"context"
	"doodocsbackendchallenge/internal/config"
	"net/http"

	"github.com/go-chi/chi"
)

type Server struct {
	srv http.Server
}

func (s *Server) Start(cfg *config.Config, handlers chi.Router) error {
	s.srv = http.Server{
		Addr:         cfg.Address,
		Handler:      handlers,
		WriteTimeout: cfg.Timeout,
		ReadTimeout:  cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
