package marketplace

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"gorm.io/gorm"
	"pkg.mattglei.ch/hueport-scraper/pkg/models"
	"pkg.mattglei.ch/timber"
)

type ProcessType int

const (
	Updated ProcessType = iota
	Created ProcessType = iota
)

func ProcessExtensions(client *http.Client, database *gorm.DB) (int, int) {
	extensions, err := fetchExtensions(client)
	if err != nil {
		timber.Fatal(err)
	}

	tempDir := filepath.Join(os.TempDir(), "hueport")
	err = os.RemoveAll(tempDir)
	if err != nil {
		timber.Fatal(err, "failed to remove temporary directory")
	}

	workers, err := strconv.Atoi(os.Getenv("WORKERS"))
	if err != nil {
		timber.Fatal(err, "failed to parse number of workers")
	}
	tasks := make(chan MarketplaceExtension, len(extensions))

	go func() {
		for _, e := range extensions {
			tasks <- e
		}
		close(tasks)
	}()

	var (
		wg      sync.WaitGroup
		created int
		updated int
		mutex   sync.Mutex
	)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for extension := range tasks {
				processType, err := processExtension(client, tempDir, extension, database)
				if err != nil {
					timber.Error(err)
				}
				mutex.Lock()
				if processType == Updated {
					updated++
				} else if processType == Created {
					created++
				}
				mutex.Unlock()
			}
		}()
	}

	wg.Wait()
	return created, updated
}

func processExtension(
	client *http.Client,
	tempDir string,
	marketplaceExtension MarketplaceExtension,
	database *gorm.DB,
) (ProcessType, error) {
	var dbExtension models.Extension
	result := database.First(
		&dbExtension,
		"extension_id = ?",
		marketplaceExtension.ExtensionID,
	)
	new := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if result.Error != nil && !new {
		return -1, fmt.Errorf("%v failed to get extension from database", result.Error)
	}
	// skipping extension if it has no themes or hasn't be updated
	if !new && dbExtension.Themes == 0 ||
		marketplaceExtension.LastUpdated.Equal(dbExtension.LastUpdated) {
		return -1, nil
	}

	zipPath, err := downloadExtension(client, tempDir, marketplaceExtension)
	if err != nil {
		return -1, fmt.Errorf("%v failed to download extension", err)
	}

	extensionFolder, err := unzipExtension(zipPath, marketplaceExtension)
	if err != nil {
		return -1, fmt.Errorf("%v failed to unzip extension", err)
	}

	themes, err := extractThemes(extensionFolder, marketplaceExtension)
	if err != nil {
		return -1, fmt.Errorf("%v failed to extract themes from extension", err)
	}

	err = os.RemoveAll(extensionFolder)
	if err != nil {
		return -1, fmt.Errorf("%v failed to remove extension folder", err)
	}
	err = os.RemoveAll(zipPath)
	if err != nil {
		return -1, fmt.Errorf("%v failed to remove zip file", err)
	}

	extension := models.Extension{
		ExtensionID: marketplaceExtension.ExtensionID,
		Name:        marketplaceExtension.DisplayName,
		LastUpdated: marketplaceExtension.LastUpdated,
		Themes:      len(themes),
	}

	if new {
		result = database.Create(&extension)
		if result.Error != nil {
			return 0, fmt.Errorf("%v failed to create extension in database", result.Error)
		}
		timber.Done("created", extension.Name)
		return Created, nil
	}

	dbExtension.LastUpdated = extension.LastUpdated
	dbExtension.Themes = extension.Themes
	dbExtension.Name = extension.Name
	result = database.Save(&dbExtension)
	if result.Error != nil {
		timber.Error(result.Error, "failed to update extension")
	}
	timber.Done("updated", extension.Name)
	return Updated, nil
}
