package db

import (
	"fmt"
	"github.com/nemirlev/zenexport/internal/config"
	"github.com/nemirlev/zenexport/internal/db/clickhouse"
	"github.com/nemirlev/zenexport/internal/logger"
)

func NewDataStore(cfg *config.Config, log logger.Log) (DataStore, error) {
	switch cfg.DatabaseType {
	case "clickhouse":
		// Инициализация и конфигурация для ClickHouse
		return &clickhouse.Store{
			Log:    log,
			Config: cfg,
		}, nil
	//case "postgres":
	//	// Инициализация и конфигурация для PostgreSQL
	//	return &postgres.PostgresStore{}, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DatabaseType)
	}
}
