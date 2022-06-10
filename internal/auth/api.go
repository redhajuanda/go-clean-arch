package auth

import (
	"go-clean-arch/configs"
	"go-clean-arch/internal/ierr"
	"go-clean-arch/shared/response"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// RegisterAPI registers a new auth module
func RegisterAPI(r echo.Group, cfg *configs.Config, service IService) {
	handler := handler{cfg, service}

	r.POST("/auth/login", handler.login)
	r.POST("/auth/token/refresh", handler.refreshToken)
}

type handler struct {
	cfg     *configs.Config
	service IService
}

// login godoc
// @Router /auth/login [post]
// @Tags Auth
// @Summary Login
// @Description Login
// @Accept json
// @Produce json
// @Param payload body RequestLogin false " "
// @Success 200 {object} response.Response{data=LoginResponse} "Success"
// @failure 400 {object} response.ErrorResponse400
// @failure 500 {object} response.ErrorResponse500
func (h handler) login(c echo.Context) error {

	ctx := c.Request().Context()
	req := RequestLogin{}
	err := c.Bind(&req)
	if err != nil {
		return response.ErrBadRequest(err)
	}

	resp, err := h.service.Login(ctx, req)
	if err != nil {
		switch errors.Cause(err) {
		case ierr.ErrUserAlreadyRegistered, ierr.ErrUserIsNotActive:
			return response.ErrBadRequest(err)
		case ierr.ErrInvalidCreds:
			return response.ErrUnauthorized(err)
		}
		return err
	}

	return response.SuccessOK(c, resp, "user authenticated")
}

// refreshToken godoc
// @Router /auth/token/refresh [post]
// @Tags Auth
// @Summary Refresh access token
// @Description Refresh access token
// @Accept json
// @Produce json
// @Param payload body RefreshTokenRequest false " "
// @Success 200 {object} response.Response{data=LoginResponse} "Refresh token success"
// @failure 400 {object} response.ErrorResponse400
// @failure 403 {object} response.ErrorResponse403
// @failure 500 {object} response.ErrorResponse500
func (h handler) refreshToken(c echo.Context) error {
	var req RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrBadRequest(err)
	}

	res, err := h.service.RefreshToken(c.Request().Context(), req)
	if err != nil {
		switch errors.Cause(err) {
		case ierr.ErrInvalidToken:
			return response.ErrBadRequest(err)
		case ierr.ErrExpiredToken:
			return response.ErrForbidden(err)
		}
		return err
	}

	return response.SuccessOK(c, res, "token refreshed")
}
