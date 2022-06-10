package user

import (
	"go-clean-arch/configs"
	"go-clean-arch/middleware"
	"go-clean-arch/shared/response"

	"github.com/labstack/echo/v4"
)

// RegisterAPI registers a new user api
func RegisterAPI(r echo.Group, cfg *configs.Config, service IService) {
	handler := handler{cfg, service}

	// Private endpoint
	r.Use(middleware.MustLoggedIn(cfg.JWT.SigningKey))

	r.GET("/me", handler.get)
}

type handler struct {
	cfg     *configs.Config
	service IService
}

// get godoc
// @Router /me [get]
// @Tags User
// @Summary Get me
// @Description Get me
// @Accept json
// @Produce json
// @Security BearerToken
// @Success 200 {object} response.Response{data=domain.User} "Success"
// @failure 500 {object} response.ErrorResponse500
func (h handler) get(c echo.Context) error {
	ctx := c.Request().Context()
	user, err := h.service.Get(ctx)
	if err != nil {
		return err
	}
	return response.SuccessOK(c, user)
}
