package service

import (
	"archive/zip"
	"doodocsbackendchallenge/models"
	"encoding/binary"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type UploadFileService struct{}

type Upload interface {
	UploadFile(filename string) (*models.Archive, error)
}

func newUploadService() *UploadFileService {
	return &UploadFileService{}
}

func (upload *UploadFileService) UploadFile(filename string) (*models.Archive, error) {
	const op = "service.UploadFile"
	zipfile, err := zip.OpenReader(filename)
	defer zipfile.Close()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	dest := "archive"
	files := []models.File{}
	totalsizearchive := 0.0
	archivebytessize, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	bits := binary.LittleEndian.Uint64(archivebytessize)
	archivesizefloat := math.Float64frombits(bits)
	for _, file := range zipfile.File {
		fp := filepath.Join(dest, file.Name)
		strings.HasPrefix(fp, filepath.Clean(dest)+string(os.PathSeparator))
		bytessize, err := os.ReadFile(file.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		mtype := http.DetectContentType(bytessize)
		bits := binary.LittleEndian.Uint64(bytessize)
		sizefloat := math.Float64frombits(bits)
		filestruct := models.File{
			File_path: fp,
			Size:      sizefloat,
			Mimetype:  mtype,
		}
		files = append(files, filestruct)
		totalsizearchive += sizefloat
	}
	archive := &models.Archive{
		Filename:     filename,
		Archive_size: archivesizefloat,
		Total_size:   totalsizearchive,
		Files:        files,
	}
	return archive, nil
}
