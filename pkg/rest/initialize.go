package rest

type NotionClient struct {
	token string
}

func NewNotionClient(fns ...notionClientConfigOpt) *NotionClient {
	var ret NotionClient

	for _, fn := range fns {
		fn(&ret)
	}

	return &ret
}

type notionClientConfigOpt func(*NotionClient)

func WithSecretToken(token string) notionClientConfigOpt {
	return func(nc *NotionClient) {
		nc.token = token
	}
}
