package blocks

import (
	"encoding/json"
	"testing"

	"github.com/surajssd/libnotion/api"
)

// jsonRoundTrip marshals v to JSON and unmarshals into a new value of the same type.
func jsonRoundTrip[T any](t *testing.T, v T) T {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	return result
}

func TestBlockResponseList_JSON(t *testing.T) {
	bt := BTParagraph
	list := BlockResponseList{
		Response: api.Response{Object: "list", HasMore: true, NextCursor: "cursor-123"},
		Results: []Block{
			{
				CommonObject: api.CommonObject{ID: "block-1", Object: "block"},
				Type:         &bt,
			},
		},
	}
	got := jsonRoundTrip(t, list)
	if got.Object != "list" {
		t.Errorf("expected Object %q, got %q", "list", got.Object)
	}
	if !got.HasMore {
		t.Error("expected HasMore true")
	}
	if got.NextCursor != "cursor-123" {
		t.Errorf("expected NextCursor %q, got %q", "cursor-123", got.NextCursor)
	}
	if len(got.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got.Results))
	}
	if got.Results[0].ID != "block-1" {
		t.Errorf("expected ID %q, got %q", "block-1", got.Results[0].ID)
	}
}

func TestBlock_Paragraph(t *testing.T) {
	bt := BTParagraph
	block := Block{
		CommonObject: api.CommonObject{ID: "b-1", Object: "block"},
		Type:         &bt,
		HasChildren:  true,
		Paragraph: &Property{
			Text: []FullText{{PlainText: "Hello world"}},
		},
	}
	got := jsonRoundTrip(t, block)
	if got.Type == nil || *got.Type != BTParagraph {
		t.Errorf("expected type paragraph, got %v", got.Type)
	}
	if !got.HasChildren {
		t.Error("expected HasChildren true")
	}
	if got.Paragraph == nil || len(got.Paragraph.Text) != 1 {
		t.Fatal("expected paragraph text")
	}
	if got.Paragraph.Text[0].PlainText != "Hello world" {
		t.Errorf("expected %q, got %q", "Hello world", got.Paragraph.Text[0].PlainText)
	}
}

func TestBlock_Headings(t *testing.T) {
	tests := []struct {
		name string
		bt   BT
		set  func(*Block, *Property)
		get  func(Block) *Property
	}{
		{"heading_1", BTHeading1, func(b *Block, p *Property) { b.Heading1 = p }, func(b Block) *Property { return b.Heading1 }},
		{"heading_2", BTHeading2, func(b *Block, p *Property) { b.Heading2 = p }, func(b Block) *Property { return b.Heading2 }},
		{"heading_3", BTHeading3, func(b *Block, p *Property) { b.Heading3 = p }, func(b Block) *Property { return b.Heading3 }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bt := tt.bt
			block := Block{
				CommonObject: api.CommonObject{ID: "b-h"},
				Type:         &bt,
			}
			tt.set(&block, &Property{Text: []FullText{{PlainText: "Heading"}}})
			got := jsonRoundTrip(t, block)
			prop := tt.get(got)
			if prop == nil || len(prop.Text) != 1 || prop.Text[0].PlainText != "Heading" {
				t.Errorf("unexpected heading content: %+v", prop)
			}
		})
	}
}

func TestBlock_Lists(t *testing.T) {
	tests := []struct {
		name string
		bt   BT
		set  func(*Block, *Property)
		get  func(Block) *Property
	}{
		{"bulleted_list_item", BTBulletListItem, func(b *Block, p *Property) { b.BulletedListItem = p }, func(b Block) *Property { return b.BulletedListItem }},
		{"numbered_list_item", BTNumberedListItem, func(b *Block, p *Property) { b.NumberedListItem = p }, func(b Block) *Property { return b.NumberedListItem }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bt := tt.bt
			block := Block{Type: &bt}
			tt.set(&block, &Property{Text: []FullText{{PlainText: "Item"}}})
			got := jsonRoundTrip(t, block)
			prop := tt.get(got)
			if prop == nil || len(prop.Text) != 1 {
				t.Errorf("unexpected list item: %+v", prop)
			}
		})
	}
}

