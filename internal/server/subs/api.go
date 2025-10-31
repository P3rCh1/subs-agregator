package subs

import (
	"log/slog"

	"github.com/P3rCh1/subs-agregator/internal/config"
)

type ServerAPI struct {
	Logger *slog.Logger
	Config *config.Config
}

func NewServerAPI(logger *slog.Logger, config *config.Config) *ServerAPI {
	return &ServerAPI{
		Logger: logger,
		Config: config,
	}
}
