package delivery

import (
	"log/slog"
	"net/http"
)

func (h *Handler) UploadFileInArchive(w http.ResponseWriter, r *http.Request) {
	const op = "delivery.uploadfileinarchive"
	log := h.log.With(
		slog.String("op", op),
	)
	switch r.Method {
	case "POST":
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Error(err.Error())
			h.EncodeJSON(w, r, http.StatusInternalServerError, err)
			return
		}
		file, fileheader, err := r.FormFile("myFile")
		if err != nil {
			log.Error(err.Error())
			h.EncodeJSON(w, r, http.StatusInternalServerError, err)
			return
		}
		defer file.Close()
		archive, err := h.services.UploadFileGetJSON(file, fileheader)
		if err != nil {
			log.Error(err.Error())
			h.EncodeJSON(w, r, http.StatusInternalServerError, err)
			return
		}
		h.EncodeJSON(w, r, http.StatusOK, archive)
		log.Debug("archive", archive)
	default:
		h.EncodeJSON(w, r, http.StatusMethodNotAllowed, nil)
	}
}
