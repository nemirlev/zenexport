package db

import (
	"github.com/nemirlev/zenapi"
	"github.com/nemirlev/zenexport/internal/config"
)

// DataStore это интерфейс для базы данных. Методы специфичны для работы с данными ДзенМани.
type DataStore interface {
	Save(cfg *config.Config, data *zenapi.Response) error
	Update(cfg *config.Config, data interface{}) error
	Delete(cfg *config.Config, data *zenapi.Deletion) error
}
