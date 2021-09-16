package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/surajssd/libnotion/api"
)

func (nc *NotionClient) FindDatabase(name string) (*api.Database, error) {
	hasMore := true
	startCursor := ""
	client := &http.Client{}

	u, err := url.Parse(APIURL)
	if err != nil {
		return nil, fmt.Errorf("parsing the APIURL: %w", err)
	}

	u.Path = path.Join(u.Path, SubPathDatabases)
	listURL := u.String()

	for hasMore {
		req, err := http.NewRequest("GET", listURL, nil)
		if err != nil {
			return nil, fmt.Errorf("building query: %w", err)
		}

		req.Header.Add("Notion-Version", NotionVersion)
		req.Header.Add("Authorization", "Bearer "+nc.token)

		if startCursor != "" {
			q := req.URL.Query()
			q.Add("start_cursor", startCursor)
			req.URL.RawQuery = q.Encode()
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("listing databases: %w", err)
		}

		defer resp.Body.Close()
		data, respErr := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			failedResp := api.FailureResponse{}

			if respErr != nil {
				log.Debugf("reading the response: %v", err)
			} else {
				if err := json.Unmarshal(data, &failedResp); err != nil {
					log.Debugf("unmarshalling failure response: %v", err)
				}
			}

			return nil, fmt.Errorf("http request returned non-200 response: %q. Message: %s",
				resp.Status, failedResp.Message)
		}

		// Check if there is any error while reading the response Data.
		if respErr != nil {
			return nil, fmt.Errorf("reading the response: %w", err)
		}

		databases := api.DatabaseResponseList{}
		if err := json.Unmarshal(data, &databases); err != nil {
			return nil, fmt.Errorf("could not unmarshal response, %w", err)
		}

		hasMore = databases.HasMore
		startCursor = databases.NextCursor

		for _, db := range databases.Results {
			if len(db.Title) == 0 {
				continue
			}

			foundName := db.Title[0].Text.Content
			if foundName == "" {
				continue
			}

			if foundName == name {
				return &db, nil
			}
		}
	}

	return nil, fmt.Errorf("database %q not found", name)
}
