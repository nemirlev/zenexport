package clickhouse

import "github.com/nemirlev/zenexport/internal/config"

func (s *Store) Update(cfg *config.Config, data interface{}) error {
	// TODO: Требуется разработать, по остаточному принципу, так как обновление в ClickHouse не тривиально
	panic("implement me")
}
