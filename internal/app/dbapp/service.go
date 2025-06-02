package dbapp

import (
	"fmt"

	"github.com/antlko/golitedb/internal/db"
	"github.com/antlko/golitedb/internal/logger"
	"github.com/antlko/golitedb/internal/server"
)

func Start(cfg Config) error {
	logger.InitLogger(logger.Config{
		AppName:  cfg.ApplicationName,
		Hostname: cfg.Hostname,
	})

	dbInstance, err := db.NewDB(cfg.DB)
	if err != nil {
		return fmt.Errorf("new db: %w", err)
	}

	if err = server.Start(cfg.Server, dbInstance); err != nil {
		return fmt.Errorf("start server: %w", err)
	}
	return nil
}
