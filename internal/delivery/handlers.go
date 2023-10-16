package delivery

import (
	"doodocsbackendchallenge/internal/service"

	"github.com/go-chi/chi"
	"golang.org/x/exp/slog"
)

//go:generate moq -out handler_mock_test.go . Handler
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

func (h *Handler) Handlers() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", h.Homepage)
	r.Post("/uploadfileinarchive", h.UploadFileInArchive)
	r.Post("/acrhivefiles", h.ArchiveFiles)
	r.Post("/sendemailandfile", h.SendEmailsAndFile)
	return r
}
