package rest

const (
	// APIURL is the base URL of Notion API.
	APIURL = "https://api.notion.com"

	// NotionVersion is the Notion API release version.
	NotionVersion = "2025-09-03"

	// SubPathPages is the Notion API sub path for querying pages.
	SubPathPages = "v1/pages"

	// SubPathDatabases is the Notion API sub path for querying databases.
	SubPathDatabases = "v1/databases"

	// SubPathBlocks is the Notion API sub path for querying blocks.
	SubPathBlocks = "v1/blocks"

	// SubPathDataSources is the Notion API sub path for querying data sources.
	// Added in API version 2025-09-03 for multi-source databases.
	SubPathDataSources = "v1/data_sources"

	// SubPathUsers is the Notion API sub path for querying users.
	SubPathUsers = "v1/users"
)
