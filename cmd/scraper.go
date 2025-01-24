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

func main() {
	err := godotenv.Load()
	if err != nil {
		timber.Fatal(err, "Failed to load environment variables")
	}

	setupLogger()
	timber.Info("booted")

	database := db.Connect()
	client := http.DefaultClient

	cycleRate := 5 * time.Minute
	for {
		fmt.Println()
		created, updated := marketplace.ProcessExtensions(client, database)
		timber.Done("Loaded in", created, "extensions")
		timber.Done("Updated", updated, "extensions")
		timber.Info("Next cycle will be at", time.Now().Add(cycleRate).Format(TIME_FORMAT))
		time.Sleep(cycleRate)
	}
}

func setupLogger() {
	nytime, err := time.LoadLocation("America/New_York")
	if err != nil {
		timber.Fatal(err, "failed to load new york timezone")
	}
	timber.SetTimezone(nytime)
	timber.SetTimeFormat(TIME_FORMAT)
}
