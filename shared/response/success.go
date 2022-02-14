package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response struct
type Response struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"success"`
	Data    interface{} `json:"data"`
}

// Success responses with JSON formatresponseMsg
func Success(c echo.Context, code int, data interface{}, msg ...string) error {

	responseMsg := buildResponseMsg("Success", msg...)

	if data == nil {
		data = map[string]interface{}{}
	}

	res := Response{
		Success: true,
		Message: responseMsg,
		Data:    data,
	}
	return c.JSON(code, res)
}

// SuccessOK returns code 200
func SuccessOK(c echo.Context, data interface{}, msg ...string) error {
	return Success(c, http.StatusOK, data, msg...)
}

// SuccessCreated returns code 201
func SuccessCreated(c echo.Context, data interface{}, msg ...string) error {
	return Success(c, http.StatusCreated, data, msg...)
}
