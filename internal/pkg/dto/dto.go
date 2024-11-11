package dto

import (
	"database/sql"
	"net/http"
)

// Package contains Data Transfer Objects (DTO).
// DTOs used for tasering data from one part of the app to another.

type loggerKey int
type userKey int
type requestIDKey int

// Context keys
const (
	LoggerContextKey    loggerKey    = 1
	UserContextKey      userKey      = 2
	RequestIDContextKey requestIDKey = 3
)

// Error messages
const (
	MsgUnableToGetUserFromContext = "unable to get user from context, please check didn't you forget to add middleware.RequireAuthorization"
	MsgMethodNotAllowed           = "method not allowed"
	MsgInvalidJSON                = "invalid json"
	MsgUnauthorized               = "user unauthorized"
	MsgDataBaseError              = "database error"
	MsgAccessDenied               = "no permissions to perform this action"
	MsgNoCookie                   = "no cookie"
	MsgBadCookie                  = "bad cookie"
	MsgBadUserType                = "got unknown user type"
	MsgNoUserWithSession          = "no user with this session"
	MsgWrongLoginOrPassword       = "wrong login or password"
	MsgUserAlreadyExists          = "user already exists" // TODO: implement check in repository
	MsgUnableToMarshalJSON        = "unable to marshal json"
)

const (
	// UserTypeApplicant is a constant for "applicant" user type
	UserTypeApplicant = "applicant"
	// UserTypeEmployer is a constant for "employer" user type
	UserTypeEmployer = "employer"
)

const SessionIDName = "session_id1"

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// JSONResponse is a standard form of response from backend to frontend
type JSONResponse struct {
	HTTPStatus int         `json:"statusCode"`
	Body       interface{} `json:"body"`
	Error      string      `json:"error"`
}

// JSONUserBody is a struct that used as a field 'Body' in struct JsonResponse
type JSONUserBody struct {
	UserType string `json:"userType"` // use constants UserType
	ID       uint64 `json:"id"`
}

// JSONLoginForm is a struct that receives login's form data from frontend
type JSONLoginForm struct {
	UserType string `json:"userType"` // use constants UserType
	Email    string `json:"email"`
	Password string `json:"password"`
}

// JSONLoutForm accepts user type when someone log outs
type JSONLogoutForm struct {
	UserType string `json:"userType"` // use constants UserType
}

// JSONApplicantRegistrationForm is a struct that receives applicant registration's form data from frontend

// JSONRegistrationForm is a struct that receives employer registration's form data from frontend
type JSONEmployerRegistrationForm struct {
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Position           string `json:"position"`
	Company            string `json:"companyName"`
	CompanyDescription string `json:"companyDescription"`
	CompanyWebsite     string `json:"companyWebsite"`
	Email              string `json:"email"`
	Password           string `json:"password"`
}

// JSONApplicantRegistrationForm is a struct that recipes applicant registration's form data from frontend

type JSONApplicantRegistrationForm struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	BirthDate string `json:"birthDate"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// JSONEmployer is a default representation of employer
type JSONEmployer struct {
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Position           string `json:"position"`
	Company            string `json:"company"`
	CompanyDescription string `json:"companyDescription"`
	CompanyWebsite     string `json:"companyWebsite"`
	Email              string `json:"email"`
}

// JSONEmployer is a default representation of employer
type ApplicantInput struct {
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	CityName            string `json:"cityName"`
	BirthDate           string `json:"birthDate"`
	PathToProfileAvatar string `json:"pathToProfileAvatar"`
	Contacts            string `json:"contacts"`
	Education           string `json:"education"`
	Email               string `json:"email"`
	Password            string `json:"password"`
}

type ApplicantOutput struct {
	UserType            string `json:"userType"`
	ID                  uint64 `json:"id"`
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	CityName            string `json:"cityName"`
	BirthDate           string `json:"birthDate"`
	PathToProfileAvatar string `json:"pathToProfileAvatar"`
	Contacts            string `json:"contacts"`
	Education           string `json:"education"`
	Email               string `json:"email"`
	PasswordHash        string `json:"-"`
	CreatedAt           string `json:"createdAt"`
	UpdatedAt           string `json:"updatedAt"`
}

type ApplicantWithNull struct {
	ID                  uint64
	FirstName           string
	LastName            string
	CityName            sql.NullString
	BirthDate           string
	PathToProfileAvatar string
	Contacts            sql.NullString
	Education           sql.NullString
	Email               string
	PasswordHash        string
	CreatedAt           string
	UpdatedAt           string
}

type EmployerInput struct {
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	CityName            string `json:"cityName"`
	Position            string `json:"position"`
	CompanyName         string `json:"companyName"`
	CompanyDescription  string `json:"companyDescription"`
	CompanyWebsite      string `json:"companyWebsite"`
	PathToProfileAvatar string `json:"pathToProfileAvatar"`
	Contacts            string `json:"contacts"`
	Email               string `json:"email"`
	Password            string `json:"password"`
}

type EmployerWithNull struct {
	ID                  uint64
	FirstName           string
	LastName            string
	CityName            sql.NullString
	Position            string
	CompanyName         string
	CompanyDescription  string
	CompanyWebsite      string
	PathToProfileAvatar string
	Contacts            sql.NullString
	Email               string
	PasswordHash        string
	CreatedAt           string
	UpdatedAt           string
}

type EmployerOutput struct {
	UserType            string `json:"userType"`
	ID                  uint64 `json:"id"`
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	CityName            string `json:"cityName"`
	Position            string `json:"position"`
	CompanyName         string `json:"companyName"`
	CompanyDescription  string `json:"companyDescription"`
	CompanyWebsite      string `json:"companyWebsite"`
	PathToProfileAvatar string `json:"pathToProfileAvatar"`
	Contacts            string `json:"contacts"`
	Email               string `json:"email"`
	PasswordHash        string `json:"-"`
	CreatedAt           string `json:"createdAt"`
	UpdatedAt           string `json:"updatedAt"`
}
type UserIDAndType struct {
	ID       uint64
	UserType string
}

type UserWithSession struct {
	ID        uint64
	UserType  string
	SessionID string
}

type JSONGetEmployerProfile struct {
	ID                 uint64 `json:"id"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	City               string `json:"city"`
	Position           string `json:"position"`
	Company            string `json:"companyName"`
	CompanyDescription string `json:"companyDescription"`
	CompanyWebsite     string `json:"companyWebsite"`
	Contacts           string `json:"contacts"`
	Avatar             string `json:"avatar"`
}

