package api

import "strings"

// maxTextContentLength is the Notion API limit on the content of a single text
// object.
const maxTextContentLength = 2000

// ChildBlock is a minimal write-only block model for creating pages with body
// content. Responses are parsed with the richer read model in api/blocks.
type ChildBlock struct {
	Object    string        `json:"object,omitempty"`
	Type      string        `json:"type"`
	Paragraph *RichTextBody `json:"paragraph,omitempty"`
}

// RichTextBody holds the rich text of a block. RichText intentionally has no
// omitempty: an empty paragraph must marshal as "rich_text": [].
type RichTextBody struct {
	RichText []Title `json:"rich_text"`
}

// NewParagraphBlock returns a paragraph child block for the given text,
// chunking it into text objects of at most 2000 characters each. Empty text
// yields an empty paragraph.
func NewParagraphBlock(text string) ChildBlock {
	richText := []Title{}

	runes := []rune(text)
	for start := 0; start < len(runes); start += maxTextContentLength {
		end := start + maxTextContentLength
		if end > len(runes) {
			end = len(runes)
		}

		richText = append(richText, Title{
			Type: "text",
			Text: Text{Content: string(runes[start:end])},
		})
	}

	return ChildBlock{
		Object:    "block",
		Type:      "paragraph",
		Paragraph: &RichTextBody{RichText: richText},
	}
}

// ParagraphBlocksFromText converts free-form text into one paragraph block per
// line. Blank lines become empty paragraphs; an empty string yields nil.
func ParagraphBlocksFromText(text string) []ChildBlock {
	if text == "" {
		return nil
	}

	lines := strings.Split(text, "\n")

	blocks := make([]ChildBlock, 0, len(lines))
	for _, line := range lines {
		blocks = append(blocks, NewParagraphBlock(strings.TrimRight(line, "\r")))
	}

	return blocks
}
