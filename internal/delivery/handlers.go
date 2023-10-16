package delivery

import (
	"doodocsbackendchallenge/internal/service"

	"github.com/go-chi/chi"
	"golang.org/x/exp/slog"
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

func (h *Handler) Handlers() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", h.Homepage)
	r.Post("/uploadfile", h.UploadFile)
	r.Post("/acrhivefiles", h.ArchiveFiles)
	r.Post("/sendemailandfile", h.SendEmailsAndFile)
	return r
}
