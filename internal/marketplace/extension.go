package marketplace

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"pkg.mattglei.ch/timber"
)

type extensionQueryResponse struct {
	Results []struct {
		Extensions []MarketplaceExtension `json:"extensions"`
	} `json:"results"`
}

type MarketplaceExtension struct {
	Publisher struct {
		DisplayName   string `json:"displayName"`
		PublisherName string `json:"publisherName"`
	} `json:"publisher"`
	ExtensionID      string    `json:"extensionId"`
	ExtensionName    string    `json:"extensionName"`
	DisplayName      string    `json:"displayName"`
	Flags            string    `json:"flags"`
	LastUpdated      time.Time `json:"lastUpdated"`
	ShortDescription string    `json:"shortDescription"`
	Versions         []struct {
		Version     string    `json:"version"`
		Flags       string    `json:"flags"`
		LastUpdated time.Time `json:"lastUpdated"`
		Files       []struct {
			AssetType string `json:"assetType"`
			Source    string `json:"source"`
		} `json:"files"`
		Properties []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"properties"`
	} `json:"versions"`
	Tags       []string `json:"tags"`
	Statistics []struct {
		StatisticName string  `json:"statisticName"`
		Value         float32 `json:"value"`
	} `json:"statistics"`
}

type extensionQuery struct {
	AssetTypes []string `json:"assetTypes"`
	Filters    []filter `json:"filters"`
	Flags      int      `json:"flags"`
}

type filter struct {
	Criteria    []criteria  `json:"criteria"`
	Direction   int         `json:"direction"`
	PageSize    int         `json:"pageSize"`
	PageNumber  int         `json:"pageNumber"`
	SortBy      int         `json:"sortBy"`
	SortOrder   int         `json:"sortOrder"`
	PagingToken interface{} `json:"pagingToken"`
}

type criteria struct {
	FilterType int    `json:"filterType"`
	Value      string `json:"value"`
}

func fetchExtensions(client *http.Client) ([]MarketplaceExtension, error) {
	extensions := []MarketplaceExtension{}
	for i := 1; true; i++ {
		reqBody, err := json.Marshal(extensionQuery{
			AssetTypes: []string{
				"Microsoft.VisualStudio.Services.Icons.Default",
				"Microsoft.VisualStudio.Services.Icons.Branding",
				"Microsoft.VisualStudio.Services.Icons.Small",
			},
			Filters: []filter{{
				Criteria: []criteria{
					{FilterType: 8, Value: "Microsoft.VisualStudio.Code"},
					{FilterType: 10, Value: `target:"Microsoft.VisualStudio.Code" `},
					{FilterType: 12, Value: "37888"},
					{FilterType: 5, Value: "Themes"},
				},
				Direction:   2,
				PageSize:    1000,
				PageNumber:  i,
				SortBy:      4,
				SortOrder:   0,
				PagingToken: nil,
			}},
			Flags: 870,
		})
		if err != nil {
			return []MarketplaceExtension{}, fmt.Errorf("%v failed to marshal JSON body", err)
		}

		req, err := http.NewRequest(
			http.MethodPost,
			"https://marketplace.visualstudio.com/_apis/public/gallery/extensionquery",
			bytes.NewBuffer(reqBody),
		)
		if err != nil {
			return []MarketplaceExtension{}, fmt.Errorf("%v failed to make new request", err)
		}
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add(
			"Cookie",
			`VstsSession=%7B%22PersistentSessionId%22%3A%223428b623-854a-4549-a399-bb57078a5cce%22%2C%22PendingAuthenticationSessionId%22%3A%2200000000-0000-0000-0000-000000000000%22%2C%22CurrentAuthenticationSessionId%22%3A%2200000000-0000-0000-0000-000000000000%22%2C%22SignInState%22%3A%7B%7D%7D; Gallery-Service-UserIdentifier=78e578a5-a429-4e31-a599-4c742b906427`,
		)

		resp, err := client.Do(req)
		if err != nil {
			return []MarketplaceExtension{}, fmt.Errorf("%v failed to execute request", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return []MarketplaceExtension{}, fmt.Errorf("%v reading response body failed", err)
		}

		var data extensionQueryResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			timber.Debug(string(body))
			return []MarketplaceExtension{}, fmt.Errorf("%v failed to parse json", err)
		}

		if len(data.Results[0].Extensions) == 0 {
			timber.Done("Finished fetching all extensions")
			break
		}

		extensions = append(extensions, data.Results[0].Extensions...)
		timber.Done(
			"Fetched",
			len(data.Results[0].Extensions),
			"extensions. Total is at",
			len(extensions),
		)
	}
	return extensions, nil
}

// Downloads the extension to a temporary directory so it can be processed
func downloadExtension(
	client *http.Client,
	tempDir string,
	extension MarketplaceExtension,
) (string, error) {
	u, err := url.JoinPath(
		"https://marketplace.visualstudio.com/_apis/public/gallery/publishers/",
		extension.Publisher.PublisherName,
		"vsextensions",
		extension.ExtensionName, extension.Versions[0].Version, "vspackage",
	)
	if err != nil {
		return "", fmt.Errorf(
			"%v failed to URL encode for extension: %s",
			err,
			extension.DisplayName,
		)
	}

	resp, err := client.Get(u)
	if err != nil {
		return "", fmt.Errorf("%v failed to fetch extension %s", err, extension.DisplayName)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%v reading response body failed", err)
	}

	loc := filepath.Join(tempDir, fmt.Sprintf("%s.zip", extension.ExtensionID))
	folder := path.Dir(loc)

	if _, err = os.Stat(folder); os.IsNotExist(err) {
		err = os.MkdirAll(folder, 0777)
		if err != nil {
			return "", fmt.Errorf("%v failed to create directory", err)
		}
	}

	err = os.WriteFile(loc, body, 0655)
	if err != nil {
		return "", fmt.Errorf("%v failed to write VSIX file", err)
	}
	return loc, nil
}

func unzipExtension(zipPath string, extension MarketplaceExtension) (string, error) {
	folder := strings.TrimSuffix(zipPath, ".zip")

	if err := os.MkdirAll(folder, 0755); err != nil {
		return folder, fmt.Errorf(
			"failed to create folder %s for %s: %w",
			folder,
			extension.DisplayName,
			err,
		)
	}
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return folder, fmt.Errorf(
			"failed to open zip %s for %s: %w",
			zipPath,
			extension.DisplayName,
			err,
		)
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(folder, f.Name)
		if !strings.HasPrefix(
			filepath.Clean(fpath),
			filepath.Clean(folder)+string(os.PathSeparator),
		) {
			return folder, fmt.Errorf("illegal file path: %s", fpath)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, f.Mode()); err != nil {
				return folder, fmt.Errorf(
					"failed to create directory %s for %s: %w",
					fpath,
					extension.DisplayName,
					err,
				)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return folder, fmt.Errorf("failed to create directory for %s: %w", fpath, err)
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return folder, fmt.Errorf(
				"failed to open file %s for writing for %s: %w",
				fpath,
				extension.DisplayName,
				err,
			)
		}
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return folder, fmt.Errorf(
				"failed to open zip entry %s for %s: %w",
				f.Name,
				extension.DisplayName,
				err,
			)
		}
		if _, err := io.Copy(outFile, rc); err != nil {
			rc.Close()
			outFile.Close()
			return folder, fmt.Errorf(
				"failed to write file %s for %s: %w",
				fpath,
				extension.DisplayName,
				err,
			)
		}
		rc.Close()
		outFile.Close()
	}

	return folder, nil
}
