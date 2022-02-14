package response

// ErrorResponse400 example for swagger doc
type ErrorResponse400 struct {
	Success   bool   `json:"success" example:"false"`
	Message   string `json:"message" example:"your request is in a bad format"`
	ErrorCode string `json:"error_code,omitempty" example:"00002"`
} //@name Bad Request

// ErrorResponse401 example for swagger doc
type ErrorResponse401 struct {
	Success   bool   `json:"success" example:"false"`
	Message   string `json:"message" example:"you are not authorized to perform the requested action"`
	ErrorCode string `json:"error_code,omitempty" example:"00003"`
} //@name Unauthorized

// ErrorResponse403 example for swagger doc
type ErrorResponse403 struct {
	Success   bool   `json:"success" example:"false"`
	Message   string `json:"message" example:"you don't have access to this resource"`
	ErrorCode string `json:"error_code,omitempty" example:"00004"`
} //@name Forbidden

// ErrorResponseWrongOTPCode example for swagger doc
type ErrorResponseWrongOTPCode struct {
	Success   bool   `json:"success" example:"false"`
	Message   string `json:"message" example:"wrong OTP code"`
	ErrorCode string `json:"error_code,omitempty" example:"10022"`
} //@name WrongOTPCode

// ErrorResponse404 example for swagger doc
type ErrorResponse404 struct {
	Success   bool   `json:"success" example:"false"`
	Message   string `json:"message" example:"the requested resource was not found"`
	ErrorCode string `json:"error_code,omitempty" example:"00001"`
} //@name Not Found

// ErrorResponse500 example for swagger doc
type ErrorResponse500 struct {
	Success   bool   `json:"success" example:"false"`
	Message   string `json:"message" example:"we encountered an error while processing your request (internal server error)"`
	ErrorCode string `json:"error_code,omitempty" example:"00000"`
} //@name Internal Server Error
