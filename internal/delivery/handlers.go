package delivery

import (
	"log/slog"
	"net/http"

	"doodocsbackendchallenge/internal/service"
)

type Handler struct {
	services *service.Service
	log      *slog.Logger
}

func NewHandler(services *service.Service, log *slog.Logger) *Handler {
	return &Handler{
		services: services,
		log:      log,
	}
}

func (h *Handler) Handlers() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/uploadfilearchive", h.uploadfile_inarchive)
	return mux
}