func TestBlock_Callout(t *testing.T) {
	bt := BTCallout
	block := Block{
		Type: &bt,
		Callout: &Property{
			Text: []FullText{{PlainText: "Note"}},
			Icon: &Icon{Type: "emoji", Emoji: "💡"},
		},
	}
	got := jsonRoundTrip(t, block)
	if got.Callout == nil {
		t.Fatal("expected Callout")
	}
	if got.Callout.Icon == nil || got.Callout.Icon.Emoji != "💡" {
		t.Errorf("unexpected icon: %+v", got.Callout.Icon)
	}
}

func TestBlock_Todo(t *testing.T) {
	bt := BTTodo
	block := Block{
		Type: &bt,
		Todo: &Property{
			Text:    []FullText{{PlainText: "Task"}},
			Checked: true,
		},
	}
	got := jsonRoundTrip(t, block)
	if got.Todo == nil {
		t.Fatal("expected Todo")
	}
	if !got.Todo.Checked {
		t.Error("expected Checked true")
	}
}

func TestBlock_Quote(t *testing.T) {
	bt := BTQuote
	block := Block{
		Type:  &bt,
		Quote: &Property{Text: []FullText{{PlainText: "A quote"}}},
	}
	got := jsonRoundTrip(t, block)
	if got.Quote == nil || len(got.Quote.Text) != 1 {
		t.Fatal("expected Quote with text")
	}
}

func TestBlock_Toggle(t *testing.T) {
	bt := BTToggle
	block := Block{
		Type:        &bt,
		HasChildren: true,
		Toggle:      &Property{Text: []FullText{{PlainText: "Toggle"}}},
	}
	got := jsonRoundTrip(t, block)
	if got.Toggle == nil {
		t.Fatal("expected Toggle")
	}
}

func TestBlock_ArchivedAndInTrash(t *testing.T) {
	bt := BTParagraph
	block := Block{
		Type:     &bt,
		Archived: true,
		InTrash:  true,
	}
	got := jsonRoundTrip(t, block)
	if !got.Archived {
		t.Error("expected Archived true")
	}
	if !got.InTrash {
		t.Error("expected InTrash true")
	}
}

func TestBlock_MediaTypes(t *testing.T) {
	tests := []struct {
		name string
		bt   BT
		set  func(*Block, *FileBlock)
		get  func(Block) *FileBlock
	}{
		{"image", BTImage, func(b *Block, f *FileBlock) { b.Image = f }, func(b Block) *FileBlock { return b.Image }},
		{"video", BTVideo, func(b *Block, f *FileBlock) { b.Video = f }, func(b Block) *FileBlock { return b.Video }},
		{"audio", BTAudio, func(b *Block, f *FileBlock) { b.Audio = f }, func(b Block) *FileBlock { return b.Audio }},
		{"file", BTFile, func(b *Block, f *FileBlock) { b.File = f }, func(b Block) *FileBlock { return b.File }},
		{"pdf", BTPDF, func(b *Block, f *FileBlock) { b.PDF = f }, func(b Block) *FileBlock { return b.PDF }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bt := tt.bt
			block := Block{Type: &bt}
			tt.set(&block, &FileBlock{
				Type:     "external",
				External: &External{URL: "https://example.com/file.png"},
				Caption:  []FullText{{PlainText: "A caption"}},
				Name:     "file.png",
			})
			got := jsonRoundTrip(t, block)
			fb := tt.get(got)
			if fb == nil {
				t.Fatal("expected FileBlock")
			}
			if fb.Type != "external" {
				t.Errorf("expected type %q, got %q", "external", fb.Type)
			}
			if fb.External == nil || fb.External.URL != "https://example.com/file.png" {
				t.Errorf("unexpected External: %+v", fb.External)
			}
			if fb.Name != "file.png" {
				t.Errorf("expected Name %q, got %q", "file.png", fb.Name)
			}
		})
	}
}

