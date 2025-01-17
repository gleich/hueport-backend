package main

import (
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"pkg.mattglei.ch/hueport-scraper/internal/db"
	"pkg.mattglei.ch/hueport-scraper/internal/marketplace"
	"pkg.mattglei.ch/timber"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		timber.Fatal(err, "Failed to load environment variables")
	}

	setupLogger()
	timber.Info("booted")

	database := db.Connect()
	client := http.DefaultClient

	marketplace.ProcessExtensions(client, database)
}

func setupLogger() {
	nytime, err := time.LoadLocation("America/New_York")
	if err != nil {
		timber.Fatal(err, "failed to load new york timezone")
	}
	timber.SetTimezone(nytime)
	timber.SetTimeFormat("01/02 03:04:05 PM MST")
}
