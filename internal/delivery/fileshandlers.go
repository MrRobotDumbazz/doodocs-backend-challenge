package delivery

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.uploadfile"
	log := h.log.With(
		slog.String("op", op),
	)
	switch r.Method {
	case "GET":
		temp, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}
		err = temp.Execute(w, nil)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}
		render.JSON(w, r, OK())
	case "POST":
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}
		file, fileheader, err := r.FormFile("myFile")
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}
		defer file.Close()
		archive, err := h.services.UploadFileGetJSON(file, fileheader)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}
		render.JSON(w, r, archive)
		log.Info("archive", archive)
	default:
		Error("MethodNotallowed")
	}
}

func (h *Handler) ArchiveFiles(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.Archivefiles"
	log := h.log.With(
		slog.String("op", op),
	)
	switch r.Method {
	case "POST":
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}
		files := r.MultipartForm.File["filestoarchive"]
		buf, err := h.services.ArchiveInFiles(files)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Diposition", "attachment; filename=archive.zip")
		w.Write(buf.Bytes())
	default:
		Error("MethodNotallowed")
	}
}

func (h *Handler) Homepage(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.homepage"
	log := h.log.With(
		slog.String("op", op),
	)
	switch r.Method {
	case "GET":
		temp, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}
		err = temp.Execute(w, nil)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}
	default:
		Error("MethodNotallowed")
	}
}

func (h *Handler) SendEmailsAndFile(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.SendEmailsAndFile"
	log := h.log.With(
		slog.String("op", op),
	)
	switch r.Method {
	case "POST":
		file, fileHeader, err := r.FormFile("fileToGetEmail")
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}
		defer file.Close()
		emails := strings.Split(r.FormValue("emails"), ", ")
		log.Info("emails:", emails)
		err = h.services.GetEmailsAndFileSendEmail(emails, file, fileHeader.Filename)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Error(err.Error()))
			return
		}
	default:
		Error("MethodNotallowed")
	}
}