func TestBlock_Bookmark(t *testing.T) {
	bt := BTBookmark
	block := Block{
		Type: &bt,
		Bookmark: &BookmarkBlock{
			URL:     "https://example.com",
			Caption: []FullText{{PlainText: "A bookmark"}},
		},
	}
	got := jsonRoundTrip(t, block)
	if got.Bookmark == nil {
		t.Fatal("expected Bookmark")
	}
	if got.Bookmark.URL != "https://example.com" {
		t.Errorf("expected URL %q, got %q", "https://example.com", got.Bookmark.URL)
	}
	if len(got.Bookmark.Caption) != 1 {
		t.Errorf("expected 1 caption, got %d", len(got.Bookmark.Caption))
	}
}

func TestBlock_Embed(t *testing.T) {
	bt := BTEmbed
	block := Block{
		Type:  &bt,
		Embed: &EmbedBlock{URL: "https://embed.example.com"},
	}
	got := jsonRoundTrip(t, block)
	if got.Embed == nil || got.Embed.URL != "https://embed.example.com" {
		t.Errorf("unexpected Embed: %+v", got.Embed)
	}
}

func TestBlock_LinkPreview(t *testing.T) {
	bt := BTLinkPreview
	block := Block{
		Type:        &bt,
		LinkPreview: &LinkPreview{URL: "https://preview.example.com"},
	}
	got := jsonRoundTrip(t, block)
	if got.LinkPreview == nil || got.LinkPreview.URL != "https://preview.example.com" {
		t.Errorf("unexpected LinkPreview: %+v", got.LinkPreview)
	}
}

func TestBlock_Code(t *testing.T) {
	bt := BTCode
	block := Block{
		Type: &bt,
		Code: &CodeBlock{
			Language: "go",
			RichText: []FullText{{PlainText: "fmt.Println(\"hello\")"}},
			Caption:  []FullText{{PlainText: "Example code"}},
		},
	}
	got := jsonRoundTrip(t, block)
	if got.Code == nil {
		t.Fatal("expected Code")
	}
	if got.Code.Language != "go" {
		t.Errorf("expected language %q, got %q", "go", got.Code.Language)
	}
	if len(got.Code.RichText) != 1 {
		t.Fatalf("expected 1 rich text, got %d", len(got.Code.RichText))
	}
	if len(got.Code.Caption) != 1 {
		t.Fatalf("expected 1 caption, got %d", len(got.Code.Caption))
	}
}

func TestBlock_SyncedBlock(t *testing.T) {
	bt := BTSyncedBlock
	block := Block{
		Type: &bt,
		SyncedBlock: &SyncedBlock{
			SyncedFrom: &SyncedFrom{Type: "block_id", BlockID: "source-block-123"},
			Children: []Block{
				{CommonObject: api.CommonObject{ID: "child-1"}},
			},
		},
	}
	got := jsonRoundTrip(t, block)
	if got.SyncedBlock == nil {
		t.Fatal("expected SyncedBlock")
	}
	if got.SyncedBlock.SyncedFrom == nil || got.SyncedBlock.SyncedFrom.BlockID != "source-block-123" {
		t.Errorf("unexpected SyncedFrom: %+v", got.SyncedBlock.SyncedFrom)
	}
	if len(got.SyncedBlock.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(got.SyncedBlock.Children))
	}
}

func TestBlock_ChildPage(t *testing.T) {
	bt := BTChildPage
	block := Block{
		Type:      &bt,
		ChildPage: &ChildPage{Title: "My Child Page"},
	}
	got := jsonRoundTrip(t, block)
	if got.ChildPage == nil || got.ChildPage.Title != "My Child Page" {
		t.Errorf("unexpected ChildPage: %+v", got.ChildPage)
	}
}

