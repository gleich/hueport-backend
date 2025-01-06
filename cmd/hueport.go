package main

import (
	"net/http"

	"github.com/gleich/hueport-backend/internal/marketplace"
	"github.com/gleich/lumber/v3"
)

func main() {
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
