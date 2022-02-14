package user

import (
	"go-clean-arch/config"
	"go-clean-arch/middleware"
	"go-clean-arch/shared/response"

	"github.com/labstack/echo/v4"
)

// RegisterModule registers a new update request module
func RegisterModule(r echo.Group, cfg *config.Config, service IService) {
	handler := handler{cfg, service}

	// Private endpoint
	r.Use(middleware.MustLoggedIn(cfg.JWT.SigningKeyMitra))

	r.GET("/me", handler.get)
}

type handler struct {
	cfg     *config.Config
	service IService
}

// get godoc
// @Router /me [get]
// @Tags Mitra
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