func TestBlock_ChildDatabase(t *testing.T) {
	bt := BTChildDatabase
	block := Block{
		Type:          &bt,
		ChildDatabase: &ChildDatabase{Title: "My Child DB"},
	}
	got := jsonRoundTrip(t, block)
	if got.ChildDatabase == nil || got.ChildDatabase.Title != "My Child DB" {
		t.Errorf("unexpected ChildDatabase: %+v", got.ChildDatabase)
	}
}

func TestBlock_ColumnList(t *testing.T) {
	bt := BTColumnList
	block := Block{
		Type: &bt,
		ColumnList: &ColumnList{
			Children: []Block{{CommonObject: api.CommonObject{ID: "col-1"}}},
		},
	}
	got := jsonRoundTrip(t, block)
	if got.ColumnList == nil || len(got.ColumnList.Children) != 1 {
		t.Errorf("unexpected ColumnList: %+v", got.ColumnList)
	}
}

func TestBlock_Column(t *testing.T) {
	bt := BTColumn
	block := Block{
		Type: &bt,
		Column: &Column{
			Children: []Block{{CommonObject: api.CommonObject{ID: "col-child-1"}}},
		},
	}
	got := jsonRoundTrip(t, block)
	if got.Column == nil || len(got.Column.Children) != 1 {
		t.Errorf("unexpected Column: %+v", got.Column)
	}
}

func TestBlock_TableOfContents(t *testing.T) {
	bt := BTTableOfContents
	block := Block{
		Type:            &bt,
		TableOfContents: &TableOfContents{Color: "gray"},
	}
	got := jsonRoundTrip(t, block)
	if got.TableOfContents == nil || got.TableOfContents.Color != "gray" {
		t.Errorf("unexpected TableOfContents: %+v", got.TableOfContents)
	}
}

func TestBlock_Breadcrumb(t *testing.T) {
	bt := BTBreadcrumb
	block := Block{
		Type:       &bt,
		Breadcrumb: &Breadcrumb{},
	}
	got := jsonRoundTrip(t, block)
	if got.Breadcrumb == nil {
		t.Error("expected Breadcrumb to be non-nil")
	}
}

func TestBlock_Divider(t *testing.T) {
	bt := BTDivider
	block := Block{
		Type:    &bt,
		Divider: &Divider{},
	}
	got := jsonRoundTrip(t, block)
	if got.Divider == nil {
		t.Error("expected Divider to be non-nil")
	}
}

func TestBlock_Table(t *testing.T) {
	bt := BTTable
	block := Block{
		Type: &bt,
		Table: &TableBlock{
			TableWidth:      3,
			HasColumnHeader: true,
			HasRowHeader:    true,
		},
	}
	got := jsonRoundTrip(t, block)
	if got.Table == nil {
		t.Fatal("expected Table")
	}
	if got.Table.TableWidth != 3 {
		t.Errorf("expected TableWidth 3, got %d", got.Table.TableWidth)
	}
	if !got.Table.HasColumnHeader {
		t.Error("expected HasColumnHeader true")
	}
	if !got.Table.HasRowHeader {
		t.Error("expected HasRowHeader true")
	}
}

func TestBlock_TableRow(t *testing.T) {
	bt := BTTableRow
	block := Block{
		Type: &bt,
		TableRow: &TableRow{
			Cells: [][]FullText{
				{{PlainText: "Cell 1"}, {PlainText: "Cell 1b"}},
				{{PlainText: "Cell 2"}},
			},
		},
	}
	got := jsonRoundTrip(t, block)
	if got.TableRow == nil {
		t.Fatal("expected TableRow")
	}
	if len(got.TableRow.Cells) != 2 {
		t.Fatalf("expected 2 cells, got %d", len(got.TableRow.Cells))
	}
	if len(got.TableRow.Cells[0]) != 2 {
		t.Fatalf("expected 2 text items in first cell, got %d", len(got.TableRow.Cells[0]))
	}
	if got.TableRow.Cells[0][0].PlainText != "Cell 1" {
		t.Errorf("expected %q, got %q", "Cell 1", got.TableRow.Cells[0][0].PlainText)
	}
}

