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
	ErrInternal         = Error{Code: "00000", Message: "we encountered an error while processing your request (internal server error)"}
	ErrResourceNotFound = Error{Code: "00001", Message: "the requested resource was not found"}
	ErrBadRequest       = Error{Code: "00002", Message: "your request is in a bad format"}
	ErrUnauthorized     = Error{Code: "00003", Message: "you are not authorized to perform the requested action"}
	ErrForbidden        = Error{Code: "00004", Message: "you don't have access to this resource"}
)

var (
	ErrNoRowsAffected = Error{Code: "00006", Message: "no rows affected"}
)

var (
	ErrUserAlreadyRegistered = Error{Code: "10020", Message: "you're already registered"}
	ErrInvalidCreds          = Error{Code: "10021", Message: "invalid username or password"}
	ErrUserIsNotActive       = Error{Code: "10022", Message: "user is not active"}
	ErrPinAlreadySet         = Error{Code: "10020", Message: "pin already set"}
	ErrPinIsNotSet           = Error{Code: "10021", Message: "pin is not set"}
	ErrWrongOTP              = Error{Code: "10022", Message: "wrong OTP code"}
	ErrExpiredOTP            = Error{Code: "10023", Message: "expired OTP code"}
	ErrInvalidToken          = Error{Code: "10024", Message: "token is invalid"}
	ErrExpiredToken          = Error{Code: "10025", Message: "token has expired"}
	ErrEmailAlreadyVerified  = Error{Code: "10026", Message: "email has been verified"}
	ErrInvalidPhoneNumber    = Error{Code: "10027", Message: "phone number is invalid"}

	ErrSimulationNotFound             = Error{Code: "20000", Message: "simulation not found"}
	ErrTxNotFound                     = Error{Code: "20001", Message: "transaction not found"}
	ErrTxCannotBeCancel               = Error{Code: "20002", Message: "transaction cannot be cancel"}
	ErrPaymentCallbackNotFound        = Error{Code: "20003", Message: "payment callback not found"}
	ErrPlafondDoesnotMeet             = Error{Code: "20004", Message: "plafond does not meet"}
	ErrRefIDAlreadyUsed               = Error{Code: "20005", Message: "you've used the same ref_id for another transaction"}
	ErrDuplicateVoucherWithNoReversal = Error{Code: "20006", Message: "duplicate voucher with no reversal"}
	ErrAlreadyReversal                = Error{Code: "20007", Message: "this voucher is already did reversal"}

	ErrCustomerNotFound    = Error{Code: "30000", Message: "customer not found"}
	ErrCustomerIsSuspended = Error{Code: "30001", Message: "customer is suspended"}

	ErrProductNotFound = Error{Code: "40000", Message: "product not found"}
)
