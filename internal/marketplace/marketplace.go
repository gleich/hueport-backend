package marketplace

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gleich/lumber/v3"
	"gorm.io/gorm"
	"pkg.mattglei.ch/hueport-scraper/internal/db"
)

func ProcessExtensions(client *http.Client, database *gorm.DB) {
	marketplaceExtensions, err := fetchExtensions(client)
	if err != nil {
		lumber.Fatal(err)
	}

	tempDir := resetProcessingFolder()

	for i, marketplaceExtension := range marketplaceExtensions {
		// check to make sure that extension doesn't already exist
		var extension db.Extension
		result := database.First(&extension, "extension_id = ?", marketplaceExtension.ExtensionID)
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			continue
		}

		fmt.Println()
		lumber.Info(
			"Processing",
			marketplaceExtension.DisplayName,
			fmt.Sprintf("(%d/%d)", i+1, len(marketplaceExtensions)),
		)
		zipPath, err := downloadExtension(client, tempDir, marketplaceExtension)
		if err != nil {
			lumber.Error(err, "failed to download extension")
			return
		}
		lumber.Done("✔︎ Downloaded")

		extensionFolder, err := unzipExtension(zipPath, marketplaceExtension)
		if err != nil {
			lumber.Error(err, "failed to unzip extension")
			return
		}
		lumber.Done("✔︎ Unzipped VSIX package")

		themes, err := extractThemes(extensionFolder, marketplaceExtension)
		if err != nil {
			lumber.Error(err, "failed to extract themes from extension")
			return
		}

		lumber.Done("✔︎ Extracted", len(themes), "themes")

		database.Create(
			&db.Extension{
				ExtensionID: marketplaceExtension.ExtensionID,
				Name:        marketplaceExtension.DisplayName,
				LastUpdated: marketplaceExtension.LastUpdated,
			},
		)

		lumber.Done("✔︎ Created", marketplaceExtension.DisplayName, "in database")

		lumber.Done("Finished processing", marketplaceExtension.DisplayName)
		resetProcessingFolder()
	}
}

func resetProcessingFolder() string {
	tempDir := filepath.Join(os.TempDir(), "hueport")
	err := os.RemoveAll(tempDir)
	if err != nil {
		lumber.Fatal(err, "failed to remove temporary dir", tempDir)
	}
	return tempDir
}
