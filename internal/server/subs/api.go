package subs

import (
	"log/slog"

	"github.com/P3rCh1/subs-agregator/internal/config"
)

type APIServer struct {
	Logger *slog.Logger
	Config *config.Config
}

func NewAPIServer(logger *slog.Logger, config *config.Config) *APIServer {
	return &APIServer{
		Logger: logger,
		Config: config,
	}
}
