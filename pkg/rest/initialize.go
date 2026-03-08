package rest

// NotionClient is used to interact with Notion.
type NotionClient struct {
	token   string
	baseURL string
}

// NewNotionClient is used to initialize the Notion client.
//
// e.g., usage:
// rest.NewNotionClient(
//     rest.WithSecretToken(token)
// )
func NewNotionClient(fns ...notionClientConfigOpt) *NotionClient {
	var ret NotionClient

	for _, fn := range fns {
		fn(&ret)
	}

	return &ret
}

type notionClientConfigOpt func(*NotionClient)

// WithSecretToken is used to add Notion secret token to the Notion client during initialization.
func WithSecretToken(token string) notionClientConfigOpt {
	return func(nc *NotionClient) {
		nc.token = token
	}
}

// WithBaseURL is used to override the default Notion API base URL during initialization.
// This is useful for testing with a local server.
func WithBaseURL(baseURL string) notionClientConfigOpt {
	return func(nc *NotionClient) {
		nc.baseURL = baseURL
	}
}

// getBaseURL returns the base URL for the Notion API. If a custom base URL is set, it returns that;
// otherwise, it returns the default APIURL.
func (nc *NotionClient) getBaseURL() string {
	if nc.baseURL != "" {
		return nc.baseURL
	}
	return APIURL
}
