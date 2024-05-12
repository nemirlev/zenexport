package db

import (
	"fmt"
	"github.com/nemirlev/zenexport/internal/config"
	"github.com/nemirlev/zenexport/internal/db/clickhouse"
)

func NewDataStore(cfg *config.Config) (DataStore, error) {
	switch cfg.DatabaseType {
	case "clickhouse":
		// Инициализация и конфигурация для ClickHouse
		return &clickhouse.Store{}, nil
	//case "postgres":
	//	// Инициализация и конфигурация для PostgreSQL
	//	return &postgres.PostgresStore{}, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DatabaseType)
	}
}
