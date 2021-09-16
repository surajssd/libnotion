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

// AddPage takes a page resource and adds it to the database.
func (nc *NotionClient) AddPage(pg api.Page) (*api.Page, error) {
	client := &http.Client{}

	u, err := url.Parse(APIURL)
	if err != nil {
		return nil, fmt.Errorf("parsing the APIURL: %w", err)
	}

	u.Path = path.Join(u.Path, SubPathPages)
	postURL := u.String()

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(pg); err != nil {
		return nil, fmt.Errorf("encoding page for request: %w", err)
	}

	req, err := http.NewRequest("POST", postURL, b)
	if err != nil {
		return nil, fmt.Errorf("building query: %w", err)
	}

	req.Header.Add("Notion-Version", NotionVersion)
	req.Header.Add("Authorization", "Bearer "+nc.token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("adding a new page: %w", err)
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

	page := api.Page{}
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("could not unmarshal response, %w", err)
	}

	return &page, nil
}
