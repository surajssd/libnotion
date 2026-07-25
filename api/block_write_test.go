package api

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewParagraphBlock_Chunking(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		wantChunks []string
	}{
		{
			name:       "empty text yields empty rich_text",
			text:       "",
			wantChunks: []string{},
		},
		{
			name:       "short text is a single chunk",
			text:       "hello world",
			wantChunks: []string{"hello world"},
		},
		{
			name:       "exactly at the limit is a single chunk",
			text:       strings.Repeat("a", 2000),
			wantChunks: []string{strings.Repeat("a", 2000)},
		},
		{
			name: "one over the limit splits into two chunks",
			text: strings.Repeat("a", 2001),
			wantChunks: []string{
				strings.Repeat("a", 2000),
				"a",
			},
		},
		{
			name: "long text splits into multiple chunks",
			text: strings.Repeat("b", 4500),
			wantChunks: []string{
				strings.Repeat("b", 2000),
				strings.Repeat("b", 2000),
				strings.Repeat("b", 500),
			},
		},
		{
			name:       "multibyte runes are counted as characters not bytes",
			text:       strings.Repeat("日", 2000),
			wantChunks: []string{strings.Repeat("日", 2000)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := NewParagraphBlock(tt.text)

			if block.Object != "block" {
				t.Errorf("expected object %q, got %q", "block", block.Object)
			}
			if block.Type != "paragraph" {
				t.Errorf("expected type %q, got %q", "paragraph", block.Type)
			}
			if block.Paragraph == nil {
				t.Fatal("expected non-nil paragraph")
			}
			if len(block.Paragraph.RichText) != len(tt.wantChunks) {
				t.Fatalf("expected %d chunks, got %d", len(tt.wantChunks), len(block.Paragraph.RichText))
			}
			for i, want := range tt.wantChunks {
				if got := block.Paragraph.RichText[i].Text.Content; got != want {
					t.Errorf("chunk %d: expected %d chars, got %d chars", i, len(want), len(got))
				}
			}
		})
	}
}

func TestNewParagraphBlock_EmptyMarshalsEmptyRichText(t *testing.T) {
	data, err := json.Marshal(NewParagraphBlock(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `{"object":"block","type":"paragraph","paragraph":{"rich_text":[]}}`
	if string(data) != want {
		t.Errorf("expected %s, got %s", want, string(data))
	}
}

func TestNewParagraphBlock_TextMarshal(t *testing.T) {
	data, err := json.Marshal(NewParagraphBlock("line one"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `{"object":"block","type":"paragraph","paragraph":{"rich_text":[{"type":"text","text":{"content":"line one"},"annotations":{}}]}}`
	if string(data) != want {
		t.Errorf("expected %s, got %s", want, string(data))
	}
}

func TestParagraphBlocksFromText(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		wantLines []string
	}{
		{
			name:      "empty string yields nil",
			text:      "",
			wantLines: nil,
		},
		{
			name:      "single line",
			text:      "just one line",
			wantLines: []string{"just one line"},
		},
		{
			name:      "multiple lines with a blank line in between",
			text:      "line one\n\nline three",
			wantLines: []string{"line one", "", "line three"},
		},
		{
			name:      "carriage returns are trimmed",
			text:      "line one\r\nline two",
			wantLines: []string{"line one", "line two"},
		},
		{
			name:      "trailing newline yields trailing empty paragraph",
			text:      "line one\n",
			wantLines: []string{"line one", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocks := ParagraphBlocksFromText(tt.text)

			if tt.wantLines == nil {
				if blocks != nil {
					t.Fatalf("expected nil blocks, got %d", len(blocks))
				}
				return
			}

			if len(blocks) != len(tt.wantLines) {
				t.Fatalf("expected %d blocks, got %d", len(tt.wantLines), len(blocks))
			}

			for i, want := range tt.wantLines {
				rt := blocks[i].Paragraph.RichText
				if want == "" {
					if len(rt) != 0 {
						t.Errorf("block %d: expected empty rich_text, got %d segments", i, len(rt))
					}
					continue
				}
				if len(rt) != 1 {
					t.Fatalf("block %d: expected 1 segment, got %d", i, len(rt))
				}
				if got := rt[0].Text.Content; got != want {
					t.Errorf("block %d: expected %q, got %q", i, want, got)
				}
			}
		})
	}
}

func TestPage_WithoutChildrenMarshalsNoChildrenKey(t *testing.T) {
	data, err := json.Marshal(Page{
		Parent: Parent{
			Type:         ParentTypeDataSource,
			DataSourceID: "ds-123",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(string(data), `"children"`) {
		t.Errorf("expected no children key in marshalled page, got %s", string(data))
	}
}

func TestPage_WithChildrenMarshalsChildrenKey(t *testing.T) {
	data, err := json.Marshal(Page{
		Children: ParagraphBlocksFromText("some details"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(string(data), `"children"`) {
		t.Errorf("expected children key in marshalled page, got %s", string(data))
	}
}
