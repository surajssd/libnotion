package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/surajssd/libnotion/api"
)

// QueryDatabase takes database id and a query object and returns list of pages based on the query.
// Set appropriate parameters in the query object to get the relevant results.
func (nc *NotionClient) QueryDatabase(id string, query *api.QueryDB) ([]api.Page, error) {
	hasMore := true
	startCursor := ""
	client := &http.Client{}
	var ret []api.Page

	u, err := url.Parse(APIURL)
	if err != nil {
		return nil, fmt.Errorf("parsing the APIURL: %w", err)
	}

	u.Path = path.Join(u.Path, SubPathDataSources, id, "query")
	listURL := u.String()

	var body io.Reader

	for hasMore {
		if query != nil {
			query.StartCursor = startCursor

			bodyBytes, err := json.Marshal(query)
			if err != nil {
				return nil, fmt.Errorf("marshalling to bytes: %w", err)
			}

			body = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequest("POST", listURL, body)
		if err != nil {
			return nil, fmt.Errorf("building query: %w", err)
		}

		req.Header.Add("Notion-Version", NotionVersion)
		req.Header.Add("Authorization", "Bearer "+nc.token)
		req.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("listing database entries: %w", err)
		}

		defer resp.Body.Close()

		data, respErr := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			failedResp := api.FailureResponse{}

			if respErr != nil {
				log.Debugf("reading the response: %v", respErr)
			} else {
				if err := json.Unmarshal(data, &failedResp); err != nil {
					log.Debugf("unmarshalling failure response: %v", err)
				}
			}

			return nil, fmt.Errorf("http request returned non-200 response: %s. Message: %s",
				resp.Status, failedResp.Message)
		}

		if respErr != nil {
			return nil, fmt.Errorf("reading the response: %w", err)
		}

		pages := api.PageResponseList{}
		if err := json.Unmarshal(data, &pages); err != nil {
			return nil, fmt.Errorf("could not unmarshal response, %w", err)
		}

		hasMore = pages.HasMore
		startCursor = pages.NextCursor

		ret = append(ret, pages.Results...)
	}

	return ret, nil
}
