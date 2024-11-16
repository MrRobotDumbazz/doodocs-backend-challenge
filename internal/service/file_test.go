package service

import (
	"archive/zip"
	"bytes"
	"log/slog"
	"mime/multipart"
	"os"
	"testing"

	"doodocsbackendchallenge/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestZip создает тестовый zip файл с заданным содержимым
func createTestZip(t *testing.T, files map[string][]byte) (*bytes.Buffer, int64) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	var totalSize int64
	for name, content := range files {
		writer, err := zipWriter.Create(name)
		require.NoError(t, err)

		_, err = writer.Write(content)
		require.NoError(t, err)

		totalSize += int64(len(content))
	}

	require.NoError(t, zipWriter.Close())
	return buf, totalSize
}

// createMultipartFile создает multipart.File из буфера
func createMultipartFile(t *testing.T, buf *bytes.Buffer, filename string) (multipart.File, *multipart.FileHeader) {
	file := &multipart.FileHeader{
		Filename: filename,
		Size:     int64(buf.Len()),
	}

	return MultipartFileFromBuffer(buf), file
}

// MultipartFileFromBuffer реализует интерфейс multipart.File
type MultipartFileFromBuffermodel struct {
	*bytes.Reader
}

func (m *MultipartFileFromBuffermodel) Close() error {
	return nil
}

func MultipartFileFromBuffer(buf *bytes.Buffer) multipart.File {
	return &MultipartFileFromBuffermodel{bytes.NewReader(buf.Bytes())}
}

func TestUploadFileGetJSON(t *testing.T) {
	tests := []struct {
		name        string
		files       map[string][]byte
		wantErr     bool
		checkResult func(*testing.T, models.Archive)
	}{
		{
			name: "успешная загрузка с одним файлом",
			files: map[string][]byte{
				"test.txt": []byte("test content"),
			},
			wantErr: false,
			checkResult: func(t *testing.T, archive models.Archive) {
				assert.Equal(t, "test.zip", archive.Filename)
				assert.Equal(t, float64(1), archive.Total_files)
				assert.Len(t, archive.Files, 1)
				assert.Equal(t, "archive/test.txt", archive.Files[0].File_path)
				assert.Equal(t, float64(len("test content")), archive.Files[0].Size)
				assert.Equal(t, "text/plain", archive.Files[0].Mimetype)
			},
		},
		{
			name: "успешная загрузка с несколькими файлами",
			files: map[string][]byte{
				"test1.txt":        []byte("content1"),
				"test2.txt":        []byte("content2"),
				"folder/test3.txt": []byte("content3"),
			},
			wantErr: false,
			checkResult: func(t *testing.T, archive models.Archive) {
				assert.Equal(t, float64(3), archive.Total_files)
				assert.Len(t, archive.Files, 3)
				totalSize := float64(len("content1") + len("content2") + len("content3"))
				assert.Equal(t, totalSize, archive.Total_size)
			},
		},
		{
			name:    "пустой архив",
			files:   map[string][]byte{},
			wantErr: false,
			checkResult: func(t *testing.T, archive models.Archive) {
				assert.Equal(t, float64(0), archive.Total_files)
				assert.Len(t, archive.Files, 0)
				assert.Equal(t, float64(0), archive.Total_size)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый zip файл
			buf, _ := createTestZip(t, tt.files)
			file, header := createMultipartFile(t, buf, "test.zip")

			// Создаем сервис
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
			service := &FileService{log: logger}

			// Вызываем тестируемый метод
			result, err := service.UploadFileGetJSON(file, header)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			tt.checkResult(t, result)
		})
	}
}

func TestIsZIPFile(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
		wantErr bool
	}{
		{
			name:    "валидный ZIP файл",
			content: createValidZipContent(t),
			wantErr: false,
		},
		{
			name:    "невалидный файл (текст)",
			content: []byte("это не zip файл"),
			wantErr: true,
		},
		{
			name:    "пустой файл",
			content: []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := MultipartFileFromBuffer(bytes.NewBuffer(tt.content))
			err := isZIPFile(file)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetFileinArchive(t *testing.T) {
	tests := []struct {
		name          string
		files         map[string][]byte
		expectedFiles int
		expectedSize  float64
	}{
		{
			name: "обычные файлы",
			files: map[string][]byte{
				"file1.txt": []byte("content1"),
				"file2.txt": []byte("content2"),
			},
			expectedFiles: 2,
			expectedSize:  float64(len("content1") + len("content2")),
		},
		{
			name: "файлы в подпапках",
			files: map[string][]byte{
				"folder1/file1.txt": []byte("content1"),
				"folder2/file2.txt": []byte("content2"),
			},
			expectedFiles: 2,
			expectedSize:  float64(len("content1") + len("content2")),
		},
		{
			name:          "пустой архив",
			files:         map[string][]byte{},
			expectedFiles: 0,
			expectedSize:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый zip
			buf, _ := createTestZip(t, tt.files)

			// Создаем zip.Reader
			zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
			require.NoError(t, err)

			// Вызываем тестируемую функцию
			files, totalSize, err := getFileinArchive(zipReader)
			require.NoError(t, err)

			// Проверяем результаты
			assert.Len(t, files, tt.expectedFiles)
			assert.Equal(t, tt.expectedSize, totalSize)

			// Проверяем корректность путей и размеров файлов
			for _, file := range files {
				assert.Contains(t, file.File_path, "archive/")
				assert.Greater(t, file.Size, float64(0))
				assert.NotEmpty(t, file.Mimetype)
			}
		})
	}
}

// createValidZipContent создает минимальный валидный ZIP-файл
func createValidZipContent(t *testing.T) []byte {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	writer, err := zipWriter.Create("test.txt")
	require.NoError(t, err)
	_, err = writer.Write([]byte("test"))
	require.NoError(t, err)
	require.NoError(t, zipWriter.Close())
	return buf.Bytes()
}
