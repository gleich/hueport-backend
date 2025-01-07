package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gleich/hueport-backend/internal/marketplace"
	"github.com/gleich/lumber/v3"
)

func main() {
	setupLogger()
	lumber.Info("booted")

	tempDir := resetProcessingFolder()

	client := http.DefaultClient
	marketplaceExtensions, err := marketplace.FetchExtensions(client)
	if err != nil {
		lumber.Fatal(err)
	}

	for i, marketplaceExtension := range marketplaceExtensions[1:] {
		fmt.Println()
		lumber.Info(
			"Processing",
			marketplaceExtension.DisplayName,
			fmt.Sprintf("(%d/%d)", i+1, len(marketplaceExtensions)),
		)
		zipPath, err := marketplace.DownloadExtension(client, tempDir, marketplaceExtension)
		if err != nil {
			lumber.Error(err, "failed to download extension")
			return
		}
		lumber.Done("✔︎ Downloaded")

		extensionFolder, err := marketplace.UnzipExtension(zipPath, marketplaceExtension)
		if err != nil {
			lumber.Error(err, "failed to unzip extension")
			return
		}
		lumber.Done("✔︎ Unzipped VSIX package")

		themes, err := marketplace.ExtractThemes(extensionFolder, marketplaceExtension)
		if err != nil {
			lumber.Error(err, "failed to extract themes from extension")
			return
		}

		lumber.Done("✔︎ Extracted", len(themes), "themes")
		lumber.Done("Finished Processed", marketplaceExtension.DisplayName)
		resetProcessingFolder()
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

func resetProcessingFolder() string {
	tempDir := filepath.Join(os.TempDir(), "hueport")
	err := os.RemoveAll(tempDir)
	if err != nil {
		lumber.Fatal(err, "failed to remove temporary dir", tempDir)
	}
	return tempDir
}
