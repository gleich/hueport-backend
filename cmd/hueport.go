package main

import (
	"net/http"
	"time"

	"github.com/gleich/hueport-backend/internal/marketplace"
	"github.com/gleich/lumber/v3"
)

func main() {
	setupLogger()
	lumber.Info("booted")

	client := http.DefaultClient
	extensions, err := marketplace.FetchExtensions(client)
	if err != nil {
		lumber.Fatal(err)
	}

	for i, e := range extensions {
		lumber.Debug(i, e.DisplayName)
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
