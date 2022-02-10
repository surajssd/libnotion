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
	"github.com/surajssd/libnotion/api/blocks"
)

func (nc *NotionClient) ListBlocks(id string) ([]blocks.Block, error) {
	hasMore := true
	startCursor := ""
	var ret []blocks.Block
	client := &http.Client{}

	u, err := url.Parse(APIURL)
	if err != nil {
		return nil, fmt.Errorf("parsing the APIURL: %w", err)
	}

	u.Path = path.Join(u.Path, SubPathBlocks, id, "children")
	listURL := u.String()

	for hasMore {
		req, err := http.NewRequest("GET", listURL, nil)
		if err != nil {
			return nil, fmt.Errorf("building query: %w", err)
		}

		req.Header.Add("Notion-Version", NotionVersion)
		req.Header.Add("Authorization", "Bearer "+nc.token)
		req.Header.Add("Content-Type", "application/json")

		if startCursor != "" {
			q := req.URL.Query()
			q.Add("start_cursor", startCursor)
			q.Add("page_size", "100")
			req.URL.RawQuery = q.Encode()
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("listing block entries: %w", err)
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

			return nil, fmt.Errorf("http request returned non-200 response: %s. Message: %s",
				resp.Status, failedResp.Message)
		}

		if respErr != nil {
			return nil, fmt.Errorf("reading the response: %w", err)
		}

		bl := blocks.BlockResponseList{}
		if err := json.Unmarshal(data, &bl); err != nil {
			return nil, fmt.Errorf("could not unmarshal response: %w", err)
		}

		hasMore = bl.HasMore
		startCursor = bl.NextCursor

		ret = append(ret, bl.Results...)
	}

	return ret, nil
}
