package app

import (
	"go-clean-arch/config"
	"go-clean-arch/infra/httplog"
	"go-clean-arch/infra/logger"
	"go-clean-arch/internal/auth"
	"go-clean-arch/internal/repository/postgres"
	"go-clean-arch/internal/user"
	"net/http"

	customMiddleware "go-clean-arch/middleware"

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// App struct
type App struct {
	cfg    *config.Config
	router *echo.Echo
	db     *pg.DB
	log    logger.Logger
}

// New inits a new app
func New(
	cfg *config.Config,
	router *echo.Echo,
	db *pg.DB,
	log logger.Logger,
) *App {
	return &App{
		cfg,
		router,
		db,
		log,
	}
}

// BuildHandler builds handlers and returns router
func (app App) BuildHandler() *echo.Echo {

	app.configRouter()

	// Endpoint for swagger documentations
	app.router.GET("/swagger/*", echoSwagger.WrapHandler)

	repoRegistry := postgres.NewRepositoryRegistry(app.db)

	auth.RegisterModule(
		*app.router.Group(""),
		app.cfg,
		auth.NewService(
			app.cfg,
			repoRegistry,
		),
	)

	user.RegisterModule(
		*app.router.Group(""),
		app.cfg,
		user.NewService(
			app.cfg,
			repoRegistry,
		),
	)

	app.router.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"app": app.cfg.Server.NAME,
			"env": app.cfg.Server.ENV.String(),
		})
	})

	app.router.Any("", func(c echo.Context) error {
		return echo.NotFoundHandler(c)
	})

	app.router.Any("/*", func(c echo.Context) error {
		return echo.NotFoundHandler(c)
	})

	return app.router
}

func (app App) configRouter() {

	app.router.Pre(middleware.RemoveTrailingSlash())
	// app.router.Use(middleware.RequestID())
	app.router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	app.router.Use(customMiddleware.RequestIDContext()) // middleware for insert request id into context

	// Setup custom HTTP error handler
	app.router.HTTPErrorHandler = CustomHTTPErrorHandler(app.cfg, app.log, httplog.NewHTTPLog(app.db))

	// Register middleware recover from panic

	// Setup access log
	app.router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","request_id":"${id}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status},"error":"${error}","latency":"${latency}",` +
			`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
			`"service":"` + app.cfg.Server.NAME + `",` +
			`"environment":"` + app.cfg.Server.ENV.String() + `",` +
			`"type":"access",` +
			`"bytes_out":${bytes_out}}` + "\n",
	}))

	app.router.Use(customMiddleware.Recover(app.log))

}
