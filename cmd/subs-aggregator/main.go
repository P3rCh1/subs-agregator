package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/P3rCh1/subs-agregator/internal/config"
	"github.com/P3rCh1/subs-agregator/internal/logger"
	"github.com/P3rCh1/subs-agregator/internal/server/handlers/subs"
	mw "github.com/P3rCh1/subs-agregator/internal/server/middleware"
	"github.com/P3rCh1/subs-agregator/internal/storage/postgres"
	"github.com/labstack/echo"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config.yaml", "config path")
	flag.Parse()

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config parse fail: %s\n", err)
		os.Exit(1)
	}

	logger, err := logger.Setup(&cfg.Logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "setup logger fail: %s\n", err)
		os.Exit(1)
	}

	db, err := postgres.NewSubsAPI(&cfg.Postgres)
	if err != nil {
		logger.Error(
			"postgres connection fail",
			"error", err,
		)
		os.Exit(1)
	}
	defer db.Close()

	subs := subs.NewServerAPI(logger, cfg, db)

	router := SetupServer(subs)

	go func() {
		logger.Info("start server")
		err := router.Start(router.Server.Addr)
		if err != nil && err != http.ErrServerClosed {
			logger.Error(
				"listen",
				"error", err,
			)

			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("starting shutdown server")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := router.Shutdown(ctx); err != nil {
		logger.Error(
			"shutdown server",
			"error", err,
		)

		return
	}

	logger.Info("server stopped gracefully")
}

func SetupServer(subs *subs.ServerAPI) *echo.Echo {
	router := echo.New()

	router.Debug = false
	router.HideBanner = true
	router.Logger.SetOutput(io.Discard)

	router.Use(mw.Recover(subs.Logger))
	router.Use(mw.Logger(subs.Logger))

	router.POST("/subs", subs.Create)
	router.GET("/subs/:id", subs.Read)
	router.PUT("/subs/:id", subs.Update)
	router.DELETE("/subs/:id", subs.Delete)
	router.GET("/subs/list/:id", subs.List)
	router.POST("/subs/summary", subs.Summary)

	router.Server.Addr = subs.Config.HTTP.Host + ":" + subs.Config.HTTP.Port
	router.Server.ReadTimeout = subs.Config.HTTP.ReadTimeout
	router.Server.WriteTimeout = subs.Config.HTTP.WriteTimeout
	router.Server.IdleTimeout = subs.Config.HTTP.IdleTimeout

	return router
}
