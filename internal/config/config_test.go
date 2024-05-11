package config

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// Тестируем при назначении всех переменных
func TestFromEnvAllEnv(t *testing.T) {
	// Установка переменных окружения
	os.Setenv("ZENMONEY_TOKEN", "test_token")
	os.Setenv("IS_DAEMON", "true")
	os.Setenv("DATABASE_TYPE", "postgresql")
	os.Setenv("DATABASE_SERVER", "localhost")
	os.Setenv("DATABASE_USER", "test_user")
	os.Setenv("DATABASE_PASSWORD", "test_password")
	os.Setenv("DATABASE_NAME", "test_db")
	os.Setenv("CLICKHOUSE_SERVER", "localhost")
	os.Setenv("CLICKHOUSE_USER", "test_user")
	os.Setenv("CLICKHOUSE_DB", "test_db")
	os.Setenv("CLICKHOUSE_PASSWORD", "test_password")
	os.Setenv("INTERVAL", "1")

	// Вызов функции FromEnv
	cfg, err := FromEnv()

	// Проверка, что функция не возвращает ошибку
	assert.NoError(t, err)

	// Проверка, что значения в структуре Config соответствуют установленным переменным окружения
	assert.Equal(t, "test_token", cfg.ZenMoneyToken)
	assert.Equal(t, true, cfg.IsDaemon)
	assert.Equal(t, "postgresql", cfg.DatabaseType)
	assert.Equal(t, "localhost", cfg.DatabaseServer)
	assert.Equal(t, "test_user", cfg.DatabaseUser)
	assert.Equal(t, "test_password", cfg.DatabasePassword)
	assert.Equal(t, "test_db", cfg.DatabaseName)
	assert.Equal(t, "localhost", cfg.ClickhouseServer)
	assert.Equal(t, "test_user", cfg.ClickhouseUser)
	assert.Equal(t, "test_db", cfg.ClickhouseDB)
	assert.Equal(t, "test_password", cfg.ClickhousePassword)
	assert.Equal(t, 1, cfg.Interval)

	// Очистка переменных окружения
	os.Clearenv()
	// Сброс флагов
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
}

func TestFromEnvMissingClickhouseUser(t *testing.T) {
	// Установка переменных окружения
	os.Setenv("ZENMONEY_TOKEN", "test_token")
	os.Setenv("IS_DAEMON", "true")
	os.Setenv("DATABASE_TYPE", "postgresql")
	os.Setenv("DATABASE_SERVER", "localhost")
	os.Setenv("DATABASE_USER", "test_user")
	os.Setenv("DATABASE_PASSWORD", "test_password")
	os.Setenv("DATABASE_NAME", "test_db")
	os.Setenv("CLICKHOUSE_SERVER", "localhost")
	// Пропускаем CLICKHOUSE_USER
	os.Setenv("CLICKHOUSE_DB", "test_db")
	os.Setenv("CLICKHOUSE_PASSWORD", "test_password")
	os.Setenv("INTERVAL", "1")

	// Вызов функции FromEnv
	cfg, err := FromEnv()

	// Проверка, что функция возвращает ошибку
	assert.Error(t, err)
	// Проверка, что возвращаемое значение nil
	assert.Nil(t, cfg)
	// Проверка, что ошибка - это ожидаемая ошибка
	assert.Equal(t, "CLICKHOUSE_USER is required", err.Error())

	// Очистка переменных окружения
	os.Clearenv()
	// Сброс флагов
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
}

func TestFromEnvMissingClickhouseDB(t *testing.T) {
	// Установка переменных окружения
	os.Setenv("ZENMONEY_TOKEN", "test_token")
	os.Setenv("IS_DAEMON", "true")
	os.Setenv("DATABASE_TYPE", "postgresql")
	os.Setenv("DATABASE_SERVER", "localhost")
	os.Setenv("DATABASE_USER", "test_user")
	os.Setenv("DATABASE_PASSWORD", "test_password")
	os.Setenv("DATABASE_NAME", "test_db")
	os.Setenv("CLICKHOUSE_SERVER", "localhost")
	os.Setenv("CLICKHOUSE_USER", "test_user")
	// Пропускаем CLICKHOUSE_DB
	os.Setenv("CLICKHOUSE_PASSWORD", "test_password")
	os.Setenv("INTERVAL", "1")

	// Вызов функции FromEnv
	cfg, err := FromEnv()

	// Проверка, что функция возвращает ошибку
	assert.Error(t, err)
	// Проверка, что возвращаемое значение nil
	assert.Nil(t, cfg)
	// Проверка, что ошибка - это ожидаемая ошибка
	assert.Equal(t, "CLICKHOUSE_DB is required", err.Error())

	// Очистка переменных окружения
	os.Clearenv()
	// Сброс флагов
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
}

func TestFromEnvMissingClickhousePassword(t *testing.T) {
	// Установка переменных окружения
	os.Setenv("ZENMONEY_TOKEN", "test_token")
	os.Setenv("IS_DAEMON", "true")
	os.Setenv("DATABASE_TYPE", "postgresql")
	os.Setenv("DATABASE_SERVER", "localhost")
	os.Setenv("DATABASE_USER", "test_user")
	os.Setenv("DATABASE_PASSWORD", "test_password")
	os.Setenv("DATABASE_NAME", "test_db")
	os.Setenv("CLICKHOUSE_SERVER", "localhost")
	os.Setenv("CLICKHOUSE_USER", "test_user")
	os.Setenv("CLICKHOUSE_DB", "test_db")
	// Пропускаем CLICKHOUSE_PASSWORD
	os.Setenv("INTERVAL", "1")

	// Вызов функции FromEnv
	cfg, err := FromEnv()

	// Проверка, что функция возвращает ошибку
	assert.Error(t, err)
	// Проверка, что возвращаемое значение nil
	assert.Nil(t, cfg)
	// Проверка, что ошибка - это ожидаемая ошибка
	assert.Equal(t, "CLICKHOUSE_PASSWORD is required", err.Error())

	// Очистка переменных окружения
	os.Clearenv()
	// Сброс флагов
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
}
