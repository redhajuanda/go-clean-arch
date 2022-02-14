package cmd

import (
	"context"
	"fmt"
	"go-clean-arch/app"
	"go-clean-arch/config"
	"go-clean-arch/docs"
	"go-clean-arch/infra/db"
	"go-clean-arch/infra/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var restCmd = &cobra.Command{
	Use: "rest",
	Run: func(_ *cobra.Command, _ []string) {
		startRestServer()
	},
}

func startRestServer() {
	cfg := config.LoadDefault()

	// set swagger config
	docs.SwaggerInfo.Host = cfg.Server.BASEURL
	docs.SwaggerInfo.Version = Version

	log := logger.New(cfg.Server.NAME, Version)
	logger.SetFormatter(&logrus.JSONFormatter{})

	db := db.NewGoPG(cfg)
	defer func() {
		log.Info("closing db connection")
		db.Close()
	}()

	// Init application
	application := app.New(cfg, echo.New(), db, log)

	// handle graceful exit
	ctx, cancel := context.WithCancel(context.Background())
	handleSigterm(func() {
		cancel()
	})

	// Start server
	server := http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Server.PORT),
		Handler: application.BuildHandler(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	log.Infof("server is running at port: %v [env: %v, version: %v]", cfg.Server.PORT, cfg.Server.ENV, Version)

	gracefulShutdownServer(ctx, &server, log)

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
