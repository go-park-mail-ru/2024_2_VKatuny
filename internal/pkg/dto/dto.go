package dto

// Package contains Data Transfer Objects (DTO).
// DTOs used for trnsfering data from one part of the app to another.

type loggerKey int

// LoggerContextKey is a key for logger
const LoggerContextKey loggerKey = 1

const (
	// UserTypeApplicant is a constant for "applicant" user type
	UserTypeApplicant = "applicant"
	// UserTypeEmployer is a constant for "employer" user type
	UserTypeEmployer = "employer"
)

// JSONResponse is a standart form of response from backend to frontend
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

// JSONLoginForm is a struct that recives login's form data from frontend
type JSONLoginForm struct {
	UserType string `json:"userType"` // use constants UserType
	Email    string `json:"email"`
	Password string `json:"password"`
}

// JSONLoutForm accepts user type when somene log outs
type JSONLogoutForm struct {
	UserType string `json:"userType"` // use constants UserType
}

// JSONRegistrationForm is a struct that recives employer registration's form data from frontend
type JSONEmployerRegistrationForm struct {
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Position           string `json:"position"`
	Company            string `json:"Name"`
	CompanyDescription string `json:"companyDescription"`
	CompanyWebsite     string `json:"comapnyWebsite"`
	Email              string `json:"email"`
	Password           string `json:"password"`
}

// JSONApplicantRegistrationForm is a struct that recives applicant registration's form data from frontend
type JSONApplicantRegistrationForm struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	BirthDate string `json:"birthDate"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// JSONEmployer is a default represenation of employer
type JSONEmployer struct {
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Position           string `json:"position"`
	Company            string `json:"company"`
	CompanyDescription string `json:"companyDescription"`
	CompanyWebsite     string `json:"companyWebsite"`
	Email              string `json:"email"`
}
