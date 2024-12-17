// Code generated by MockGen. DO NOT EDIT.
// Source: microservices/auth/auth.go
//
// Generated by this command:
//
//	mockgen -source=microservices/auth/auth.go -destination=microservices/auth/mock/mock_delivery.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	auth "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth"
	__auth_microservice "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/auth/gen"
	gomock "go.uber.org/mock/gomock"
)

// MockIAuthorizationDelivery is a mock of IAuthorizationDelivery interface.
type MockIAuthorizationDelivery struct {
	ctrl     *gomock.Controller
	recorder *MockIAuthorizationDeliveryMockRecorder
	isgomock struct{}
}

// MockIAuthorizationDeliveryMockRecorder is the mock recorder for MockIAuthorizationDelivery.
type MockIAuthorizationDeliveryMockRecorder struct {
	mock *MockIAuthorizationDelivery
}

// NewMockIAuthorizationDelivery creates a new mock instance.
func NewMockIAuthorizationDelivery(ctrl *gomock.Controller) *MockIAuthorizationDelivery {
	mock := &MockIAuthorizationDelivery{ctrl: ctrl}
	mock.recorder = &MockIAuthorizationDeliveryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAuthorizationDelivery) EXPECT() *MockIAuthorizationDeliveryMockRecorder {
	return m.recorder
}

// Auth mocks base method.
func (m *MockIAuthorizationDelivery) Auth(ctx context.Context, req *__auth_microservice.AuthRequest) (*__auth_microservice.AuthResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Auth", ctx, req)
	ret0, _ := ret[0].(*__auth_microservice.AuthResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Auth indicates an expected call of Auth.
func (mr *MockIAuthorizationDeliveryMockRecorder) Auth(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Auth", reflect.TypeOf((*MockIAuthorizationDelivery)(nil).Auth), ctx, req)
}

// CheckAuth mocks base method.
func (m *MockIAuthorizationDelivery) CheckAuth(ctx context.Context, req *__auth_microservice.CheckAuthRequest) (*__auth_microservice.CheckAuthResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAuth", ctx, req)
	ret0, _ := ret[0].(*__auth_microservice.CheckAuthResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckAuth indicates an expected call of CheckAuth.
func (mr *MockIAuthorizationDeliveryMockRecorder) CheckAuth(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAuth", reflect.TypeOf((*MockIAuthorizationDelivery)(nil).CheckAuth), ctx, req)
}

// Deauth mocks base method.
func (m *MockIAuthorizationDelivery) Deauth(ctx context.Context, req *__auth_microservice.DeauthRequest) (*__auth_microservice.DeauthResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Deauth", ctx, req)
	ret0, _ := ret[0].(*__auth_microservice.DeauthResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Deauth indicates an expected call of Deauth.
func (mr *MockIAuthorizationDeliveryMockRecorder) Deauth(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deauth", reflect.TypeOf((*MockIAuthorizationDelivery)(nil).Deauth), ctx, req)
}

// MockIAuthorizationRepository is a mock of IAuthorizationRepository interface.
type MockIAuthorizationRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIAuthorizationRepositoryMockRecorder
	isgomock struct{}
}

// MockIAuthorizationRepositoryMockRecorder is the mock recorder for MockIAuthorizationRepository.
type MockIAuthorizationRepositoryMockRecorder struct {
	mock *MockIAuthorizationRepository
}

// NewMockIAuthorizationRepository creates a new mock instance.
func NewMockIAuthorizationRepository(ctrl *gomock.Controller) *MockIAuthorizationRepository {
	mock := &MockIAuthorizationRepository{ctrl: ctrl}
	mock.recorder = &MockIAuthorizationRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAuthorizationRepository) EXPECT() *MockIAuthorizationRepositoryMockRecorder {
	return m.recorder
}

// CreateSession mocks base method.
func (m *MockIAuthorizationRepository) CreateSession(arg0 uint64, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockIAuthorizationRepositoryMockRecorder) CreateSession(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockIAuthorizationRepository)(nil).CreateSession), arg0, arg1)
}

// DeleteSession mocks base method.
func (m *MockIAuthorizationRepository) DeleteSession(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSession", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSession indicates an expected call of DeleteSession.
func (mr *MockIAuthorizationRepositoryMockRecorder) DeleteSession(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSession", reflect.TypeOf((*MockIAuthorizationRepository)(nil).DeleteSession), arg0)
}

// GetUser mocks base method.
func (m *MockIAuthorizationRepository) GetUser(userType, email string) (*auth.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", userType, email)
	ret0, _ := ret[0].(*auth.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockIAuthorizationRepositoryMockRecorder) GetUser(userType, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockIAuthorizationRepository)(nil).GetUser), userType, email)
}

// GetUserIdBySession mocks base method.
func (m *MockIAuthorizationRepository) GetUserIdBySession(arg0 string) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserIdBySession", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIdBySession indicates an expected call of GetUserIdBySession.
func (mr *MockIAuthorizationRepositoryMockRecorder) GetUserIdBySession(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIdBySession", reflect.TypeOf((*MockIAuthorizationRepository)(nil).GetUserIdBySession), arg0)
}