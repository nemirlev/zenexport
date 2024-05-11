// Package export экспортирует все данные из Zen Money в БД ClickHouse
package main

import (
	"fmt"
	"github.com/nemirlev/zenapi"
	"github.com/nemirlev/zenexport/db"
	"github.com/nemirlev/zenexport/db/clickhouse"
	"github.com/nemirlev/zenexport/internal/config"
	"os"
	"time"
)

func createClient(token string) (*zenapi.Client, error) {
	return zenapi.NewClient(token)
}

func setupDatabase() (*clickhouse.ClickHouse, error) {
	return &clickhouse.ClickHouse{}, nil
}

func runSyncAndSave(cfg *config.Config, client *zenapi.Client, db db.DB) {
	fmt.Println("Starting import...")
	resBody, err := client.FullSync()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = db.Save(cfg, &resBody)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Import completed.")
}

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		//log.WithError(err, "get cfg")
		os.Exit(1)
	}

	client, err := createClient(cfg.ZenMoneyToken)
	if err != nil {
		fmt.Println(err, "failed to create client")
		return
	}

	dbase, err := setupDatabase()
	if err != nil {
		fmt.Println(err)
		return
	}

	if cfg.IsDaemon {
		interval := time.Duration(cfg.Interval) * time.Minute

		ticker := time.NewTicker(interval)

		for range ticker.C {
			start := time.Now()
			runSyncAndSave(cfg, client, dbase)

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
		runSyncAndSave(cfg, client, dbase)
	}
}
