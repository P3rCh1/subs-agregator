package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/P3rCh1/subs-agregator/internal/config"
	"github.com/P3rCh1/subs-agregator/internal/logger"
	"github.com/P3rCh1/subs-agregator/internal/server/subs"
	"github.com/gin-gonic/gin"
)

func main() {
	cfgPath := flag.String("c", "config.yaml", "config path")
	flag.Parse()

	cfg, err := config.ParseFile(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config parse fail: %s\n", err)
		os.Exit(1)
	}

	logger, err := logger.Setup(cfg.Logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "setup logger fail: %s\n", err)
		os.Exit(1)
	}

	subs := subs.NewServerAPI(logger, cfg)

	router := gin.New()
	router.Use(gin.Recovery())

	router.POST("/subs", subs.Create)
	router.GET("/subs/:id", subs.Read)
	router.PUT("/subs/:id", subs.Update)
	router.DELETE("/subs/:id", subs.Delete)
	router.GET("/subs", subs.List)

	server := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		err := server.ListenAndServe()
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

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error(
			"shutdown server",
			"error", err,
		)

		return
	}

	logger.Info("server stopped gracefully")
}
