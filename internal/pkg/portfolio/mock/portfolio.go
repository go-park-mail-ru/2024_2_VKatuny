// Code generated by MockGen. DO NOT EDIT.
// Source: internal/pkg/portfolio/portfolio.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	dto "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	models "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
	gomock "github.com/golang/mock/gomock"
)

// MockIPortfolioRepository is a mock of IPortfolioRepository interface.
type MockIPortfolioRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIPortfolioRepositoryMockRecorder
}

// MockIPortfolioRepositoryMockRecorder is the mock recorder for MockIPortfolioRepository.
type MockIPortfolioRepositoryMockRecorder struct {
	mock *MockIPortfolioRepository
}

// NewMockIPortfolioRepository creates a new mock instance.
func NewMockIPortfolioRepository(ctrl *gomock.Controller) *MockIPortfolioRepository {
	mock := &MockIPortfolioRepository{ctrl: ctrl}
	mock.recorder = &MockIPortfolioRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPortfolioRepository) EXPECT() *MockIPortfolioRepositoryMockRecorder {
	return m.recorder
}

// GetPortfoliosByApplicantID mocks base method.
func (m *MockIPortfolioRepository) GetPortfoliosByApplicantID(applicantID uint64) ([]*models.Portfolio, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPortfoliosByApplicantID", applicantID)
	ret0, _ := ret[0].([]*models.Portfolio)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPortfoliosByApplicantID indicates an expected call of GetPortfoliosByApplicantID.
func (mr *MockIPortfolioRepositoryMockRecorder) GetPortfoliosByApplicantID(applicantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPortfoliosByApplicantID", reflect.TypeOf((*MockIPortfolioRepository)(nil).GetPortfoliosByApplicantID), applicantID)
}

// MockIPortfolioUsecase is a mock of IPortfolioUsecase interface.
type MockIPortfolioUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockIPortfolioUsecaseMockRecorder
}

// MockIPortfolioUsecaseMockRecorder is the mock recorder for MockIPortfolioUsecase.
type MockIPortfolioUsecaseMockRecorder struct {
	mock *MockIPortfolioUsecase
}

// NewMockIPortfolioUsecase creates a new mock instance.
func NewMockIPortfolioUsecase(ctrl *gomock.Controller) *MockIPortfolioUsecase {
	mock := &MockIPortfolioUsecase{ctrl: ctrl}
	mock.recorder = &MockIPortfolioUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPortfolioUsecase) EXPECT() *MockIPortfolioUsecaseMockRecorder {
	return m.recorder
}

// GetApplicantPortfolios mocks base method.
func (m *MockIPortfolioUsecase) GetApplicantPortfolios(applicantID uint64) ([]*dto.JSONGetApplicantPortfolio, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApplicantPortfolios", applicantID)
	ret0, _ := ret[0].([]*dto.JSONGetApplicantPortfolio)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetApplicantPortfolios indicates an expected call of GetApplicantPortfolios.
func (mr *MockIPortfolioUsecaseMockRecorder) GetApplicantPortfolios(applicantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApplicantPortfolios", reflect.TypeOf((*MockIPortfolioUsecase)(nil).GetApplicantPortfolios), applicantID)
}
