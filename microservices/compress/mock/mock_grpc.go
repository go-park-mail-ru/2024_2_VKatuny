// Code generated by MockGen. DO NOT EDIT.
// Source: compress_grpc.pb.go
//
// Generated by this command:
//
//	mockgen -source=compress_grpc.pb.go -destination=../mock/mock_grpc.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	compressmicroservice "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/generated"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockCompressServiceClient is a mock of CompressServiceClient interface.
type MockCompressServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockCompressServiceClientMockRecorder
	isgomock struct{}
}

// MockCompressServiceClientMockRecorder is the mock recorder for MockCompressServiceClient.
type MockCompressServiceClientMockRecorder struct {
	mock *MockCompressServiceClient
}

// NewMockCompressServiceClient creates a new mock instance.
func NewMockCompressServiceClient(ctrl *gomock.Controller) *MockCompressServiceClient {
	mock := &MockCompressServiceClient{ctrl: ctrl}
	mock.recorder = &MockCompressServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCompressServiceClient) EXPECT() *MockCompressServiceClientMockRecorder {
	return m.recorder
}

// CompressAndSaveFile mocks base method.
func (m *MockCompressServiceClient) CompressAndSaveFile(ctx context.Context, in *compressmicroservice.CompressAndSaveFileInput, opts ...grpc.CallOption) (*compressmicroservice.Nothing, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CompressAndSaveFile", varargs...)
	ret0, _ := ret[0].(*compressmicroservice.Nothing)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CompressAndSaveFile indicates an expected call of CompressAndSaveFile.
func (mr *MockCompressServiceClientMockRecorder) CompressAndSaveFile(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompressAndSaveFile", reflect.TypeOf((*MockCompressServiceClient)(nil).CompressAndSaveFile), varargs...)
}

// DeleteFile mocks base method.
func (m *MockCompressServiceClient) DeleteFile(ctx context.Context, in *compressmicroservice.DeleteFileInput, opts ...grpc.CallOption) (*compressmicroservice.Nothing, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteFile", varargs...)
	ret0, _ := ret[0].(*compressmicroservice.Nothing)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFile indicates an expected call of DeleteFile.
func (mr *MockCompressServiceClientMockRecorder) DeleteFile(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFile", reflect.TypeOf((*MockCompressServiceClient)(nil).DeleteFile), varargs...)
}

// MockCompressServiceServer is a mock of CompressServiceServer interface.
type MockCompressServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockCompressServiceServerMockRecorder
	isgomock struct{}
}

// MockCompressServiceServerMockRecorder is the mock recorder for MockCompressServiceServer.
type MockCompressServiceServerMockRecorder struct {
	mock *MockCompressServiceServer
}

// NewMockCompressServiceServer creates a new mock instance.
func NewMockCompressServiceServer(ctrl *gomock.Controller) *MockCompressServiceServer {
	mock := &MockCompressServiceServer{ctrl: ctrl}
	mock.recorder = &MockCompressServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCompressServiceServer) EXPECT() *MockCompressServiceServerMockRecorder {
	return m.recorder
}

// CompressAndSaveFile mocks base method.
func (m *MockCompressServiceServer) CompressAndSaveFile(arg0 context.Context, arg1 *compressmicroservice.CompressAndSaveFileInput) (*compressmicroservice.Nothing, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompressAndSaveFile", arg0, arg1)
	ret0, _ := ret[0].(*compressmicroservice.Nothing)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CompressAndSaveFile indicates an expected call of CompressAndSaveFile.
func (mr *MockCompressServiceServerMockRecorder) CompressAndSaveFile(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompressAndSaveFile", reflect.TypeOf((*MockCompressServiceServer)(nil).CompressAndSaveFile), arg0, arg1)
}

// DeleteFile mocks base method.
func (m *MockCompressServiceServer) DeleteFile(arg0 context.Context, arg1 *compressmicroservice.DeleteFileInput) (*compressmicroservice.Nothing, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFile", arg0, arg1)
	ret0, _ := ret[0].(*compressmicroservice.Nothing)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFile indicates an expected call of DeleteFile.
func (mr *MockCompressServiceServerMockRecorder) DeleteFile(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFile", reflect.TypeOf((*MockCompressServiceServer)(nil).DeleteFile), arg0, arg1)
}

// mustEmbedUnimplementedCompressServiceServer mocks base method.
func (m *MockCompressServiceServer) mustEmbedUnimplementedCompressServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedCompressServiceServer")
}

// mustEmbedUnimplementedCompressServiceServer indicates an expected call of mustEmbedUnimplementedCompressServiceServer.
func (mr *MockCompressServiceServerMockRecorder) mustEmbedUnimplementedCompressServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedCompressServiceServer", reflect.TypeOf((*MockCompressServiceServer)(nil).mustEmbedUnimplementedCompressServiceServer))
}

// MockUnsafeCompressServiceServer is a mock of UnsafeCompressServiceServer interface.
type MockUnsafeCompressServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeCompressServiceServerMockRecorder
	isgomock struct{}
}

// MockUnsafeCompressServiceServerMockRecorder is the mock recorder for MockUnsafeCompressServiceServer.
type MockUnsafeCompressServiceServerMockRecorder struct {
	mock *MockUnsafeCompressServiceServer
}

// NewMockUnsafeCompressServiceServer creates a new mock instance.
func NewMockUnsafeCompressServiceServer(ctrl *gomock.Controller) *MockUnsafeCompressServiceServer {
	mock := &MockUnsafeCompressServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeCompressServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeCompressServiceServer) EXPECT() *MockUnsafeCompressServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedCompressServiceServer mocks base method.
func (m *MockUnsafeCompressServiceServer) mustEmbedUnimplementedCompressServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedCompressServiceServer")
}

// mustEmbedUnimplementedCompressServiceServer indicates an expected call of mustEmbedUnimplementedCompressServiceServer.
func (mr *MockUnsafeCompressServiceServerMockRecorder) mustEmbedUnimplementedCompressServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedCompressServiceServer", reflect.TypeOf((*MockUnsafeCompressServiceServer)(nil).mustEmbedUnimplementedCompressServiceServer))
}
