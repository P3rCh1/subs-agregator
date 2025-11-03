package subs

import (
	"log/slog"

	"github.com/P3rCh1/subs-agregator/internal/config"
	"github.com/P3rCh1/subs-agregator/internal/storage/postgres"
)

type ServerAPI struct {
	Logger *slog.Logger
	Config *config.Config
	DB     *postgres.SubsAPI
}

func NewServerAPI(logger *slog.Logger, config *config.Config, db *postgres.SubsAPI) *ServerAPI {
	return &ServerAPI{
		Logger: logger,
		Config: config,
		DB:     db,
	}
}
