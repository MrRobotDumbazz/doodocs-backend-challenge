package service

import (
	"archive/zip"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"doodocsbackendchallenge/internal/config"
	"doodocsbackendchallenge/models"

	"github.com/gabriel-vasile/mimetype"
	"github.com/jordan-wright/email"
)

type FileService struct {
	log *slog.Logger
	cfg *config.Config
}

type File interface {
	UploadFileGetJSON(f multipart.File, header *multipart.FileHeader) (models.Archive, error)
	ArchiveInFiles(files []*multipart.FileHeader) (*bytes.Buffer, error)
	GetEmailAndFileSendEmail(emails []string, file io.Reader, filename string) error
}

func newFileService(log *slog.Logger, cfg *config.Config) *FileService {
	return &FileService{
		log: log,
		cfg: cfg,
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
		log.Printf("%s: %s", op, filepath.Ext(file.Filename))
		mtype := mime.TypeByExtension(filepath.Ext(file.Filename))
		log.Printf("%s: %s", op, mtype)
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
	if err := zipWriter.Close(); err != nil {
		log.Printf("%s: %v", op, err)
		return nil, fmt.Errorf("%s: %w\n", op, err)
	}
	return buf, nil
}

func (flsrv *FileService) GetEmailAndFileSendEmail(emails []string, file io.Reader, filename string) error {
	mtype := mime.TypeByExtension(filepath.Ext(filename))
	const op = "service.getemailsandfilesendemail"
	switch mtype {
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/pdf":
		smtpstruct := struct {
			smtpServer string
			smtpPort   string
		}{
			smtpServer: "smtp.gmail.com",
			smtpPort:   "587",
		}
		auth, err := startsmtp(smtpstruct, flsrv.cfg)
		if err != nil {
			log.Printf("%s: %v\n", op, err)
			return fmt.Errorf("%s: %v\n", op, err)
		}

		if err := getemail(emails, smtpstruct, filename, file, auth); err != nil {
			log.Printf("%s: %v\n", op, err)
			return fmt.Errorf("%s: %v\n", op, err)
		}
	}
	return nil
}

func startsmtp(smtpstruct struct {
	smtpServer string
	smtpPort   string
},
	config *config.Config,
) (smtp.Auth, error) {
	auth := smtp.PlainAuth("", config.Email, config.Password, smtpstruct.smtpServer)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpstruct.smtpServer,
	}
	client, err := smtp.Dial(smtpstruct.smtpServer + ":" + smtpstruct.smtpPort)
	if err != nil {
		return nil, err
	}
	defer client.Quit()
	err = client.StartTLS(tlsConfig)
	if err != nil {
		return nil, err
	}
	err = client.Auth(auth)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func getemail(emails []string,
	smtpstruct struct {
		smtpServer string
		smtpPort   string
	},
	filename string, file io.Reader, auth smtp.Auth,
) error {
	for _, oneemail := range emails {
		log.Printf("email: %s\n", oneemail)
		validEmail, err := regexp.MatchString(`[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`, oneemail)
		if err != nil {
			return err
		}
		if !validEmail {
			return fmt.Errorf("wrong email")
		}
		e := email.NewEmail()
		log.Println()
		e.From = "sansskeleton415@gmail.com"
		e.To = []string{oneemail}
		e.Subject = "Doodocs Backend Challenge"
		e.Text = []byte("Hello i'm testing smtp.")
		_, fn := filepath.Split(filename)
		log.Printf("Emails %s\n", emails)
		e.Attach(file, fn, "application/octet-stream")
		err = e.Send(smtpstruct.smtpServer+":"+smtpstruct.smtpPort, auth)
		if err != nil {
			return err
		}
	}
	return nil
}
