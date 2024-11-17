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
	mux.HandleFunc("/archivefiles", h.archive_files)
	mux.HandleFunc("/sendemailandfile", h.sendemails_file)
	return mux
}