func TestBlock_Template(t *testing.T) {
	bt := BTTemplate
	block := Block{
		Type: &bt,
		Template: &TemplateBlock{
			RichText: []FullText{{PlainText: "Template Title"}},
			Children: []Block{
				{CommonObject: api.CommonObject{ID: "tmpl-child-1"}},
			},
		},
	}
	got := jsonRoundTrip(t, block)
	if got.Template == nil {
		t.Fatal("expected Template")
	}
	if len(got.Template.RichText) != 1 || got.Template.RichText[0].PlainText != "Template Title" {
		t.Errorf("unexpected template RichText: %+v", got.Template.RichText)
	}
	if len(got.Template.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(got.Template.Children))
	}
}

func TestProperty_JSON(t *testing.T) {
	prop := Property{
		Text:    []FullText{{PlainText: "text content"}},
		Icon:    &Icon{Type: "emoji", Emoji: "🔥"},
		Checked: true,
	}
	got := jsonRoundTrip(t, prop)
	if len(got.Text) != 1 || got.Text[0].PlainText != "text content" {
		t.Errorf("unexpected Text: %+v", got.Text)
	}
	if got.Icon == nil || got.Icon.Emoji != "🔥" {
		t.Errorf("unexpected Icon: %+v", got.Icon)
	}
	if !got.Checked {
		t.Error("expected Checked true")
	}
}

func TestFullText_JSON(t *testing.T) {
	ft := FullText{
		Type:        "text",
		Text:        &Text{Content: "Hello", Link: &api.Link{URL: "https://example.com"}},
		Annotations: &Annotations{Bold: true, Italic: true, Code: true, Color: "red"},
		PlainText:   "Hello",
	}
	got := jsonRoundTrip(t, ft)
	if got.Type != "text" {
		t.Errorf("expected type %q, got %q", "text", got.Type)
	}
	if got.Text == nil || got.Text.Content != "Hello" {
		t.Errorf("unexpected Text: %+v", got.Text)
	}
	if got.Text.Link == nil || got.Text.Link.URL != "https://example.com" {
		t.Errorf("unexpected Link: %+v", got.Text.Link)
	}
	if got.Annotations == nil || !got.Annotations.Bold || !got.Annotations.Italic || !got.Annotations.Code {
		t.Errorf("unexpected Annotations: %+v", got.Annotations)
	}
	if got.Annotations.Color != "red" {
		t.Errorf("expected color %q, got %q", "red", got.Annotations.Color)
	}
}

func TestFullText_Annotations_All(t *testing.T) {
	ann := Annotations{
		Bold:          true,
		Italic:        true,
		Strikethrough: true,
		Underline:     true,
		Code:          true,
		Color:         "blue",
	}
	got := jsonRoundTrip(t, ann)
	if !got.Bold || !got.Italic || !got.Strikethrough || !got.Underline || !got.Code {
		t.Error("expected all annotations to be true")
	}
	if got.Color != "blue" {
		t.Errorf("expected color %q, got %q", "blue", got.Color)
	}
}

func TestFileBlock_JSON(t *testing.T) {
	fb := FileBlock{
		Caption:  []FullText{{PlainText: "A file"}},
		Type:     "file",
		File:     &NotionFile{URL: "https://notion.so/file.pdf", ExpiryTime: "2025-01-01T00:00:00Z"},
		Name:     "document.pdf",
	}
	got := jsonRoundTrip(t, fb)
	if got.Type != "file" {
		t.Errorf("expected type %q, got %q", "file", got.Type)
	}
	if got.File == nil || got.File.URL != "https://notion.so/file.pdf" {
		t.Errorf("unexpected File: %+v", got.File)
	}
	if got.File.ExpiryTime != "2025-01-01T00:00:00Z" {
		t.Errorf("unexpected ExpiryTime: %q", got.File.ExpiryTime)
	}
	if got.Name != "document.pdf" {
		t.Errorf("expected Name %q, got %q", "document.pdf", got.Name)
	}
	if len(got.Caption) != 1 {
		t.Errorf("expected 1 caption, got %d", len(got.Caption))
	}
}

