package marketplace

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gleich/lumber/v3"
)

type extensionQueryResponse struct {
	Results []struct {
		Extensions []Extension `json:"extensions"`
	} `json:"results"`
}

type Extension struct {
	Publisher struct {
		DisplayName string `json:"displayName"`
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

func FetchExtensions(client *http.Client) ([]Extension, error) {
	extensions := []Extension{}
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
			lumber.Error(err, "failed to marshal JSON body")
			return []Extension{}, err
		}

		req, err := http.NewRequest(
			http.MethodPost,
			"https://marketplace.visualstudio.com/_apis/public/gallery/extensionquery",
			bytes.NewBuffer(reqBody),
		)
		if err != nil {
			lumber.Error(err, "failed to create new request")
			return []Extension{}, err
		}
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add(
			"Cookie",
			`VstsSession=%7B%22PersistentSessionId%22%3A%223428b623-854a-4549-a399-bb57078a5cce%22%2C%22PendingAuthenticationSessionId%22%3A%2200000000-0000-0000-0000-000000000000%22%2C%22CurrentAuthenticationSessionId%22%3A%2200000000-0000-0000-0000-000000000000%22%2C%22SignInState%22%3A%7B%7D%7D; Gallery-Service-UserIdentifier=78e578a5-a429-4e31-a599-4c742b906427`,
		)

		resp, err := client.Do(req)
		if err != nil {
			lumber.Error(err, "failed to execute request")
			return []Extension{}, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lumber.Error(err, "reading response body failed")
			return []Extension{}, err
		}

		var data extensionQueryResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			lumber.Error(err, "failed to parse json")
			lumber.Debug(string(body))
			return []Extension{}, err
		}

		if len(data.Results[0].Extensions) == 0 {
			lumber.Done("Finished fetching all extensions")
			break
		}

		extensions = append(extensions, data.Results[0].Extensions...)
		lumber.Done(
			"Fetched",
			len(data.Results[0].Extensions),
			"extensions. Total is at",
			len(extensions),
		)
	}
	return extensions, nil
}
