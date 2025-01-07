package main

import (
	"fmt"
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

	for i, extension := range extensions {
		fmt.Println()
		lumber.Info(
			"Processing",
			extension.DisplayName,
			fmt.Sprintf("(%d/%d)", i+1, len(extensions)),
		)
		path, err := marketplace.DownloadExtension(client, extension)
		if err != nil {
			lumber.Error(err, "failed to download extension")
		}
		lumber.Done("✔︎ Downloaded")

		err = marketplace.UnzipExtension(path, extension)
		if err != nil {
			lumber.Error(err, "failed to unzip extension")
		}
		lumber.Done("✔︎ Unzipped VSIX package")

		lumber.Done("Processed", extension.DisplayName)
		break
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
