package rest

// NotionClient is used to interact with Notion.
type NotionClient struct {
	token string
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
