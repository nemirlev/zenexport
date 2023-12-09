// Package export экспортирует все данные из Zen Money в БД ClickHouse
package main

import (
	"flag"
	"fmt"
	"github.com/nemirlev/zenapi"
	"github.com/nemirlev/zenexport/db"
	"github.com/nemirlev/zenexport/db/clickhouse"
	"os"
	"time"
)

type Flags struct {
	IntervalMinutes int
	Token           string
	Daemon          bool
	DBType          string
	DatabaseFlags
}

type DatabaseFlags struct {
	Server   string
	User     string
	DB       string
	Password string
}

func parseFlags() Flags {
	intervalMinutes := flag.Int("interval", 30, "The interval in minutes to wait between syncs")
	tokenFlag := flag.String("token", "", "The ZenMoney token. Get it from https://zerro.app/token")
	daemon := flag.Bool("d", false, "Run as a daemon")
	dbType := flag.String("dbtype", "clickhouse", "The type of the database")
	serverFlag := flag.String("server", "", "The ClickHouse server")
	userFlag := flag.String("user", "", "The ClickHouse user")
	dbFlag := flag.String("db", "", "The ClickHouse database")
	passwordFlag := flag.String("password", "", "The ClickHouse password")

	flag.Parse()

	return Flags{
		IntervalMinutes: *intervalMinutes,
		Token:           *tokenFlag,
		Daemon:          *daemon,
		DBType:          *dbType,
		DatabaseFlags: DatabaseFlags{
			Server:   *serverFlag,
			User:     *userFlag,
			DB:       *dbFlag,
			Password: *passwordFlag,
		},
	}
}

func createClient(token string) (*zenapi.Client, error) {
	if token == "" {
		token = os.Getenv("ZENMONEY_TOKEN")
	}
	return zenapi.NewClient(token)
}

func setupDatabase(flags DatabaseFlags) (*clickhouse.ClickHouse, error) {
	if flags.Server != "" {
		err := os.Setenv("CLICKHOUSE_SERVER", flags.Server)
		if err != nil {
			return nil, err
		}
	}
	if flags.User != "" {
		err := os.Setenv("CLICKHOUSE_USER", flags.User)
		if err != nil {
			return nil, err
		}
	}
	if flags.DB != "" {
		err := os.Setenv("CLICKHOUSE_DB", flags.DB)
		if err != nil {
			return nil, err
		}
	}
	if flags.Password != "" {
		err := os.Setenv("CLICKHOUSE_PASSWORD", flags.Password)
		if err != nil {
			return nil, err
		}
	}

	return &clickhouse.ClickHouse{}, nil
}

func runSyncAndSave(client *zenapi.Client, db db.DB) {
	fmt.Println("Starting import...")
	resBody, err := client.FullSync()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = db.Save(&resBody)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Import completed.")
}

func main() {
	flags := parseFlags()

	if flags.DBType != "clickhouse" {
		fmt.Println("Only 'clickhouse' is supported at the moment, but you can open a PR.")
		return
	}

	client, err := createClient(flags.Token)
	if err != nil {
		fmt.Println(err, "failed to create client")
		return
	}

	dbase, err := setupDatabase(flags.DatabaseFlags)
	if err != nil {
		fmt.Println(err)
		return
	}

	if flags.Daemon {
		interval := time.Duration(flags.IntervalMinutes) * time.Minute

		ticker := time.NewTicker(interval)

		for range ticker.C {
			start := time.Now()
			runSyncAndSave(client, dbase)

			nextTick := start.Add(interval)

			timer := time.NewTicker(1 * time.Second)
			go func() {
				for range timer.C {
					timeLeft := time.Until(nextTick).Round(time.Second)
					fmt.Printf("\rNext run in %v", timeLeft)
				}
			}()

			<-time.After(time.Until(nextTick))
			timer.Stop()
			fmt.Println()
		}
	} else {
		runSyncAndSave(client, dbase)
	}
}
