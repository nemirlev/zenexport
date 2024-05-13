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

func runSyncAndSave(log logger.Log, client *zenapi.Client, db db.DataStore) error {
	fmt.Println("Get data from ZenMoney...")
	resBody, err := client.FullSync()
	fmt.Println("Finished getting data from ZenMoney.")
	if err != nil {
		log.WithError(err, "error getting ZenMoney data")
		return err
	}

	fmt.Println("Save data to Database...")
	err = db.Save(&resBody)
	if err != nil {
		log.WithError(err, "error save ZenMoney data to DB")
		return err
	}
	fmt.Println("Import completed.")
	return nil
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

	dbase, err := db.NewDataStore(cfg, log)
	if err != nil {
		log.WithError(err, "failed to setup database")
		os.Exit(1)
	}

	if cfg.IsDaemon {
		interval := time.Duration(cfg.Interval) * time.Minute

		ticker := time.NewTicker(interval)

		for range ticker.C {
			start := time.Now()
			err := runSyncAndSave(log, client, dbase)
			if err != nil {
				log.WithError(err, "error sync ZenMoney data")
			}

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
		err := runSyncAndSave(log, client, dbase)
		if err != nil {
			log.WithError(err, "error sync ZenMoney data")
		}
	}
}
