package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"pkg.mattglei.ch/hueport-scraper/internal/db"
	"pkg.mattglei.ch/hueport-scraper/internal/marketplace"
	"pkg.mattglei.ch/timber"
)

const TIME_FORMAT = "01/02 03:04:05 PM MST"

var newYork *time.Location

func main() {
	err := godotenv.Load()
	if err != nil {
		timber.Fatal(err, "Failed to load environment variables")
	}

	setupLogger()
	timber.Info("booted")

	database := db.Connect()
	client := http.DefaultClient

	cycleRate := 1 * time.Hour
	for {
		fmt.Println()
		start := time.Now()
		created, updated := marketplace.ProcessExtensions(client, database)
		timber.Info("cycle took", time.Since(start))
		timber.Done("loaded in", created, "extensions")
		timber.Done("updated", updated, "extensions")
		timber.Info(
			"next cycle will be at",
			time.Now().In(newYork).Add(cycleRate).Format(TIME_FORMAT),
		)
		time.Sleep(cycleRate)
	}
}

func setupLogger() {
	nytime, err := time.LoadLocation("America/New_York")
	if err != nil {
		timber.Fatal(err, "failed to load new york timezone")
	}
	newYork = nytime
	timber.SetTimezone(nytime)
	timber.SetTimeFormat(TIME_FORMAT)
}
