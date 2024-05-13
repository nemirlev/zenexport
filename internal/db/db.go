package db

import (
	"github.com/nemirlev/zenapi"
)

// DataStore это интерфейс для базы данных. Методы специфичны для работы с данными ДзенМани.
type DataStore interface {
	Save(data *zenapi.Response) error
	Update(data interface{}) error
	Delete(data *zenapi.Deletion) error
}
