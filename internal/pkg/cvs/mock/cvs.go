// Code generated by MockGen. DO NOT EDIT.
// Source: internal/pkg/cvs/cvs.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	dto "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	gomock "github.com/golang/mock/gomock"
)

// MockICVsRepository is a mock of ICVsRepository interface.
type MockICVsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockICVsRepositoryMockRecorder
}

// MockICVsRepositoryMockRecorder is the mock recorder for MockICVsRepository.
type MockICVsRepositoryMockRecorder struct {
	mock *MockICVsRepository
}

// NewMockICVsRepository creates a new mock instance.
func NewMockICVsRepository(ctrl *gomock.Controller) *MockICVsRepository {
	mock := &MockICVsRepository{ctrl: ctrl}
	mock.recorder = &MockICVsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockICVsRepository) EXPECT() *MockICVsRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockICVsRepository) Create(cv *dto.JSONCv) (*dto.JSONCv, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", cv)
	ret0, _ := ret[0].(*dto.JSONCv)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockICVsRepositoryMockRecorder) Create(cv interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockICVsRepository)(nil).Create), cv)
}

// Delete mocks base method.
func (m *MockICVsRepository) Delete(ID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockICVsRepositoryMockRecorder) Delete(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockICVsRepository)(nil).Delete), ID)
}

// GetByID mocks base method.
func (m *MockICVsRepository) GetByID(ID uint64) (*dto.JSONCv, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ID)
	ret0, _ := ret[0].(*dto.JSONCv)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockICVsRepositoryMockRecorder) GetByID(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockICVsRepository)(nil).GetByID), ID)
}

// GetCVsByApplicantID mocks base method.
func (m *MockICVsRepository) GetCVsByApplicantID(applicantID uint64) ([]*dto.JSONCv, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCVsByApplicantID", applicantID)
	ret0, _ := ret[0].([]*dto.JSONCv)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCVsByApplicantID indicates an expected call of GetCVsByApplicantID.
func (mr *MockICVsRepositoryMockRecorder) GetCVsByApplicantID(applicantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCVsByApplicantID", reflect.TypeOf((*MockICVsRepository)(nil).GetCVsByApplicantID), applicantID)
}

// Update mocks base method.
func (m *MockICVsRepository) Update(ID uint64, updatedCv *dto.JSONCv) (*dto.JSONCv, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ID, updatedCv)
	ret0, _ := ret[0].(*dto.JSONCv)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockICVsRepositoryMockRecorder) Update(ID, updatedCv interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockICVsRepository)(nil).Update), ID, updatedCv)
}

// MockICVsUsecase is a mock of ICVsUsecase interface.
type MockICVsUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockICVsUsecaseMockRecorder
}

// MockICVsUsecaseMockRecorder is the mock recorder for MockICVsUsecase.
type MockICVsUsecaseMockRecorder struct {
	mock *MockICVsUsecase
}

// NewMockICVsUsecase creates a new mock instance.
func NewMockICVsUsecase(ctrl *gomock.Controller) *MockICVsUsecase {
	mock := &MockICVsUsecase{ctrl: ctrl}
	mock.recorder = &MockICVsUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockICVsUsecase) EXPECT() *MockICVsUsecaseMockRecorder {
	return m.recorder
}

// CreateCV mocks base method.
func (m *MockICVsUsecase) CreateCV(cv *dto.JSONCv, currentUser *dto.UserFromSession) (*dto.JSONCv, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCV", cv, currentUser)
	ret0, _ := ret[0].(*dto.JSONCv)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCV indicates an expected call of CreateCV.
func (mr *MockICVsUsecaseMockRecorder) CreateCV(cv, currentUser interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCV", reflect.TypeOf((*MockICVsUsecase)(nil).CreateCV), cv, currentUser)
}

// DeleteCV mocks base method.
func (m *MockICVsUsecase) DeleteCV(ID uint64, currentUser *dto.UserFromSession) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCV", ID, currentUser)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCV indicates an expected call of DeleteCV.
func (mr *MockICVsUsecaseMockRecorder) DeleteCV(ID, currentUser interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCV", reflect.TypeOf((*MockICVsUsecase)(nil).DeleteCV), ID, currentUser)
}

// GetApplicantCVs mocks base method.
func (m *MockICVsUsecase) GetApplicantCVs(applicantID uint64) ([]*dto.JSONGetApplicantCV, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApplicantCVs", applicantID)
	ret0, _ := ret[0].([]*dto.JSONGetApplicantCV)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetApplicantCVs indicates an expected call of GetApplicantCVs.
func (mr *MockICVsUsecaseMockRecorder) GetApplicantCVs(applicantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApplicantCVs", reflect.TypeOf((*MockICVsUsecase)(nil).GetApplicantCVs), applicantID)
}

// GetCV mocks base method.
func (m *MockICVsUsecase) GetCV(ID uint64) (*dto.JSONCv, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCV", ID)
	ret0, _ := ret[0].(*dto.JSONCv)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCV indicates an expected call of GetCV.
func (mr *MockICVsUsecaseMockRecorder) GetCV(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCV", reflect.TypeOf((*MockICVsUsecase)(nil).GetCV), ID)
}

// UpdateCV mocks base method.
func (m *MockICVsUsecase) UpdateCV(ID uint64, currentUser *dto.UserFromSession, cv *dto.JSONCv) (*dto.JSONCv, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCV", ID, currentUser, cv)
	ret0, _ := ret[0].(*dto.JSONCv)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCV indicates an expected call of UpdateCV.
func (mr *MockICVsUsecaseMockRecorder) UpdateCV(ID, currentUser, cv interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCV", reflect.TypeOf((*MockICVsUsecase)(nil).UpdateCV), ID, currentUser, cv)
}
