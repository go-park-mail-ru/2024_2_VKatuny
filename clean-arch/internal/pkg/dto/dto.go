package dto

// "github.com/sirupsen/logrus"

type loggerKey int

const LoggerContextKey loggerKey = 1

const (
	UserTypeApplicant = "applicant"
	UserTypeEmployer  = "employer"
)

// Standart json response from backend to frontend
type JsonResponse struct {
	HttpStatus int         `json:"statusCode"`
	Body       interface{} `json:"body"`
	Error      string      `json:"error"`
}

// use this struct as a field 'Body' in struct JsonResponse
type JsonUserBody struct {
	UserType string `json:"userType"` // use constants UserType
	ID       uint64 `json:"id"`
}

type JsonLoginForm struct {
	UserType string `json:"userType"` // use constants UserType
	Email    string `json:"login"`
	Password string `json:"password"`
}
