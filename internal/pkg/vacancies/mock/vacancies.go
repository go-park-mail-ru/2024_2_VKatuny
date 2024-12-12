// Code generated by MockGen. DO NOT EDIT.
// Source: internal/pkg/vacancies/vacancies.go
//
// Generated by this command:
//
//	mockgen -source=internal/pkg/vacancies/vacancies.go -destination=internal/pkg/vacancies/mock/vacancies.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	dto "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/dto"
	models "github.com/go-park-mail-ru/2024_2_VKatuny/internal/pkg/models"
	gomock "go.uber.org/mock/gomock"
)

// MockIVacanciesRepository is a mock of IVacanciesRepository interface.
type MockIVacanciesRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIVacanciesRepositoryMockRecorder
	isgomock struct{}
}

// MockIVacanciesRepositoryMockRecorder is the mock recorder for MockIVacanciesRepository.
type MockIVacanciesRepositoryMockRecorder struct {
	mock *MockIVacanciesRepository
}

// NewMockIVacanciesRepository creates a new mock instance.
func NewMockIVacanciesRepository(ctrl *gomock.Controller) *MockIVacanciesRepository {
	mock := &MockIVacanciesRepository{ctrl: ctrl}
	mock.recorder = &MockIVacanciesRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIVacanciesRepository) EXPECT() *MockIVacanciesRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIVacanciesRepository) Create(vacancy *dto.JSONVacancy) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", vacancy)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIVacanciesRepositoryMockRecorder) Create(vacancy any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIVacanciesRepository)(nil).Create), vacancy)
}

// Delete mocks base method.
func (m *MockIVacanciesRepository) Delete(ID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockIVacanciesRepositoryMockRecorder) Delete(ID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIVacanciesRepository)(nil).Delete), ID)
}

// GetApplicantFavoriteVacancies mocks base method.
func (m *MockIVacanciesRepository) GetApplicantFavoriteVacancies(applicantID uint64) ([]*dto.JSONVacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApplicantFavoriteVacancies", applicantID)
	ret0, _ := ret[0].([]*dto.JSONVacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetApplicantFavoriteVacancies indicates an expected call of GetApplicantFavoriteVacancies.
func (mr *MockIVacanciesRepositoryMockRecorder) GetApplicantFavoriteVacancies(applicantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApplicantFavoriteVacancies", reflect.TypeOf((*MockIVacanciesRepository)(nil).GetApplicantFavoriteVacancies), applicantID)
}

// GetByID mocks base method.
func (m *MockIVacanciesRepository) GetByID(ID uint64) (*dto.JSONVacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ID)
	ret0, _ := ret[0].(*dto.JSONVacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIVacanciesRepositoryMockRecorder) GetByID(ID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIVacanciesRepository)(nil).GetByID), ID)
}

// GetSubscribersCount mocks base method.
func (m *MockIVacanciesRepository) GetSubscribersCount(ID uint64) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscribersCount", ID)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscribersCount indicates an expected call of GetSubscribersCount.
func (mr *MockIVacanciesRepositoryMockRecorder) GetSubscribersCount(ID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscribersCount", reflect.TypeOf((*MockIVacanciesRepository)(nil).GetSubscribersCount), ID)
}

