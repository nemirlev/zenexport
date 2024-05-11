package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Config struct {
	ZenMoneyToken      string `mapstructure:"ZENMONEY_TOKEN"`
	IsDaemon           bool   `mapstructure:"IS_DAEMON"`
	DatabaseType       string `mapstructure:"DATABASE_TYPE"`
	DatabaseServer     string `mapstructure:"DATABASE_SERVER"`
	DatabaseUser       string `mapstructure:"DATABASE_USER"`
	DatabasePassword   string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseName       string `mapstructure:"DATABASE_NAME"`
	ClickhouseServer   string `mapstructure:"CLICKHOUSE_SERVER"`
	ClickhouseUser     string `mapstructure:"CLICKHOUSE_USER"`
	ClickhouseDB       string `mapstructure:"CLICKHOUSE_DB"`
	ClickhousePassword string `mapstructure:"CLICKHOUSE_PASSWORD"`
	Interval           int    `mapstructure:"INTERVAL"`
}

// DatabaseURL возращает строку подключения к базе данных
func (c Config) DatabaseURL() string {
	return fmt.Sprintf("%s://%s:%s@%s/%s", c.DatabaseType, c.DatabaseUser, c.DatabasePassword, c.DatabaseServer, c.DatabaseName)
}

// initViper инициализирует viper
func initViper() *viper.Viper {
	v := viper.New()
	v.SetConfigType("env")
	v.AutomaticEnv()

	v.SetDefault("ZENMONEY_TOKEN", "")
	v.SetDefault("IS_DAEMON", false)
	v.SetDefault("DATABASE_TYPE", "clickhouse")
	v.SetDefault("DATABASE_SERVER", "127.0.0.1")
	v.SetDefault("DATABASE_USER", "")
	v.SetDefault("DATABASE_PASSWORD", "")
	v.SetDefault("DATABASE_NAME", "")
	v.SetDefault("CLICKHOUSE_SERVER", "127.0.0.1")
	v.SetDefault("CLICKHOUSE_USER", "")
	v.SetDefault("CLICKHOUSE_DB", "")
	v.SetDefault("CLICKHOUSE_PASSWORD", "")
	v.SetDefault("INTERVAL", 1)

	return v
}

// FromEnv загружает конфигурацию из переменных окружения и флагов
func FromEnv() (*Config, error) {
	v := initViper()

	// Проверка, не тест ли это
	if !isTestEnvironment() {
		if !flag.Parsed() {
			defineFlags()
			flag.Parse()
		}
	}

	// Переопределение переменных окружения значениями флагов при их наличии
	applyFlagOverrides(v)

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	if cfg.ClickhouseUser == "" {
		return nil, fmt.Errorf("CLICKHOUSE_USER is required")
	}

	if cfg.ClickhouseDB == "" {
		return nil, fmt.Errorf("CLICKHOUSE_DB is required")
	}

	if cfg.ClickhousePassword == "" {
		return nil, fmt.Errorf("CLICKHOUSE_PASSWORD is required")
	}

	return cfg, nil
}

// defineFlags определяет флаги
func defineFlags() {
	flag.Int("interval", 0, "The interval in minutes to wait between syncs")
	flag.String("token", "", "The ZenMoney token. Get it from https://zerro.app/token")
	flag.String("dbtype", "", "The type of the database")
	flag.Bool("d", false, "Run as a daemon")
	flag.String("server", "", "The ClickHouse server")
	flag.String("user", "", "The ClickHouse user")
	flag.String("db", "", "The ClickHouse database")
	flag.String("password", "", "The ClickHouse password")
}

// applyFlagOverrides переопределяет переменные окружения значениями флагов при их наличии
func applyFlagOverrides(v *viper.Viper) {
	intervalFlag := flag.Lookup("interval")
	if intervalFlag != nil {
		intervalVal, ok := intervalFlag.Value.(flag.Getter)
		if ok && intervalVal.Get().(int) != 0 {
			v.Set("INTERVAL", intervalVal.Get().(int))
		}
	}

	tokenFlag := flag.Lookup("token")
	if tokenFlag != nil {
		tokenVal, ok := tokenFlag.Value.(flag.Getter)
		if ok && tokenVal.Get().(string) != "" {
			v.Set("ZENMONEY_TOKEN", tokenVal.Get().(string))
		}
	}

	dbTypeFlag := flag.Lookup("dbtype")
	if dbTypeFlag != nil {
		dbTypeVal, ok := dbTypeFlag.Value.(flag.Getter)
		if ok && dbTypeVal.Get().(string) != "" {
			v.Set("DATABASE_TYPE", dbTypeVal.Get().(string))
		}
	}

	// Server
	serverFlag := flag.Lookup("server")
	if serverFlag != nil {
		serverVal, ok := serverFlag.Value.(flag.Getter)
		if ok && serverVal.Get().(string) != "" {
			v.Set("CLICKHOUSE_SERVER", serverVal.Get().(string))
		}
	}

	userFlag := flag.Lookup("user")
	if userFlag != nil {
		userVal, ok := userFlag.Value.(flag.Getter)
		if ok && userVal.Get().(string) != "" {
			v.Set("CLICKHOUSE_USER", userVal.Get().(string))
		}
	}

	dbFlag := flag.Lookup("db")
	if dbFlag != nil {
		dbVal, ok := dbFlag.Value.(flag.Getter)
		if ok && dbVal.Get().(string) != "" {
			v.Set("CLICKHOUSE_DB", dbVal.Get().(string))
		}
	}

	passwordFlag := flag.Lookup("password")
	if passwordFlag != nil {
		passwordVal, ok := passwordFlag.Value.(flag.Getter)
		if ok && passwordVal.Get().(string) != "" {
			v.Set("CLICKHOUSE_PASSWORD", passwordVal.Get().(string))
		}
	}

	daemonFlag := flag.Lookup("d")
	if daemonFlag != nil {
		daemonVal, ok := daemonFlag.Value.(flag.Getter)
		if ok && daemonVal.Get().(bool) {
			v.Set("IS_DAEMON", daemonVal.Get().(bool))
		}
	}
}

func isTestEnvironment() bool {
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}
	return false
}
