// Package export экспортирует все данные из Zen Money в БД ClickHouse
package main

import (
	"fmt"
	"github.com/nemirlev/zenapi"
	"github.com/nemirlev/zenexport/internal/config"
	"github.com/nemirlev/zenexport/internal/db"
	"github.com/nemirlev/zenexport/internal/logger"
	"os"
	"time"
)

func createClient(token string) (*zenapi.Client, error) {
	return zenapi.NewClient(token)
}

func runSyncAndSave(cfg *config.Config, client *zenapi.Client, db db.DataStore) {
	fmt.Println("Get data from ZenMoney...")
	resBody, err := client.FullSync()
	fmt.Println("Finished getting data from ZenMoney.")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Save data to Database...")
	err = db.Save(cfg, &resBody)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Import completed.")
}

func main() {
	log := logger.New()
	cfg, err := config.FromEnv()
	if err != nil {
		log.WithError(err, "get cfg")
		os.Exit(1)
	}

	client, err := createClient(cfg.ZenMoneyToken)
	if err != nil {
		log.WithError(err, "failed to create client")
		os.Exit(1)
	}

	dbase, err := db.NewDataStore(cfg)
	if err != nil {
		log.WithError(err, "failed to setup database")
		os.Exit(1)
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