// GetSubscribersList mocks base method.
func (m *MockIVacanciesRepository) GetSubscribersList(ID uint64) ([]*models.Applicant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscribersList", ID)
	ret0, _ := ret[0].([]*models.Applicant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscribersList indicates an expected call of GetSubscribersList.
func (mr *MockIVacanciesRepositoryMockRecorder) GetSubscribersList(ID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscribersList", reflect.TypeOf((*MockIVacanciesRepository)(nil).GetSubscribersList), ID)
}

// GetSubscriptionStatus mocks base method.
func (m *MockIVacanciesRepository) GetSubscriptionStatus(ID, applicantID uint64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscriptionStatus", ID, applicantID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscriptionStatus indicates an expected call of GetSubscriptionStatus.
func (mr *MockIVacanciesRepositoryMockRecorder) GetSubscriptionStatus(ID, applicantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscriptionStatus", reflect.TypeOf((*MockIVacanciesRepository)(nil).GetSubscriptionStatus), ID, applicantID)
}

// GetVacanciesByEmployerID mocks base method.
func (m *MockIVacanciesRepository) GetVacanciesByEmployerID(employerID uint64) ([]*dto.JSONVacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVacanciesByEmployerID", employerID)
	ret0, _ := ret[0].([]*dto.JSONVacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVacanciesByEmployerID indicates an expected call of GetVacanciesByEmployerID.
func (mr *MockIVacanciesRepositoryMockRecorder) GetVacanciesByEmployerID(employerID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVacanciesByEmployerID", reflect.TypeOf((*MockIVacanciesRepository)(nil).GetVacanciesByEmployerID), employerID)
}

// MakeFavorite mocks base method.
func (m *MockIVacanciesRepository) MakeFavorite(ID, applicantID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeFavorite", ID, applicantID)
	ret0, _ := ret[0].(error)
	return ret0
}

// MakeFavorite indicates an expected call of MakeFavorite.
func (mr *MockIVacanciesRepositoryMockRecorder) MakeFavorite(ID, applicantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeFavorite", reflect.TypeOf((*MockIVacanciesRepository)(nil).MakeFavorite), ID, applicantID)
}

// SearchAll mocks base method.
func (m *MockIVacanciesRepository) SearchAll(offset, num uint64, searchStr, group, searchBy string) ([]*dto.JSONVacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchAll", offset, num, searchStr, group, searchBy)
	ret0, _ := ret[0].([]*dto.JSONVacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchAll indicates an expected call of SearchAll.
func (mr *MockIVacanciesRepositoryMockRecorder) SearchAll(offset, num, searchStr, group, searchBy any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchAll", reflect.TypeOf((*MockIVacanciesRepository)(nil).SearchAll), offset, num, searchStr, group, searchBy)
}

// Subscribe mocks base method.
func (m *MockIVacanciesRepository) Subscribe(ID, applicantID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", ID, applicantID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockIVacanciesRepositoryMockRecorder) Subscribe(ID, applicantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockIVacanciesRepository)(nil).Subscribe), ID, applicantID)
}

// Unsubscribe mocks base method.
func (m *MockIVacanciesRepository) Unsubscribe(ID, applicantID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unsubscribe", ID, applicantID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unsubscribe indicates an expected call of Unsubscribe.
func (mr *MockIVacanciesRepositoryMockRecorder) Unsubscribe(ID, applicantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsubscribe", reflect.TypeOf((*MockIVacanciesRepository)(nil).Unsubscribe), ID, applicantID)
}

// Update mocks base method.
func (m *MockIVacanciesRepository) Update(ID uint64, updatedVacancy *dto.JSONVacancy) (*dto.JSONVacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ID, updatedVacancy)
	ret0, _ := ret[0].(*dto.JSONVacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockIVacanciesRepositoryMockRecorder) Update(ID, updatedVacancy any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIVacanciesRepository)(nil).Update), ID, updatedVacancy)
}

// MockIVacanciesUsecase is a mock of IVacanciesUsecase interface.
type MockIVacanciesUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockIVacanciesUsecaseMockRecorder
	isgomock struct{}
}

// MockIVacanciesUsecaseMockRecorder is the mock recorder for MockIVacanciesUsecase.
type MockIVacanciesUsecaseMockRecorder struct {
	mock *MockIVacanciesUsecase
}

// NewMockIVacanciesUsecase creates a new mock instance.
func NewMockIVacanciesUsecase(ctrl *gomock.Controller) *MockIVacanciesUsecase {
	mock := &MockIVacanciesUsecase{ctrl: ctrl}
	mock.recorder = &MockIVacanciesUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIVacanciesUsecase) EXPECT() *MockIVacanciesUsecaseMockRecorder {
	return m.recorder
}

// AddIntoFavorite mocks base method.
func (m *MockIVacanciesUsecase) AddIntoFavorite(ID uint64, currentUser *dto.UserFromSession) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddIntoFavorite", ID, currentUser)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddIntoFavorite indicates an expected call of AddIntoFavorite.
func (mr *MockIVacanciesUsecaseMockRecorder) AddIntoFavorite(ID, currentUser any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddIntoFavorite", reflect.TypeOf((*MockIVacanciesUsecase)(nil).AddIntoFavorite), ID, currentUser)
}

// CreateVacancy mocks base method.
func (m *MockIVacanciesUsecase) CreateVacancy(vacancy *dto.JSONVacancy, currentUser *dto.UserFromSession) (*dto.JSONVacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateVacancy", vacancy, currentUser)
	ret0, _ := ret[0].(*dto.JSONVacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateVacancy indicates an expected call of CreateVacancy.
func (mr *MockIVacanciesUsecaseMockRecorder) CreateVacancy(vacancy, currentUser any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVacancy", reflect.TypeOf((*MockIVacanciesUsecase)(nil).CreateVacancy), vacancy, currentUser)
}

// DeleteVacancy mocks base method.
func (m *MockIVacanciesUsecase) DeleteVacancy(ID uint64, currentUser *dto.UserFromSession) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteVacancy", ID, currentUser)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteVacancy indicates an expected call of DeleteVacancy.
func (mr *MockIVacanciesUsecaseMockRecorder) DeleteVacancy(ID, currentUser any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVacancy", reflect.TypeOf((*MockIVacanciesUsecase)(nil).DeleteVacancy), ID, currentUser)
}

// GetApplicantFavoriteVacancies mocks base method.
func (m *MockIVacanciesUsecase) GetApplicantFavoriteVacancies(applicantID uint64) ([]*dto.JSONGetEmployerVacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApplicantFavoriteVacancies", applicantID)
	ret0, _ := ret[0].([]*dto.JSONGetEmployerVacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetApplicantFavoriteVacancies indicates an expected call of GetApplicantFavoriteVacancies.
func (mr *MockIVacanciesUsecaseMockRecorder) GetApplicantFavoriteVacancies(applicantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApplicantFavoriteVacancies", reflect.TypeOf((*MockIVacanciesUsecase)(nil).GetApplicantFavoriteVacancies), applicantID)
}

// GetSubscriptionInfo mocks base method.
func (m *MockIVacanciesUsecase) GetSubscriptionInfo(ID, applicantID uint64) (*dto.JSONVacancySubscriptionStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscriptionInfo", ID, applicantID)
	ret0, _ := ret[0].(*dto.JSONVacancySubscriptionStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscriptionInfo indicates an expected call of GetSubscriptionInfo.
func (mr *MockIVacanciesUsecaseMockRecorder) GetSubscriptionInfo(ID, applicantID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscriptionInfo", reflect.TypeOf((*MockIVacanciesUsecase)(nil).GetSubscriptionInfo), ID, applicantID)
}

// GetVacanciesByEmployerID mocks base method.
func (m *MockIVacanciesUsecase) GetVacanciesByEmployerID(employerID uint64) ([]*dto.JSONGetEmployerVacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVacanciesByEmployerID", employerID)
	ret0, _ := ret[0].([]*dto.JSONGetEmployerVacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVacanciesByEmployerID indicates an expected call of GetVacanciesByEmployerID.
func (mr *MockIVacanciesUsecaseMockRecorder) GetVacanciesByEmployerID(employerID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVacanciesByEmployerID", reflect.TypeOf((*MockIVacanciesUsecase)(nil).GetVacanciesByEmployerID), employerID)
}

// GetVacancy mocks base method.
func (m *MockIVacanciesUsecase) GetVacancy(ID uint64) (*dto.JSONVacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVacancy", ID)
	ret0, _ := ret[0].(*dto.JSONVacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVacancy indicates an expected call of GetVacancy.
func (mr *MockIVacanciesUsecaseMockRecorder) GetVacancy(ID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVacancy", reflect.TypeOf((*MockIVacanciesUsecase)(nil).GetVacancy), ID)
}

// GetVacancySubscribers mocks base method.
func (m *MockIVacanciesUsecase) GetVacancySubscribers(ID uint64, currentUser *dto.UserFromSession) (*dto.JSONVacancySubscribers, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVacancySubscribers", ID, currentUser)
	ret0, _ := ret[0].(*dto.JSONVacancySubscribers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVacancySubscribers indicates an expected call of GetVacancySubscribers.
func (mr *MockIVacanciesUsecaseMockRecorder) GetVacancySubscribers(ID, currentUser any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVacancySubscribers", reflect.TypeOf((*MockIVacanciesUsecase)(nil).GetVacancySubscribers), ID, currentUser)
}

// SearchVacancies mocks base method.
func (m *MockIVacanciesUsecase) SearchVacancies(offsetStr, numStr, searchStr, group, searchBy string) ([]*dto.JSONVacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchVacancies", offsetStr, numStr, searchStr, group, searchBy)
	ret0, _ := ret[0].([]*dto.JSONVacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchVacancies indicates an expected call of SearchVacancies.
func (mr *MockIVacanciesUsecaseMockRecorder) SearchVacancies(offsetStr, numStr, searchStr, group, searchBy any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchVacancies", reflect.TypeOf((*MockIVacanciesUsecase)(nil).SearchVacancies), offsetStr, numStr, searchStr, group, searchBy)
}

// SubscribeOnVacancy mocks base method.
func (m *MockIVacanciesUsecase) SubscribeOnVacancy(ID uint64, currentUser *dto.UserFromSession) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeOnVacancy", ID, currentUser)
	ret0, _ := ret[0].(error)
	return ret0
}

// SubscribeOnVacancy indicates an expected call of SubscribeOnVacancy.
func (mr *MockIVacanciesUsecaseMockRecorder) SubscribeOnVacancy(ID, currentUser any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeOnVacancy", reflect.TypeOf((*MockIVacanciesUsecase)(nil).SubscribeOnVacancy), ID, currentUser)
}

// UnsubscribeFromVacancy mocks base method.
func (m *MockIVacanciesUsecase) UnsubscribeFromVacancy(ID uint64, currentUser *dto.UserFromSession) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnsubscribeFromVacancy", ID, currentUser)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnsubscribeFromVacancy indicates an expected call of UnsubscribeFromVacancy.
func (mr *MockIVacanciesUsecaseMockRecorder) UnsubscribeFromVacancy(ID, currentUser any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnsubscribeFromVacancy", reflect.TypeOf((*MockIVacanciesUsecase)(nil).UnsubscribeFromVacancy), ID, currentUser)
}

// UpdateVacancy mocks base method.
func (m *MockIVacanciesUsecase) UpdateVacancy(ID uint64, updatedVacancy *dto.JSONVacancy, currentUser *dto.UserFromSession) (*dto.JSONVacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateVacancy", ID, updatedVacancy, currentUser)
	ret0, _ := ret[0].(*dto.JSONVacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateVacancy indicates an expected call of UpdateVacancy.
func (mr *MockIVacanciesUsecaseMockRecorder) UpdateVacancy(ID, updatedVacancy, currentUser any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateVacancy", reflect.TypeOf((*MockIVacanciesUsecase)(nil).UpdateVacancy), ID, updatedVacancy, currentUser)
}

// ValidateQueryParameters mocks base method.
func (m *MockIVacanciesUsecase) ValidateQueryParameters(offset, num string) (uint64, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateQueryParameters", offset, num)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(uint64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ValidateQueryParameters indicates an expected call of ValidateQueryParameters.
func (mr *MockIVacanciesUsecaseMockRecorder) ValidateQueryParameters(offset, num any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateQueryParameters", reflect.TypeOf((*MockIVacanciesUsecase)(nil).ValidateQueryParameters), offset, num)
}
