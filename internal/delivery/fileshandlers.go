package delivery

import (
	"log/slog"
	"net/http"
	"strings"
)

func (h *Handler) uploadfile_inarchive(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) archive_files(w http.ResponseWriter, r *http.Request) {
	const op = "delivery.archive_files"
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
		files := r.MultipartForm.File["filestoarchive"]
		buf, err := h.services.ArchiveInFiles(files)
		if err != nil {
			log.Error(err.Error())
			h.EncodeJSON(w, r, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Diposition", "attachment; filename=archive.zip")
		w.Write(buf.Bytes())
	default:
		h.EncodeJSON(w, r, http.StatusMethodNotAllowed, nil)
	}
}

func (h *Handler) sendemails_file(w http.ResponseWriter, r *http.Request) {
	const op = "delivery.sendemails_file"
	log := h.log.With(
		slog.String("op", op),
	)
	switch r.Method {
	case "POST":
		file, fileHeader, err := r.FormFile("fileToGetEmail")
		if err != nil {
			log.Error(err.Error())
			h.EncodeJSON(w, r, http.StatusInternalServerError, err)
			return
		}
		defer file.Close()
		emails := strings.Split(r.FormValue("emails"), ", ")
		log.Info("emails:", emails)
		if err := h.services.GetEmailAndFileSendEmail(emails, file, fileHeader.Filename); err != nil {
			log.Error(err.Error())
			h.EncodeJSON(w, r, http.StatusInternalServerError, err)
			return
		}
	default:
		h.EncodeJSON(w, r, http.StatusMethodNotAllowed, nil)
	}
}
