package blocks

import "github.com/surajssd/libnotion/api"

// PageResponseList is used to parse the response when querying pages endpoint.
type BlockResponseList struct {
	api.Response `json:",inline"`
	Results      []Block `json:"results,omitempty"`
}

type Block struct {
	api.CommonObject `json:",inline"`
	HasChildren      bool      `json:"has_children,omitempty"`
	Archived         bool      `json:"archived,omitempty"`
	InTrash          bool      `json:"in_trash,omitempty"`
	Type             *BT       `json:"type,omitempty"`
	BulletedListItem *Property `json:"bulleted_list_item,omitempty"`
	NumberedListItem *Property `json:"numbered_list_item,omitempty"`
	Callout          *Property `json:"callout,omitempty"`
	Paragraph        *Property `json:"paragraph,omitempty"`
	Heading1         *Property `json:"heading_1,omitempty"`
	Heading2         *Property `json:"heading_2,omitempty"`
	Heading3         *Property `json:"heading_3,omitempty"`
	Todo             *Property `json:"to_do,omitempty"`
	Quote            *Property `json:"quote,omitempty"`

	// If the item is Toggle then we should go into this to explore all the blocks.
	Toggle *Property `json:"toggle,omitempty"`

	// Media block types
	Image *FileBlock `json:"image,omitempty"`
	Video *FileBlock `json:"video,omitempty"`
	Audio *FileBlock `json:"audio,omitempty"`
	File  *FileBlock `json:"file,omitempty"`
	PDF   *FileBlock `json:"pdf,omitempty"`

	// Link and embed block types
	Bookmark    *BookmarkBlock `json:"bookmark,omitempty"`
	Embed       *EmbedBlock    `json:"embed,omitempty"`
	LinkPreview *LinkPreview   `json:"link_preview,omitempty"`

	// Code block
	Code *CodeBlock `json:"code,omitempty"`

	// Synced block
	SyncedBlock *SyncedBlock `json:"synced_block,omitempty"`

	// Child page and database
	ChildPage     *ChildPage     `json:"child_page,omitempty"`
	ChildDatabase *ChildDatabase `json:"child_database,omitempty"`

	// Layout blocks
	ColumnList      *ColumnList      `json:"column_list,omitempty"`
	Column          *Column          `json:"column,omitempty"`
	TableOfContents *TableOfContents `json:"table_of_contents,omitempty"`
	Breadcrumb      *Breadcrumb      `json:"breadcrumb,omitempty"`
	Divider         *Divider         `json:"divider,omitempty"`

	// Table blocks
	Table    *TableBlock `json:"table,omitempty"`
	TableRow *TableRow   `json:"table_row,omitempty"`

	// Template block
	Template *TemplateBlock `json:"template,omitempty"`
}

type Property struct {
	Text    []FullText `json:"text,omitempty"`
	Icon    *Icon      `json:"icon,omitempty"`
	Checked bool       `json:"checked,omitempty"`
}

type Icon struct {
	Type  string `json:"type,omitempty"`
	Emoji string `json:"emoji,omitempty"`
}

type FullText struct {
	Type        string       `json:"type,omitempty"`
	Text        *Text        `json:"text,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`
	PlainText   string       `json:"plain_text,omitempty"`
}

type Text struct {
	Content string    `json:"content,omitempty"`
	Link    *api.Link `json:"link,omitempty"`
}
type Annotations struct {
	Bold          bool   `json:"bold,omitempty"`
	Italic        bool   `json:"italic,omitempty"`
	Strikethrough bool   `json:"strikethrough,omitempty"`
	Underline     bool   `json:"underline,omitempty"`
	Code          bool   `json:"code,omitempty"`
	Color         string `json:"color,omitempty"`
}

// BlockType
type BT string

