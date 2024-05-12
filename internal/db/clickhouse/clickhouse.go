package clickhouse

import (
	"context"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/nemirlev/zenexport/internal/config"
)

type Store struct {
	Conn driver.Conn
}

func (s *Store) connect(cfg *config.Config) error {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{fmt.Sprintf("%s:9000", cfg.ClickhouseServer)},
			Auth: clickhouse.Auth{
				Database: cfg.ClickhouseDB,
				Username: cfg.ClickhouseUser,
				Password: cfg.ClickhousePassword,
			},
			Debugf: func(format string, v ...interface{}) {
				fmt.Printf(format, v)
			},
		})
	)

	if err != nil {
		return err
	}

	if err := conn.Ping(ctx); err != nil {
		var exception *clickhouse.Exception
		if errors.As(err, &exception) {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return err
	}

	s.Conn = conn
	return nil
}
