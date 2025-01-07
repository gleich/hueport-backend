package marketplace

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gleich/lumber/v3"
	"github.com/muhammadmuzzammil1998/jsonc"
)

type Theme struct {
	Name   string `json:"name"`
	Colors struct {
		Foreground        string `json:"terminal.foreground"`
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

func extractThemes(loc string, extension MarketplaceExtension) ([]Theme, error) {
	folder := filepath.Join(loc, "extension", "themes")
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		lumber.Warning(extension.DisplayName, "doesn't have a themes folder")
		return []Theme{}, nil
	}
	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("%v failed to get file system entries from %s", err, folder)
	}

	themes := []Theme{}
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
				lumber.Warning(name, "did not contain proper json data")
				continue
			}
			themes = append(themes, theme)
		}
	}
	return themes, nil
}
