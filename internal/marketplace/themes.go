package marketplace

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/muhammadmuzzammil1998/jsonc"
	"pkg.mattglei.ch/hueport-scraper/pkg/models"
	"pkg.mattglei.ch/timber"
)

type Theme struct {
	Name   string `json:"name"`
	Colors struct {
		Foreground        string `json:"terminal.foreground"`
		Background        string `json:"terminal.background"`
		TabActiveBorder   string `json:"terminal.tab.activeBorder"`
		CursorBackground  string `json:"terminalCursor.background"`
		CursorForeground  string `json:"terminalCursor.foreground"`
		AnsiBrightWhite   string `json:"terminal.ansiBrightWhite"`
		AnsiWhite         string `json:"terminal.ansiWhite"`
		AnsiBrightBlack   string `json:"terminal.ansiBrightBlack"`
		AnsiBlack         string `json:"terminal.ansiBlack"`
		AnsiBlue          string `json:"terminal.ansiBlue"`
		AnsiBrightBlue    string `json:"terminal.ansiBrightBlue"`
		AnsiGreen         string `json:"terminal.ansiGreen"`
		AnsiBrightGreen   string `json:"terminal.ansiBrightGreen"`
		AnsiCyan          string `json:"terminal.ansiCyan"`
		AnsiBrightCyan    string `json:"terminal.ansiBrightCyan"`
		AnsiRed           string `json:"terminal.ansiRed"`
		AnsiBrightRed     string `json:"terminal.ansiBrightRed"`
		AnsiMagenta       string `json:"terminal.ansiMagenta"`
		AnsiBrightMagenta string `json:"terminal.ansiBrightMagenta"`
		AnsiYellow        string `json:"terminal.ansiYellow"`
		AnsiBrightYellow  string `json:"terminal.ansiBrightYellow"`
	} `json:"colors"`
}

func extractThemes(loc string, extension MarketplaceExtension) ([]models.Theme, error) {
	folder := filepath.Join(loc, "extension", "themes")
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		timber.Warning(extension.DisplayName, "doesn't have a themes folder")
		return []models.Theme{}, nil
	}
	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("%v failed to get file system entries from %s", err, folder)
	}

	themes := []models.Theme{}
	for _, e := range entries {
		name := e.Name()
		if !e.IsDir() && strings.HasSuffix(strings.ToLower(name), ".json") {
			bin, err := os.ReadFile(filepath.Join(folder, name))
			if err != nil {
				return nil, fmt.Errorf("%v failed to read from theme JSON file", err)
			}

			var theme Theme
			err = json.Unmarshal(jsonc.ToJSON(bin), &theme)
			if err != nil {
				timber.Warning(name, "did not contain proper json data")
				continue
			}
			themes = append(themes, models.Theme{
				Name:          theme.Name,
				ExtensionID:   extension.ExtensionID,
				Foreground:    theme.Colors.Foreground,
				Background:    theme.Colors.Background,
				BrightWhite:   theme.Colors.AnsiBrightWhite,
				White:         theme.Colors.AnsiWhite,
				BrightBlack:   theme.Colors.AnsiBrightBlack,
				Black:         theme.Colors.AnsiBlack,
				BrightBlue:    theme.Colors.AnsiBrightBlue,
				Blue:          theme.Colors.AnsiBlue,
				BrightGreen:   theme.Colors.AnsiBrightGreen,
				Green:         theme.Colors.AnsiGreen,
				BrightCyan:    theme.Colors.AnsiBrightCyan,
				Cyan:          theme.Colors.AnsiCyan,
				BrightRed:     theme.Colors.AnsiBrightRed,
				Red:           theme.Colors.AnsiRed,
				BrightMagenta: theme.Colors.AnsiBrightMagenta,
				Magenta:       theme.Colors.AnsiMagenta,
				BrightYellow:  theme.Colors.AnsiBrightYellow,
				Yellow:        theme.Colors.AnsiYellow,
			})
		}
	}
	return themes, nil
}
