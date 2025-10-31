package main

import (
	"flag"
	"fmt"
	"os"

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

	subs := subs.NewAPIServer(logger, cfg)

	router := gin.New()
	router.Use(gin.Recovery())

	router.POST("/subs", subs.Create)
	router.GET("/subs/:id", subs.Read)
	router.PUT("/subs/:id", subs.Update)
	router.DELETE("/subs/:id", subs.Delete)
	router.DELETE("/subs", subs.List)

	router.Run(cfg.Server.Host + ":" + cfg.Server.Port)
}
