package main

import (
	"net/http"
	"time"

	"github.com/gleich/hueport-scraper/internal/marketplace"
	"github.com/gleich/lumber/v3"
)

func main() {
	setupLogger()
	lumber.Info("booted")

	client := http.DefaultClient

	for {
		marketplace.ProcessExtensions(client)
		time.Sleep(5 * time.Minute)
	}
}

func setupLogger() {
	nytime, err := time.LoadLocation("America/New_York")
	if err != nil {
		lumber.Fatal(err, "failed to load new york timezone")
	}
	lumber.SetTimezone(nytime)
	lumber.SetTimeFormat("01/02 03:04:05 PM MST")
}
