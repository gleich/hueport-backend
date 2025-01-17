package marketplace

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"gorm.io/gorm"
	"pkg.mattglei.ch/hueport-scraper/internal/db"
	"pkg.mattglei.ch/timber"
)

func ProcessExtensions(client *http.Client, database *gorm.DB) {
	marketplaceExtensions, err := fetchExtensions(client)
	if err != nil {
		timber.Fatal(err)
	}

	tempDir := resetProcessingFolder()

	for i, marketplaceExtension := range marketplaceExtensions {
		progress := fmt.Sprintf("(%d/%d)", i+1, len(marketplaceExtensions))
		// check to make sure that extension doesn't already exist
		var extension db.Extension
		result := database.First(&extension, "extension_id = ?", marketplaceExtension.ExtensionID)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			timber.Fatal(result.Error, "failed to get extension from database")
		}

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) && extension.Themes == 0 ||
			marketplaceExtension.LastUpdated.Equal(extension.LastUpdated) {
			continue
		}

		fmt.Println()
		timber.Info("Processing", marketplaceExtension.DisplayName, progress)
		zipPath, err := downloadExtension(client, tempDir, marketplaceExtension)
		if err != nil {
			timber.Error(err, "failed to download extension")
			return
		}
		timber.Done("✔︎ Downloaded")

		extensionFolder, err := unzipExtension(zipPath, marketplaceExtension)
		if err != nil {
			timber.Error(err, "failed to unzip extension")
			return
		}
		timber.Done("✔︎ Unzipped VSIX package")

		themes, err := extractThemes(extensionFolder, marketplaceExtension)
		if err != nil {
			timber.Error(err, "failed to extract themes from extension")
			return
		}

		timber.Done("✔︎ Extracted", len(themes), "themes")

		database.Create(
			&db.Extension{
				ExtensionID: marketplaceExtension.ExtensionID,
				Name:        marketplaceExtension.DisplayName,
				LastUpdated: marketplaceExtension.LastUpdated,
				Themes:      len(themes),
			},
		)
		timber.Done("✔︎ Created", marketplaceExtension.DisplayName, "in database")

		timber.Done("Finished processing", marketplaceExtension.DisplayName)
		resetProcessingFolder()
		fmt.Println()
	}
}

func resetProcessingFolder() string {
	tempDir := filepath.Join(os.TempDir(), "hueport")
	err := os.RemoveAll(tempDir)
	if err != nil {
		timber.Fatal(err, "failed to remove temporary dir", tempDir)
	}
	return tempDir
}
