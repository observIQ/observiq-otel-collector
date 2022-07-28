// Code generated by mockery v2.12.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	protobufs "github.com/open-telemetry/opamp-go/protobufs"
)

// MockDownloadableFileManager is an autogenerated mock type for the DownloadableFileManager type
type MockDownloadableFileManager struct {
	mock.Mock
}

// FetchAndExtractArchive provides a mock function with given fields: _a0
func (_m *MockDownloadableFileManager) FetchAndExtractArchive(_a0 *protobufs.DownloadableFile) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*protobufs.DownloadableFile) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type NewDownloadableFileManagerT interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockDownloadableFileManager creates a new instance of DownloadableFileManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockDownloadableFileManager(t NewDownloadableFileManagerT) *MockDownloadableFileManager {
	mock := &MockDownloadableFileManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}