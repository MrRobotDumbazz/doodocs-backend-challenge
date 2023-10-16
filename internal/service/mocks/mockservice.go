package mock_service

import (
	gomock "github.com/golang/mock/gomock"
)

type MockFileService struct {
	ctrl     *gomock.Controller
	recorder *MockFileServiceMockRecorder
}

type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

type MockServiceMockRecorder struct {
	mock *MockService
}

type MockFileServiceMockRecorder struct {
	mock *MockFileService
}

func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

func (m *MockService) SomeMethod() string {
	// Mock your Service method behavior here.
	return "mocked result"
}

func (m *MockService) GetFileService() *MockService {
	// Return the MockFileService.
	return m
}

func NewMockFileService(ctrl *gomock.Controller) *MockFileService {
	mock := &MockFileService{ctrl: ctrl}
	mock.recorder = &MockFileServiceMockRecorder{mock}
	return mock
}

func (m *MockFileService) EXPECT() *MockFileServiceMockRecorder {
	return m.recorder
}