type JSONGetApplicantProfile struct {
	ID        uint64 `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	City      string `json:"city"`
	BirthDate string `json:"birthDate"`
	Avatar    string `json:"avatar"`
	Contacts  string `json:"contacts"`
	Education string `json:"education"`
}

type JSONUpdateEmployerProfile struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	City      string `json:"city"`
	Contacts  string `json:"contacts"`
}

type JSONUpdateApplicantProfile struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	City      string `json:"city"`
	BirthDate string `json:"birthDate"`
	Contacts  string `json:"contacts"`
	Education string `json:"education"`
}

type JSONGetEmployerVacancy struct {
	ID          uint64 `json:"id"`
	EmployerID  uint64 `json:"employer"`
	Salary      int32  `json:"salary"`
	Position    string `json:"position"`
	Location    string `json:"location"`
	Description string `json:"description"`
	WorkType    string `json:"workType"`
	Avatar      string `json:"avatar"`
	CreatedAt   string `json:"createdAt"`
}

type JSONGetApplicantPortfolio struct {
	ID          uint64 `json:"id"`
	ApplicantID uint64 `json:"applicant"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"createdAt"`
}

type JSONGetApplicantCV struct {
	ID                uint64 `json:"id"`
	ApplicantID       uint64 `json:"applicant"`
	PositionRu        string `json:"positionRu"`
	PositionEn        string `json:"positionEn"`
	Description       string `json:"description,omitempty"`
	JobSearchStatus   string `json:"jobSearchStatus"`
	WorkingExperience string `json:"workingExperience"`
	CreatedAt         string `json:"createdAt"`
}

type JSONCv struct {
	ID                  uint64 `json:"id"`
	ApplicantID         uint64 `json:"applicant"`
	PositionRu          string `json:"positionRu"`
	PositionEn          string `json:"positionEn"`
	Description         string `json:"description,omitempty"`
	JobSearchStatusName string `json:"jobSearchStatus"`
	WorkingExperience   string `json:"workingExperience"`
	Avatar              string `json:"avatar"`
	CreatedAt           string `json:"createdAt"`
	UpdatedAt           string `json:"updatedAt"`
}

type JSONVacancy struct {
	ID          uint64 `json:"id"`
	EmployerID  uint64 `json:"employer"`
	Salary      int32  `json:"salary"`
	Position    string `json:"position"`
	Location    string `json:"location"`
	Description string `json:"description"`
	WorkType    string `json:"workType"`
	Avatar      string `json:"avatar"`
	CompanyName string `json:"companyName"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type JSONVacancySubscriptionStatus struct {
	ID           uint64 `json:"id"`
	ApplicantID  uint64 `json:"applicantID"`
	IsSubscribed bool   `json:"isSubscribed"`
}

type JSONVacancySubscribers struct {
	ID          uint64                     `json:"vacancyID"`
	Subscribers []*JSONGetApplicantProfile `json:"subscribers"`
}

type SessionUser struct {
	ID       uint64
	UserType string
}
