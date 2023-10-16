package service

import (
	"archive/zip"
	"bytes"
	"crypto/tls"
	"doodocsbackendchallenge/internal/config"
	"doodocsbackendchallenge/models"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"mime"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jordan-wright/email"
)

type FileService struct {
	cfg *config.Config
}

type File interface {
	UploadFileGetJSON(file multipart.File, header *multipart.FileHeader) (models.Archive, error)
	ArchiveInFiles(files []*multipart.FileHeader) (*bytes.Buffer, error)
	GetEmailsAndFileSendEmail(emails []string, file io.Reader, filename string) error
}

func newFileService(cfg *config.Config) *FileService {
	return &FileService{
		cfg: cfg,
	}
}

type mygmail struct {
	email    string `env:"EMAIL"`
	password string `env:"PASSWORD"`
}

func (flsrv *FileService) UploadFileGetJSON(f multipart.File, header *multipart.FileHeader) (models.Archive, error) {
	const op = "service.UploadFile"
	zipfile, err := zip.NewReader(f, header.Size)
	if err != nil {
		log.Println(err)
		return models.Archive{}, fmt.Errorf("%s: %w\n", op, err)
	}
	if err := isZIPFile(f); err != nil {
		return models.Archive{}, fmt.Errorf("%s: %w\n", op, err)
	}
	dest := "archive"
	dataarchive, err := io.ReadAll(f)
	if err != nil {
		return models.Archive{}, fmt.Errorf("%s: %w\n", op, err)
	}

	files := []models.File{}
	totalsizearchive := 0.0
	if err != nil {
		return models.Archive{}, fmt.Errorf("%s: %w\n", op, err)
	}
	bitsarchive := binary.LittleEndian.Uint64(dataarchive)
	archivesizefloat := math.Float64frombits(bitsarchive)
	for _, file := range zipfile.File {
		fp := filepath.Join(dest, file.Name)
		strings.HasPrefix(fp, filepath.Clean(dest)+string(os.PathSeparator))
		mtype := mime.TypeByExtension(filepath.Ext(file.Name))
		filereader, err := file.Open()
		if err != nil {
			return models.Archive{}, fmt.Errorf("%s: %w\n", op, err)
		}
		defer filereader.Close()
		datafile, err := io.ReadAll(filereader)
		if err != nil {
			return models.Archive{}, fmt.Errorf("%s: %w\n", op, err)
		}
		bitsfiles := binary.LittleEndian.Uint64(datafile)
		filesizefloat := math.Float64frombits(bitsfiles)
		filestruct := models.File{
			File_path: fp,
			Size:      filesizefloat,
			Mimetype:  mtype,
		}
		files = append(files, filestruct)
		totalsizearchive += filesizefloat
	}
	archive := models.Archive{
		Filename:     header.Filename,
		Archive_size: archivesizefloat,
		Total_size:   totalsizearchive,
		Total_files:  float64(len(zipfile.File)),
		Files:        files,
	}
	return archive, nil
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

func (flsrv *FileService) GetEmailsAndFileSendEmail(emails []string, file io.Reader, filename string) error {
	mtype := mime.TypeByExtension(filepath.Ext(filename))
	op := "service.GetEmailsAndFileSendEmail"
	switch mtype {
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/pdf":
		smtpstruct := struct {
			smtpServer string
			smtpPort   string
		}{
			smtpServer: "smtp.gmail.com",
			smtpPort:   "587",
		}
		auth := smtp.PlainAuth("", flsrv.cfg.Email, flsrv.cfg.Password, smtpstruct.smtpServer)
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         smtpstruct.smtpServer,
		}
		client, err := smtp.Dial(smtpstruct.smtpServer + ":" + smtpstruct.smtpPort)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		defer client.Quit()
		err = client.StartTLS(tlsConfig)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		err = client.Auth(auth)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		log.Printf("Emails:%#v", emails)
		for _, oneemail := range emails {
			log.Printf("email: %s\n", oneemail)
			validEmail, err := regexp.MatchString(`[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`, oneemail)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			if !validEmail {
				return fmt.Errorf("%s: %s", op, "wrong email")
			}
			e := email.NewEmail()
			log.Println()
			e.From = "sansskeleton415@gmail.com"
			e.To = []string{oneemail}
			e.Subject = "Dodocs Backend Challenge"
			e.Text = []byte("Hello i'm testing smtp.")
			_, fn := filepath.Split(filename)
			log.Printf("emails %s\n", emails)
			e.Attach(file, fn, "application/octet-stream")
			err = e.Send(smtpstruct.smtpServer+":"+smtpstruct.smtpPort, auth)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}
	default:
		return fmt.Errorf("%s, %s", op, "Wrong mime type")
	}
	return nil
}
