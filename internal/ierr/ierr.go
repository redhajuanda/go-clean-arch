package ierr

import "fmt"

type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

var (
	ErrInternal         = Error{Code: "500000", Message: "we encountered an error while processing your request (internal server error)"}
	ErrResourceNotFound = Error{Code: "404000", Message: "the requested resource was not found"}
	ErrBadRequest       = Error{Code: "400000", Message: "your request is in a bad format"}
	ErrUnauthorized     = Error{Code: "401000", Message: "you are not authorized to perform the requested action"}
	ErrForbidden        = Error{Code: "403000", Message: "you don't have access to this resource"}
)

var (
	ErrUserAlreadyRegistered = Error{Code: "400020", Message: "you're already registered"}
	ErrInvalidCreds          = Error{Code: "400021", Message: "invalid username or password"}
	ErrUserIsNotActive       = Error{Code: "400022", Message: "user is not active"}
	ErrPinAlreadySet         = Error{Code: "400023", Message: "pin already set"}
	ErrPinIsNotSet           = Error{Code: "400024", Message: "pin is not set"}
	ErrWrongOTP              = Error{Code: "400025", Message: "wrong OTP code"}
	ErrExpiredOTP            = Error{Code: "400026", Message: "expired OTP code"}
	ErrInvalidToken          = Error{Code: "400027", Message: "token is invalid"}
	ErrExpiredToken          = Error{Code: "400028", Message: "token has expired"}
	ErrEmailAlreadyVerified  = Error{Code: "400029", Message: "email has been verified"}
	ErrInvalidPhoneNumber    = Error{Code: "400030", Message: "phone number is invalid"}
)
