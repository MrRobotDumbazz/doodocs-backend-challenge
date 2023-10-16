package delivery_test

import (
	"bytes"
	mock_delivery "doodocsbackendchallenge/internal/delivery/mocks"
	"doodocsbackendchallenge/internal/logger"
	mock_service "doodocsbackendchallenge/internal/service/mocks"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	host = "localhost:8082"
)

func TestUploadFileInArchive(t *testing.T) {
	cases := []struct {
		name      string
		filename  string
		respError string
		mockError error
	}{
		{
			name:     "Success",
			filename: "archive.zip",
		},
		{
			name:      "Not zip",
			filename:  "notzip.txt",
			respError: "Is not zip file",
		},
		{
			name:      "Empty filename",
			filename:  "",
			respError: "field file is a required",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			logger := logger.SetupLogger()
			ctrl := gomock.NewController(t)
			mock := new(mock_service.MockFile)
			header := &multipart.FileHeader{Filename: tc.filename}
			buf := new(bytes.Buffer)
			writer := multipart.NewWriter(buf)
			dataPart, err := writer.CreateFormFile("myFile", tc.filename)
			services := mock_service.NewMockService(ctrl)
			handlers := mock_delivery.NewHandler(services, logger)
			require.NoError(t, err)
			f, err := os.Open(tc.filename)
			require.NoError(t, err)
			_, err = io.Copy(dataPart, f)
			require.NoError(t, err)
			require.NoError(t, writer.Close())
			require.NoError(t, err)
			mockarchive, err := mock.UploadFileGetJSON(f, header)
			require.NoError(t, err)
			input := fmt.Sprintf(`{"filename":%s}`, tc.filename)
			req, err := http.NewRequest(http.MethodPost, "uploadfileinarchive", bytes.NewReader([]byte(input)))
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			handlers.ServeHTTP(rr, req)
			require.Equal(t, rr.Code, http.StatusOK)
			body := rr.Body.String()
			require.NoError(t, json.Unmarshal([]byte(body), &mockarchive))
		})
	}
}

//go:generate moq -out handler_mock_test.go . Handler
