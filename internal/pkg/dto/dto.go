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
	UserTypeEmployer  = "employer"
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
	Email    string `json:"login"`
	Password string `json:"password"`
}

// JSONLoutForm accepts user type when somene log outs
type JSONLogoutForm struct {
	UserType string `json:"userType"` // use constants UserType
}

// JSONRegistrationForm is a struct that recives employer registration's form data from frontend 
type JSONEmployerRegistrationForm struct {
	Name        string `json:"firstName"`
	LastName    string `json:"lastName"`
	Position    string `json:"position"`
	CompanyName string `json:"companyName"`
	Description string `json:"companyDescription"`
	Website     string `json:"website"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

// JSONApplicantRegistrationForm is a struct that recives applicant registration's form data from frontend
type JSONApplicantRegistrationForm struct {
	Name      string `json:"firstName"`
	LastName  string `json:"lastName"`
	BirthDate string `json:"birthDate"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