func TestFileBlock_External(t *testing.T) {
	fb := FileBlock{
		Type:     "external",
		External: &External{URL: "https://external.com/image.png"},
	}
	got := jsonRoundTrip(t, fb)
	if got.External == nil || got.External.URL != "https://external.com/image.png" {
		t.Errorf("unexpected External: %+v", got.External)
	}
}

func TestCodeBlock_JSON(t *testing.T) {
	cb := CodeBlock{
		Language: "python",
		RichText: []FullText{
			{PlainText: "print('hello')"},
		},
		Caption: []FullText{
			{PlainText: "Python example"},
		},
	}
	got := jsonRoundTrip(t, cb)
	if got.Language != "python" {
		t.Errorf("expected language %q, got %q", "python", got.Language)
	}
	if len(got.RichText) != 1 || got.RichText[0].PlainText != "print('hello')" {
		t.Errorf("unexpected RichText: %+v", got.RichText)
	}
	if len(got.Caption) != 1 || got.Caption[0].PlainText != "Python example" {
		t.Errorf("unexpected Caption: %+v", got.Caption)
	}
}

func TestSyncedBlock_JSON(t *testing.T) {
	sb := SyncedBlock{
		SyncedFrom: &SyncedFrom{Type: "block_id", BlockID: "source-123"},
		Children:   []Block{{CommonObject: api.CommonObject{ID: "child-1"}}},
	}
	got := jsonRoundTrip(t, sb)
	if got.SyncedFrom == nil || got.SyncedFrom.BlockID != "source-123" {
		t.Errorf("unexpected SyncedFrom: %+v", got.SyncedFrom)
	}
	if got.SyncedFrom.Type != "block_id" {
		t.Errorf("expected type %q, got %q", "block_id", got.SyncedFrom.Type)
	}
	if len(got.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(got.Children))
	}
}

func TestTableRow_JSON(t *testing.T) {
	tr := TableRow{
		Cells: [][]FullText{
			{{PlainText: "A1"}},
			{{PlainText: "B1"}, {PlainText: "B1-extra"}},
		},
	}
	got := jsonRoundTrip(t, tr)
	if len(got.Cells) != 2 {
		t.Fatalf("expected 2 cells, got %d", len(got.Cells))
	}
	if len(got.Cells[1]) != 2 {
		t.Fatalf("expected 2 texts in cell 1, got %d", len(got.Cells[1]))
	}
}

func TestBlockTypeVariables(t *testing.T) {
	types := map[string]BT{
		"bulleted_list_item": BTBulletListItem,
		"numbered_list_item": BTNumberedListItem,
		"callout":            BTCallout,
		"paragraph":          BTParagraph,
		"heading_1":          BTHeading1,
		"heading_2":          BTHeading2,
		"heading_3":          BTHeading3,
		"to_do":              BTTodo,
		"quote":              BTQuote,
		"toggle":             BTToggle,
		"unsupported":        BTUnsupported,
		"table":              BTTable,
		"divider":            BTDivider,
		"image":              BTImage,
		"video":              BTVideo,
		"audio":              BTAudio,
		"file":               BTFile,
		"pdf":                BTPDF,
		"bookmark":           BTBookmark,
		"embed":              BTEmbed,
		"link_preview":       BTLinkPreview,
		"code":               BTCode,
		"synced_block":       BTSyncedBlock,
		"child_page":         BTChildPage,
		"child_database":     BTChildDatabase,
		"column_list":        BTColumnList,
		"column":             BTColumn,
		"table_of_contents":  BTTableOfContents,
		"breadcrumb":         BTBreadcrumb,
		"table_row":          BTTableRow,
		"template":           BTTemplate,
	}
	for expected, got := range types {
		if string(got) != expected {
			t.Errorf("expected block type %q, got %q", expected, got)
		}
	}
	// Verify count — 31 block types
	if len(types) != 31 {
		t.Errorf("expected 31 block types, got %d", len(types))
	}
}
