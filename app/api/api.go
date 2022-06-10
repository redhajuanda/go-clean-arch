package api

import (
	"context"
	"fmt"
	"go-clean-arch/app"
	"go-clean-arch/configs"
	"go-clean-arch/docs"
	"go-clean-arch/internal/auth"
	"go-clean-arch/internal/repository/postgres"
	"go-clean-arch/internal/user"
	"go-clean-arch/pkg/db"
	"go-clean-arch/pkg/httplog"
	"go-clean-arch/pkg/logger"
	"go-clean-arch/pkg/otel"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	customMiddleware "go-clean-arch/middleware"

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// API struct
type API struct {
	cfg    *configs.Config
	router *echo.Echo
	db     *pg.DB
	log    logger.Logger
}

// New inits a new api
func New() *API {

	cfg := configs.LoadDefault()

	// set swagger config
	docs.SwaggerInfo.Host = cfg.Server.BASEURL
	docs.SwaggerInfo.Version = app.Version

	log := logger.New(cfg.Server.NAME, app.Version)
	logger.SetFormatter(&logrus.JSONFormatter{})

	db := db.NewGoPG(cfg)
	router := echo.New()

	return &API{
		cfg,
		router,
		db,
		log,
	}
}

// BuildHandler builds handlers and returns router
func (api API) BuildHandler() *echo.Echo {

	api.configRouter()

	// Endpoint for swagger documentations
	api.router.GET("/swagger/*", echoSwagger.WrapHandler)

	err := otel.SetTraceProvider(api.cfg.OpenTelemetry.JaegerURL, api.cfg.Server.NAME, app.Version, api.cfg.Server.ENV.String(), api.cfg.OpenTelemetry.Sampled)
	if err != nil {
		api.log.Fatal(err)
	}

	repoRegistry := postgres.NewRepositoryRegistry(api.db)

	auth.RegisterAPI(
		*api.router.Group(""),
		api.cfg,
		auth.NewService(api.cfg, repoRegistry),
	)

	user.RegisterAPI(
		*api.router.Group(""),
		api.cfg,
		user.NewService(api.cfg, repoRegistry),
	)

	api.router.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"api": api.cfg.Server.NAME,
			"env": api.cfg.Server.ENV.String(),
		})
	})

	api.router.Any("", func(c echo.Context) error {
		return echo.NotFoundHandler(c)
	})

	api.router.Any("/*", func(c echo.Context) error {
		return echo.NotFoundHandler(c)
	})

	return api.router
}

func (api API) configRouter() {

	api.router.Pre(middleware.RemoveTrailingSlash())
	// api.router.Use(middleware.RequestID())
	api.router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	api.router.Use(customMiddleware.RequestIDContext())                  // middleware for insert request id into context
	api.router.Use(customMiddleware.HandlerTracing(api.cfg.Server.NAME)) // middleware for handling opentelemetry

	// Setup custom HTTP error handler
	api.router.HTTPErrorHandler = CustomHTTPErrorHandler(api.cfg, api.log, httplog.NewHTTPLog(api.db))

	// Register middleware recover from panic

	// Setup access log
	api.router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","request_id":"${id}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status},"error":"${error}","latency":"${latency}",` +
			`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
			`"service":"` + api.cfg.Server.NAME + `",` +
			`"environment":"` + api.cfg.Server.ENV.String() + `",` +
			`"type":"access",` +
			`"bytes_out":${bytes_out}}` + "\n",
	}))

	api.router.Use(customMiddleware.Recover(api.log))

}

func (api API) Start() {
	// handle graceful exit
	ctx, cancel := context.WithCancel(context.Background())
	handleSigterm(func() {
		cancel()
	})

	// Start server
	server := http.Server{
		Addr:    fmt.Sprintf(":%v", api.cfg.Server.PORT),
		Handler: api.BuildHandler(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			api.log.Fatalf("listen:%+s\n", err)
		}
	}()

	api.log.Infof("server is running at port: %v [env: %v, version: %v]", api.cfg.Server.PORT, api.cfg.Server.ENV, app.Version)

	gracefulShutdownServer(ctx, &server, api.log)
}

func gracefulShutdownServer(ctx context.Context, srv *http.Server, log logger.Logger) {

	<-ctx.Done()

	log.Info("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server shutdown failed:%+s", err)
	}

	log.Info("server exited properly")

}

func handleSigterm(exitFunc func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		exitFunc()
	}()
}
