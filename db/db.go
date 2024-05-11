package db

import (
	"github.com/nemirlev/zenapi"
	"github.com/nemirlev/zenexport/internal/config"
)

// DB это интерфейс для базы данных. Методы специфичны для работы с данными ДзенМани.
type DB interface {
	Save(cfg *config.Config, data *zenapi.Response) error
}

//	DatabaseEntity это интерфейс для сущностей в БД. Само описания сущностей в /internal/models
//	type DatabaseEntity interface {
//	GetTableName() string
//	GetInsertQuery() string
//	GetValues() []interface{}
//}
