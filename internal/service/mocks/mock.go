// Code generated by MockGen. DO NOT EDIT.
// Source: file.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	bytes "bytes"
	models "doodocsbackendchallenge/models"
	io "io"
	multipart "mime/multipart"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFile is a mock of File interface.
type MockFile struct {
	ctrl     *gomock.Controller
	recorder *MockFileMockRecorder
}

// MockFileMockRecorder is the mock recorder for MockFile.
type MockFileMockRecorder struct {
	mock *MockFile
}

// NewMockFile creates a new mock instance.
func NewMockFile(ctrl *gomock.Controller) *MockFile {
	mock := &MockFile{ctrl: ctrl}
	mock.recorder = &MockFileMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFile) EXPECT() *MockFileMockRecorder {
	return m.recorder
}

// ArchiveInFiles mocks base method.
func (m *MockFile) ArchiveInFiles(files []*multipart.FileHeader) (*bytes.Buffer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ArchiveInFiles", files)
	ret0, _ := ret[0].(*bytes.Buffer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ArchiveInFiles indicates an expected call of ArchiveInFiles.
func (mr *MockFileMockRecorder) ArchiveInFiles(files interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ArchiveInFiles", reflect.TypeOf((*MockFile)(nil).ArchiveInFiles), files)
}

// GetEmailsAndFileSendEmail mocks base method.
func (m *MockFile) GetEmailsAndFileSendEmail(emails []string, file io.Reader, filename string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmailsAndFileSendEmail", emails, file, filename)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetEmailsAndFileSendEmail indicates an expected call of GetEmailsAndFileSendEmail.
func (mr *MockFileMockRecorder) GetEmailsAndFileSendEmail(emails, file, filename interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmailsAndFileSendEmail", reflect.TypeOf((*MockFile)(nil).GetEmailsAndFileSendEmail), emails, file, filename)
}

// UploadFileGetJSON mocks base method.
func (m *MockFile) UploadFileGetJSON(file multipart.File, header *multipart.FileHeader) (models.Archive, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFileGetJSON", file, header)
	ret0, _ := ret[0].(models.Archive)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadFileGetJSON indicates an expected call of UploadFileGetJSON.
func (mr *MockFileMockRecorder) UploadFileGetJSON(file, header interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFileGetJSON", reflect.TypeOf((*MockFile)(nil).UploadFileGetJSON), file, header)
}