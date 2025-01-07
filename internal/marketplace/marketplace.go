package marketplace

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gleich/lumber/v3"
)

func ProcessExtensions(client *http.Client) {
	marketplaceExtensions, err := fetchExtensions(client)
	if err != nil {
		lumber.Fatal(err)
	}

	tempDir := resetProcessingFolder()

	for i, marketplaceExtension := range marketplaceExtensions {
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
		lumber.Done("Finished Processed", marketplaceExtension.DisplayName)
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