var (
	BTBulletListItem  = BT("bulleted_list_item")
	BTNumberedListItem = BT("numbered_list_item")
	BTCallout         = BT("callout")
	BTParagraph       = BT("paragraph")
	BTHeading1        = BT("heading_1")
	BTHeading2        = BT("heading_2")
	BTHeading3        = BT("heading_3")
	BTTodo            = BT("to_do")
	BTQuote           = BT("quote")
	BTToggle          = BT("toggle")
	BTUnsupported     = BT("unsupported")
	BTTable           = BT("table")
	BTDivider         = BT("divider")
	// Media block types
	BTImage = BT("image")
	BTVideo = BT("video")
	BTAudio = BT("audio")
	BTFile  = BT("file")
	BTPDF   = BT("pdf")
	// Link and embed block types
	BTBookmark    = BT("bookmark")
	BTEmbed       = BT("embed")
	BTLinkPreview = BT("link_preview")
	// Code block
	BTCode = BT("code")
	// Synced block
	BTSyncedBlock = BT("synced_block")
	// Child page and database
	BTChildPage     = BT("child_page")
	BTChildDatabase = BT("child_database")
	// Layout blocks
	BTColumnList      = BT("column_list")
	BTColumn          = BT("column")
	BTTableOfContents = BT("table_of_contents")
	BTBreadcrumb      = BT("breadcrumb")
	// Table blocks
	BTTableRow = BT("table_row")
	// Template block
	BTTemplate = BT("template")
)

// FileBlock represents media blocks (image, video, audio, file, pdf).
type FileBlock struct {
	Caption  []FullText  `json:"caption,omitempty"`
	Type     string      `json:"type,omitempty"` // "external" or "file"
	External *External   `json:"external,omitempty"`
	File     *NotionFile `json:"file,omitempty"`
	Name     string      `json:"name,omitempty"`
}

// External represents an externally hosted file.
type External struct {
	URL string `json:"url,omitempty"`
}

// NotionFile represents a file hosted by Notion.
type NotionFile struct {
	URL        string `json:"url,omitempty"`
	ExpiryTime string `json:"expiry_time,omitempty"`
}

// BookmarkBlock represents a bookmark block.
type BookmarkBlock struct {
	Caption []FullText `json:"caption,omitempty"`
	URL     string     `json:"url,omitempty"`
}

// EmbedBlock represents an embed block.
type EmbedBlock struct {
	URL string `json:"url,omitempty"`
}

// LinkPreview represents a link preview block.
type LinkPreview struct {
	URL string `json:"url,omitempty"`
}

// CodeBlock represents a code block.
type CodeBlock struct {
	Caption  []FullText `json:"caption,omitempty"`
	RichText []FullText `json:"rich_text,omitempty"`
	Language string     `json:"language,omitempty"`
}

// SyncedBlock represents a synced block.
type SyncedBlock struct {
	SyncedFrom *SyncedFrom `json:"synced_from,omitempty"`
	Children   []Block     `json:"children,omitempty"`
}

// SyncedFrom represents the source of a synced block.
type SyncedFrom struct {
	Type    string `json:"type,omitempty"`
	BlockID string `json:"block_id,omitempty"`
}

// ChildPage represents a child page block.
type ChildPage struct {
	Title string `json:"title,omitempty"`
}

// ChildDatabase represents a child database block.
type ChildDatabase struct {
	Title string `json:"title,omitempty"`
}

// ColumnList represents a column list block.
type ColumnList struct {
	Children []Block `json:"children,omitempty"`
}

// Column represents a column block.
type Column struct {
	Children []Block `json:"children,omitempty"`
}

// TableOfContents represents a table of contents block.
type TableOfContents struct {
	Color string `json:"color,omitempty"`
}

// Breadcrumb represents a breadcrumb block.
type Breadcrumb struct{}

// Divider represents a divider block.
type Divider struct{}

// TableBlock represents a table block.
type TableBlock struct {
	TableWidth      int  `json:"table_width,omitempty"`
	HasColumnHeader bool `json:"has_column_header,omitempty"`
	HasRowHeader    bool `json:"has_row_header,omitempty"`
}

// TableRow represents a table row block.
type TableRow struct {
	Cells [][]FullText `json:"cells,omitempty"`
}

// TemplateBlock represents a template block.
type TemplateBlock struct {
	RichText []FullText `json:"rich_text,omitempty"`
	Children []Block    `json:"children,omitempty"`
}
