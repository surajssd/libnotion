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
	BTBulletListItem = BT("bulleted_list_item")
	NumberedListItem = BT("numbered_list_item")
	BTCallout        = BT("callout")
	BTParagraph      = BT("paragraph")
	BTHeading1       = BT("heading_1")
	BTHeading2       = BT("heading_2")
	BTHeading3       = BT("heading_3")
	BTTodo           = BT("to_do")
	BTQuote          = BT("quote")
	BTToggle         = BT("toggle")
	BTUnsupported    = BT("unsupported")
	BTTable          = BT("table")
	BTDivider        = BT("divider")
)
