package service

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"doodocsbackendchallenge/models"

	"github.com/gabriel-vasile/mimetype"
)

type FileService struct {
	log *slog.Logger
}

type File interface{}

func newFileService(log *slog.Logger) *FileService {
	return &FileService{
		log: log,
	}
}

type mygmail struct {
	email    string
	password string
}

func (flsrv *FileService) UploadFileGetJSON(f multipart.File, header *multipart.FileHeader) (models.Archive, error) {
	const op = "service.UploadFileGetJSON"
	log := flsrv.log.With(
		slog.String("op", op),
	)
	zipfile, err := zip.NewReader(f, header.Size)
	if err != nil {
		log.Error(err.Error())
		return models.Archive{}, fmt.Errorf("%s: %w\n", op, err)
	}
	archivesize := float64(header.Size)
	files, totalarchivesize, err := getFileinArchive(zipfile)
	if err != nil {
		return models.Archive{}, fmt.Errorf("%s: %w\n", op, err)
	}
	archive := models.Archive{
		Filename:     header.Filename,
		Archive_size: archivesize,
		Total_size:   totalarchivesize,
		Total_files:  float64(len(zipfile.File)),
		Files:        files,
	}
	return archive, nil
}

func isZIPFile(file multipart.File) error {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return err
	}

	fileType := http.DetectContentType(buffer)
	if fileType != "application/zip" {
		return fmt.Errorf("Is not zip file")
	}
	return nil
}

func getFileinArchive(zip *zip.Reader) ([]models.File, float64, error) {
	dest := "archive"
	totalsizearchive := 0.0
	files := []models.File{}
	for _, file := range zip.File {
		fp := filepath.Join(dest, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(fp, os.ModePerm)
			continue
		}
		os.MkdirAll(filepath.Dir(fp), os.ModePerm)
		strings.HasPrefix(fp, filepath.Clean(dest)+string(os.PathSeparator))
		rc, err := file.Open()
		if err != nil {
			return nil, 0.0, err
		}
		defer rc.Close()

		buffer := make([]byte, 512)
		n, err := rc.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, 0.0, err
		}
		mtype := strings.Split(mimetype.Detect(buffer[:n]).String(), ";")[0]

		fmt.Println(mtype)
		filesize := float64(file.FileInfo().Size())
		file := models.File{
			File_path: fp,
			Size:      filesize,
			Mimetype:  mtype,
		}
		files = append(files, file)
		totalsizearchive += filesize
	}
	return files, totalsizearchive, nil
}

func (flsrv *FileService) ArchiveInFiles(files []*multipart.FileHeader) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	const op = "service.ArchiveInFiles"
	zipWriter := zip.NewWriter(buf)
	defer zipWriter.Close()
	for _, file := range files {
		mtype := mime.TypeByExtension(filepath.Ext(file.Filename))
		switch mtype {
		case "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/xml", "image/jpeg", "image/png":
			src, err := file.Open()
			if err != nil {
				log.Println(err)
				return nil, fmt.Errorf("%s: %w\n", op, err)
			}
			defer src.Close()

			zipEntry, err := zipWriter.Create(file.Filename)
			if err != nil {
				log.Println(err)
				return nil, fmt.Errorf("%s: %w\n", op, err)
			}

			_, err = io.Copy(zipEntry, src)
			if err != nil {
				log.Println(err)
				return nil, fmt.Errorf("%s: %w\n", op, err)
			}
		default:
			return nil, fmt.Errorf("%s: %s\n", op, "Wrong mime type")
		}
	}
	return buf, nil
}
