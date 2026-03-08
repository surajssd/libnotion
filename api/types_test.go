package api

import (
	"encoding/json"
	"testing"
)

// jsonRoundTrip marshals v to JSON and unmarshals into a new value of the same type,
// returning the result. It fails the test if either step errors.
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

func TestResponse_JSON(t *testing.T) {
	resp := Response{
		Object:     "list",
		NextCursor: "cursor-abc",
		HasMore:    true,
		Status:     200,
		Code:       "success",
		Message:    "OK",
	}
	got := jsonRoundTrip(t, resp)
	if got.Object != "list" {
		t.Errorf("expected Object %q, got %q", "list", got.Object)
	}
	if got.NextCursor != "cursor-abc" {
		t.Errorf("expected NextCursor %q, got %q", "cursor-abc", got.NextCursor)
	}
	if !got.HasMore {
		t.Error("expected HasMore true")
	}
	if got.Status != 200 {
		t.Errorf("expected Status 200, got %d", got.Status)
	}
	if got.Code != "success" {
		t.Errorf("expected Code %q, got %q", "success", got.Code)
	}
	if got.Message != "OK" {
		t.Errorf("expected Message %q, got %q", "OK", got.Message)
	}
}

func TestFailureResponse_JSON(t *testing.T) {
	resp := FailureResponse{
		Object:         "error",
		Status:         400,
		Code:           "validation_error",
		Message:        "Invalid input",
		AdditionalData: map[string]interface{}{"field": "name"},
	}
	got := jsonRoundTrip(t, resp)
	if got.Object != "error" {
		t.Errorf("expected Object %q, got %q", "error", got.Object)
	}
	if got.Status != 400 {
		t.Errorf("expected Status 400, got %d", got.Status)
	}
	if got.Code != "validation_error" {
		t.Errorf("expected Code %q, got %q", "validation_error", got.Code)
	}
	if got.Message != "Invalid input" {
		t.Errorf("expected Message %q, got %q", "Invalid input", got.Message)
	}
	if got.AdditionalData["field"] != "name" {
		t.Errorf("expected AdditionalData field 'field' = 'name', got %v", got.AdditionalData["field"])
	}
}

func TestDatabase_JSON(t *testing.T) {
	db := Database{
		CommonObject: CommonObject{
			ID:             "db-123",
			Object:         "database",
			CreatedTime:    "2020-03-17T19:10:04.968Z",
			LastEditedTime: "2020-03-17T21:49:37.913Z",
		},
		Title: []Title{
			{Text: Text{Content: "My DB"}, PlainText: "My DB", Type: "text"},
		},
		Properties: map[string]Property{
			"Name": {ID: "title", Type: "title"},
		},
	}
	got := jsonRoundTrip(t, db)
	if got.ID != "db-123" {
		t.Errorf("expected ID %q, got %q", "db-123", got.ID)
	}
	if got.Object != "database" {
		t.Errorf("expected Object %q, got %q", "database", got.Object)
	}
	if len(got.Title) != 1 || got.Title[0].Text.Content != "My DB" {
		t.Errorf("unexpected Title: %+v", got.Title)
	}
	if _, ok := got.Properties["Name"]; !ok {
		t.Error("expected 'Name' property in Properties")
	}
}

func TestDatabaseResponseList_JSON(t *testing.T) {
	list := DatabaseResponseList{
		Response: Response{Object: "list", HasMore: false},
		Results: []Database{
			{CommonObject: CommonObject{ID: "db-1"}},
		},
	}
	got := jsonRoundTrip(t, list)
	if got.Object != "list" {
		t.Errorf("expected Object %q, got %q", "list", got.Object)
	}
	if len(got.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got.Results))
	}
	if got.Results[0].ID != "db-1" {
		t.Errorf("expected result ID %q, got %q", "db-1", got.Results[0].ID)
	}
}

func TestPage_JSON(t *testing.T) {
	pg := Page{
		CommonObject: CommonObject{
			ID:     "page-123",
			Object: "page",
		},
		Parent: Parent{
			Type:       ParentTypeDatabase,
			DatabaseID: "db-123",
		},
		Archived: true,
		IsLocked: true,
		Properties: map[string]ValueProperty{
			"Name": {
				Type:  ValuePropertyTypeTitle,
				Title: []Title{{PlainText: "My Page"}},
			},
		},
	}
	got := jsonRoundTrip(t, pg)
	if got.ID != "page-123" {
		t.Errorf("expected ID %q, got %q", "page-123", got.ID)
	}
	if got.Parent.Type != ParentTypeDatabase {
		t.Errorf("expected parent type %q, got %q", ParentTypeDatabase, got.Parent.Type)
	}
	if got.Parent.DatabaseID != "db-123" {
		t.Errorf("expected parent database ID %q, got %q", "db-123", got.Parent.DatabaseID)
	}
	if !got.Archived {
		t.Error("expected Archived true")
	}
	if !got.IsLocked {
		t.Error("expected IsLocked true")
	}
	if vp, ok := got.Properties["Name"]; !ok || vp.Type != ValuePropertyTypeTitle {
		t.Errorf("unexpected Properties: %+v", got.Properties)
	}
}

func TestPageResponseList_JSON(t *testing.T) {
	list := PageResponseList{
		Response: Response{Object: "list", HasMore: true, NextCursor: "cursor-1"},
		Results: []Page{
			{CommonObject: CommonObject{ID: "page-1"}},
			{CommonObject: CommonObject{ID: "page-2"}},
		},
	}
	got := jsonRoundTrip(t, list)
	if !got.HasMore {
		t.Error("expected HasMore true")
	}
	if got.NextCursor != "cursor-1" {
		t.Errorf("expected NextCursor %q, got %q", "cursor-1", got.NextCursor)
	}
	if len(got.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got.Results))
	}
}

func TestCommonObject_JSON(t *testing.T) {
	user := &User{Object: "user", ID: "user-1"}
	co := CommonObject{
		CreatedTime:    "2020-01-01T00:00:00.000Z",
		LastEditedTime: "2020-06-01T00:00:00.000Z",
		ID:             "obj-123",
		Object:         "database",
		InTrash:        true,
		CreatedBy:      user,
		LastEditedBy:   user,
	}
	got := jsonRoundTrip(t, co)
	if got.CreatedTime != "2020-01-01T00:00:00.000Z" {
		t.Errorf("unexpected CreatedTime: %q", got.CreatedTime)
	}
	if got.LastEditedTime != "2020-06-01T00:00:00.000Z" {
		t.Errorf("unexpected LastEditedTime: %q", got.LastEditedTime)
	}
	if got.ID != "obj-123" {
		t.Errorf("unexpected ID: %q", got.ID)
	}
	if !got.InTrash {
		t.Error("expected InTrash true")
	}
	if got.CreatedBy == nil || got.CreatedBy.ID != "user-1" {
		t.Errorf("unexpected CreatedBy: %+v", got.CreatedBy)
	}
	if got.LastEditedBy == nil || got.LastEditedBy.ID != "user-1" {
		t.Errorf("unexpected LastEditedBy: %+v", got.LastEditedBy)
	}
}

func TestUser_JSON(t *testing.T) {
	user := User{
		Object:    "user",
		ID:        "user-123",
		Type:      "person",
		Name:      "John Doe",
		AvatarURL: "https://example.com/avatar.png",
		Person:    &Person{Email: "john@example.com"},
	}
	got := jsonRoundTrip(t, user)
	if got.Object != "user" {
		t.Errorf("expected Object %q, got %q", "user", got.Object)
	}
	if got.ID != "user-123" {
		t.Errorf("expected ID %q, got %q", "user-123", got.ID)
	}
	if got.Type != "person" {
		t.Errorf("expected Type %q, got %q", "person", got.Type)
	}
	if got.Name != "John Doe" {
		t.Errorf("expected Name %q, got %q", "John Doe", got.Name)
	}
	if got.AvatarURL != "https://example.com/avatar.png" {
		t.Errorf("expected AvatarURL %q, got %q", "https://example.com/avatar.png", got.AvatarURL)
	}
	if got.Person == nil || got.Person.Email != "john@example.com" {
		t.Errorf("unexpected Person: %+v", got.Person)
	}
}

func TestUser_Bot_JSON(t *testing.T) {
	user := User{
		Object: "user",
		ID:     "bot-123",
		Type:   "bot",
		Name:   "My Bot",
		Bot: &Bot{
			Owner:         &Owner{Type: "workspace"},
			WorkspaceName: "Test Workspace",
		},
	}
	got := jsonRoundTrip(t, user)
	if got.Bot == nil {
		t.Fatal("expected Bot to be non-nil")
	}
	if got.Bot.Owner == nil || got.Bot.Owner.Type != "workspace" {
		t.Errorf("unexpected Bot.Owner: %+v", got.Bot.Owner)
	}
	if got.Bot.WorkspaceName != "Test Workspace" {
		t.Errorf("expected WorkspaceName %q, got %q", "Test Workspace", got.Bot.WorkspaceName)
	}
}

func TestOwner_WithUser_JSON(t *testing.T) {
	owner := Owner{
		Type: "user",
		User: &User{Object: "user", ID: "user-1"},
	}
	got := jsonRoundTrip(t, owner)
	if got.Type != "user" {
		t.Errorf("expected Type %q, got %q", "user", got.Type)
	}
	if got.User == nil || got.User.ID != "user-1" {
		t.Errorf("unexpected User: %+v", got.User)
	}
}

func TestUserResponseList_JSON(t *testing.T) {
	list := UserResponseList{
		Response: Response{Object: "list"},
		Results: []User{
			{Object: "user", ID: "user-1", Name: "Alice"},
			{Object: "user", ID: "user-2", Name: "Bob"},
		},
	}
	got := jsonRoundTrip(t, list)
	if len(got.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got.Results))
	}
	if got.Results[0].Name != "Alice" {
		t.Errorf("expected first user name %q, got %q", "Alice", got.Results[0].Name)
	}
}

func TestTitle_JSON(t *testing.T) {
	title := Title{
		Type:      "text",
		Text:      Text{Content: "Hello", Link: &Link{URL: "https://example.com"}},
		PlainText: "Hello",
		Href:      "https://example.com",
		Annotations: Annotation{
			Bold:   true,
			Italic: true,
			Color:  ColorBlue,
		},
	}
	got := jsonRoundTrip(t, title)
	if got.Type != "text" {
		t.Errorf("expected Type %q, got %q", "text", got.Type)
	}
	if got.Text.Content != "Hello" {
		t.Errorf("expected Content %q, got %q", "Hello", got.Text.Content)
	}
	if got.Text.Link == nil || got.Text.Link.URL != "https://example.com" {
		t.Errorf("unexpected Link: %+v", got.Text.Link)
	}
	if got.PlainText != "Hello" {
		t.Errorf("expected PlainText %q, got %q", "Hello", got.PlainText)
	}
	if got.Href != "https://example.com" {
		t.Errorf("expected Href %q, got %q", "https://example.com", got.Href)
	}
	if !got.Annotations.Bold || !got.Annotations.Italic {
		t.Error("expected Bold and Italic to be true")
	}
	if got.Annotations.Color != ColorBlue {
		t.Errorf("expected Color %q, got %q", ColorBlue, got.Annotations.Color)
	}
}

func TestText_JSON(t *testing.T) {
	txt := Text{Content: "content", Link: &Link{URL: "https://example.com"}}
	got := jsonRoundTrip(t, txt)
	if got.Content != "content" {
		t.Errorf("expected Content %q, got %q", "content", got.Content)
	}
	if got.Link == nil || got.Link.URL != "https://example.com" {
		t.Errorf("unexpected Link: %+v", got.Link)
	}
}

func TestAnnotation_JSON(t *testing.T) {
	ann := Annotation{
		Bold:          true,
		Italic:        true,
		Strikethrough: true,
		Underline:     true,
		Code:          true,
		Color:         ColorRed,
	}
	got := jsonRoundTrip(t, ann)
	if !got.Bold || !got.Italic || !got.Strikethrough || !got.Underline || !got.Code {
		t.Error("expected all formatting flags to be true")
	}
	if got.Color != ColorRed {
		t.Errorf("expected Color %q, got %q", ColorRed, got.Color)
	}
}

func TestProperty_Select(t *testing.T) {
	prop := Property{
		ID:   "prop-1",
		Type: "select",
		Select: Select{
			Options: []Option{
				{ID: "opt-1", Name: "Option A", Color: "blue"},
				{ID: "opt-2", Name: "Option B", Color: "red"},
			},
		},
	}
	got := jsonRoundTrip(t, prop)
	if got.ID != "prop-1" {
		t.Errorf("expected ID %q, got %q", "prop-1", got.ID)
	}
	if got.Type != "select" {
		t.Errorf("expected Type %q, got %q", "select", got.Type)
	}
	if len(got.Select.Options) != 2 {
		t.Fatalf("expected 2 options, got %d", len(got.Select.Options))
	}
	if got.Select.Options[0].Name != "Option A" {
		t.Errorf("expected option name %q, got %q", "Option A", got.Select.Options[0].Name)
	}
}

func TestProperty_MultiSelect(t *testing.T) {
	prop := Property{
		ID:   "prop-2",
		Type: "multi_select",
		MultiSelect: Select{
			Options: []Option{
				{ID: "opt-1", Name: "Tag A"},
			},
		},
	}
	got := jsonRoundTrip(t, prop)
	if len(got.MultiSelect.Options) != 1 {
		t.Fatalf("expected 1 option, got %d", len(got.MultiSelect.Options))
	}
}

func TestProperty_Number(t *testing.T) {
	prop := Property{
		ID:     "prop-3",
		Type:   "number",
		Number: Number{Format: "percent"},
	}
	got := jsonRoundTrip(t, prop)
	if got.Number.Format != "percent" {
		t.Errorf("expected format %q, got %q", "percent", got.Number.Format)
	}
}

func TestProperty_Status(t *testing.T) {
	prop := Property{
		ID:   "prop-4",
		Type: "status",
		Status: Status{
			Options: []StatusOption{
				{ID: "so-1", Name: "Not started", Color: "default"},
				{ID: "so-2", Name: "In progress", Color: "blue"},
			},
			Groups: []StatusGroup{
				{ID: "sg-1", Name: "To-do", Color: "gray", OptionIDs: []string{"so-1"}},
			},
		},
	}
	got := jsonRoundTrip(t, prop)
	if len(got.Status.Options) != 2 {
		t.Fatalf("expected 2 status options, got %d", len(got.Status.Options))
	}
	if got.Status.Options[0].Name != "Not started" {
		t.Errorf("expected option name %q, got %q", "Not started", got.Status.Options[0].Name)
	}
	if len(got.Status.Groups) != 1 {
		t.Fatalf("expected 1 status group, got %d", len(got.Status.Groups))
	}
	if got.Status.Groups[0].Name != "To-do" {
		t.Errorf("expected group name %q, got %q", "To-do", got.Status.Groups[0].Name)
	}
	if len(got.Status.Groups[0].OptionIDs) != 1 || got.Status.Groups[0].OptionIDs[0] != "so-1" {
		t.Errorf("unexpected OptionIDs: %v", got.Status.Groups[0].OptionIDs)
	}
}

func TestParent_Database(t *testing.T) {
	parent := Parent{Type: ParentTypeDatabase, DatabaseID: "db-123"}
	got := jsonRoundTrip(t, parent)
	if got.Type != ParentTypeDatabase {
		t.Errorf("expected type %q, got %q", ParentTypeDatabase, got.Type)
	}
	if got.DatabaseID != "db-123" {
		t.Errorf("expected DatabaseID %q, got %q", "db-123", got.DatabaseID)
	}
}

func TestParent_Page(t *testing.T) {
	parent := Parent{Type: ParentTypePage, PageID: "page-123"}
	got := jsonRoundTrip(t, parent)
	if got.Type != ParentTypePage {
		t.Errorf("expected type %q, got %q", ParentTypePage, got.Type)
	}
	if got.PageID != "page-123" {
		t.Errorf("expected PageID %q, got %q", "page-123", got.PageID)
	}
}

func TestParent_Workspace(t *testing.T) {
	parent := Parent{Type: ParentTypeWorkspace, Workspace: true}
	got := jsonRoundTrip(t, parent)
	if got.Type != ParentTypeWorkspace {
		t.Errorf("expected type %q, got %q", ParentTypeWorkspace, got.Type)
	}
	if !got.Workspace {
		t.Error("expected Workspace true")
	}
}

func TestValueProperty_Title(t *testing.T) {
	vp := ValueProperty{
		ID:    "vp-1",
		Type:  ValuePropertyTypeTitle,
		Title: []Title{{PlainText: "Test Title"}},
	}
	got := jsonRoundTrip(t, vp)
	if got.Type != ValuePropertyTypeTitle {
		t.Errorf("expected type %q, got %q", ValuePropertyTypeTitle, got.Type)
	}
	if len(got.Title) != 1 || got.Title[0].PlainText != "Test Title" {
		t.Errorf("unexpected Title: %+v", got.Title)
	}
}

func TestValueProperty_RichText(t *testing.T) {
	vp := ValueProperty{
		Type:     ValuePropertyTypeRichText,
		RichText: []Title{{PlainText: "Some text"}},
	}
	got := jsonRoundTrip(t, vp)
	if got.Type != ValuePropertyTypeRichText {
		t.Errorf("expected type %q, got %q", ValuePropertyTypeRichText, got.Type)
	}
	if len(got.RichText) != 1 {
		t.Fatalf("expected 1 rich text, got %d", len(got.RichText))
	}
}

func TestValueProperty_Number(t *testing.T) {
	vp := ValueProperty{
		Type:   ValuePropertyTypeNumber,
		Number: 42.5,
	}
	got := jsonRoundTrip(t, vp)
	if got.Number != 42.5 {
		t.Errorf("expected Number 42.5, got %f", got.Number)
	}
}

func TestValueProperty_Select(t *testing.T) {
	vp := ValueProperty{
		Type:   ValuePropertyTypeSelect,
		Select: &Option{ID: "opt-1", Name: "Selected", Color: "green"},
	}
	got := jsonRoundTrip(t, vp)
	if got.Select == nil || got.Select.Name != "Selected" {
		t.Errorf("unexpected Select: %+v", got.Select)
	}
}

func TestValueProperty_MultiSelect(t *testing.T) {
	vp := ValueProperty{
		Type: ValuePropertyTypeMultiSelect,
		MultiSelect: []Option{
			{ID: "opt-1", Name: "Tag1"},
			{ID: "opt-2", Name: "Tag2"},
		},
	}
	got := jsonRoundTrip(t, vp)
	if len(got.MultiSelect) != 2 {
		t.Fatalf("expected 2 multiselect options, got %d", len(got.MultiSelect))
	}
}

func TestValueProperty_Date(t *testing.T) {
	vp := ValueProperty{
		Type: ValuePropertyTypeDate,
		Date: &DateRange{Start: "2020-12-08", End: "2020-12-10"},
	}
	got := jsonRoundTrip(t, vp)
	if got.Date == nil {
		t.Fatal("expected Date to be non-nil")
	}
	if got.Date.Start != "2020-12-08" {
		t.Errorf("expected Start %q, got %q", "2020-12-08", got.Date.Start)
	}
	if got.Date.End != "2020-12-10" {
		t.Errorf("expected End %q, got %q", "2020-12-10", got.Date.End)
	}
}

func TestValueProperty_Checkbox(t *testing.T) {
	vp := ValueProperty{
		Type:     ValuePropertyTypeCheckbox,
		Checkbox: true,
	}
	got := jsonRoundTrip(t, vp)
	if !got.Checkbox {
		t.Error("expected Checkbox true")
	}
}

func TestValueProperty_URL(t *testing.T) {
	vp := ValueProperty{
		Type: ValuePropertyTypeURL,
		URL:  "https://example.com",
	}
	got := jsonRoundTrip(t, vp)
	if got.URL != "https://example.com" {
		t.Errorf("expected URL %q, got %q", "https://example.com", got.URL)
	}
}

func TestValueProperty_Relation(t *testing.T) {
	vp := ValueProperty{
		Type:     ValuePropertyTypeRelation,
		Relation: []Relation{{ID: "rel-1"}, {ID: "rel-2"}},
	}
	got := jsonRoundTrip(t, vp)
	if len(got.Relation) != 2 {
		t.Fatalf("expected 2 relations, got %d", len(got.Relation))
	}
	if got.Relation[0].ID != "rel-1" {
		t.Errorf("expected relation ID %q, got %q", "rel-1", got.Relation[0].ID)
	}
}

func TestValueProperty_Status(t *testing.T) {
	vp := ValueProperty{
		Type:   ValuePropertyTypeStatus,
		Status: &Option{ID: "opt-1", Name: "In progress", Color: "blue"},
	}
	got := jsonRoundTrip(t, vp)
	if got.Status == nil || got.Status.Name != "In progress" {
		t.Errorf("unexpected Status: %+v", got.Status)
	}
}

func TestSort_JSON(t *testing.T) {
	sort := Sort{
		Property:  "Name",
		Direction: &SortDirectionAscending,
	}
	got := jsonRoundTrip(t, sort)
	if got.Property != "Name" {
		t.Errorf("expected Property %q, got %q", "Name", got.Property)
	}
	if got.Direction == nil || *got.Direction != SortDirectionAscending {
		t.Errorf("unexpected Direction: %v", got.Direction)
	}
}

func TestSort_WithTimestamp(t *testing.T) {
	sort := Sort{
		Timestamp: &SortTimestampCreatedTime,
		Direction: &SortDirectionDescending,
	}
	got := jsonRoundTrip(t, sort)
	if got.Timestamp == nil || *got.Timestamp != SortTimestampCreatedTime {
		t.Errorf("unexpected Timestamp: %v", got.Timestamp)
	}
}

func TestQueryDB_JSON(t *testing.T) {
	query := QueryDB{
		Sorts: []Sort{
			{Property: "Name", Direction: &SortDirectionAscending},
		},
		Filter: &Filter{
			And: []PropertyFilter{
				{Property: "Status", Status: &StatusFilter{Equals: "Done"}},
			},
		},
		PageSize:    50,
		StartCursor: "cursor-123",
	}
	got := jsonRoundTrip(t, query)
	if len(got.Sorts) != 1 {
		t.Fatalf("expected 1 sort, got %d", len(got.Sorts))
	}
	if got.Filter == nil {
		t.Fatal("expected Filter to be non-nil")
	}
	if got.PageSize != 50 {
		t.Errorf("expected PageSize 50, got %d", got.PageSize)
	}
	if got.StartCursor != "cursor-123" {
		t.Errorf("expected StartCursor %q, got %q", "cursor-123", got.StartCursor)
	}
}

func TestFilter_JSON(t *testing.T) {
	filter := Filter{
		Or: []PropertyFilter{
			{Property: "Name", Title: &TextFilter{Contains: "test"}},
			{Property: "Email", Email: &TextFilter{Equals: "a@b.com"}},
		},
	}
	got := jsonRoundTrip(t, filter)
	if len(got.Or) != 2 {
		t.Fatalf("expected 2 or filters, got %d", len(got.Or))
	}
	if got.Or[0].Title == nil || got.Or[0].Title.Contains != "test" {
		t.Errorf("unexpected first filter: %+v", got.Or[0])
	}
}

func TestFilter_And(t *testing.T) {
	filter := Filter{
		And: []PropertyFilter{
			{Property: "Status", Select: &SelectFilter{Equals: "Done"}},
		},
	}
	got := jsonRoundTrip(t, filter)
	if len(got.And) != 1 {
		t.Fatalf("expected 1 and filter, got %d", len(got.And))
	}
}

func TestPropertyFilter_AllTypes(t *testing.T) {
	pf := PropertyFilter{
		Property:       "test",
		Title:          &TextFilter{Contains: "abc"},
		RichText:       &TextFilter{StartsWith: "def"},
		URL:            &TextFilter{Equals: "https://example.com"},
		Email:          &TextFilter{EndsWith: "@test.com"},
		Phone:          &TextFilter{DoesNotContain: "555"},
		Number:         &NumberFilter{GreaterThan: 10},
		Checkbox:       &CheckboxFilter{Equals: &BoolTrue},
		Select:         &SelectFilter{DoesNotEqual: "bad"},
		MultiSelect:    &MultiSelectFilter{Contains: "tag"},
		Date:           &DateFilter{Before: "2020-01-01"},
		CreatedTime:    &DateFilter{After: "2019-01-01"},
		LastEditedTime: &DateFilter{OnOrBefore: "2021-01-01"},
		Status:         &StatusFilter{Equals: "Done"},
	}
	got := jsonRoundTrip(t, pf)
	if got.Title == nil || got.Title.Contains != "abc" {
		t.Errorf("unexpected Title filter: %+v", got.Title)
	}
	if got.RichText == nil || got.RichText.StartsWith != "def" {
		t.Errorf("unexpected RichText filter: %+v", got.RichText)
	}
	if got.URL == nil || got.URL.Equals != "https://example.com" {
		t.Errorf("unexpected URL filter: %+v", got.URL)
	}
	if got.Email == nil || got.Email.EndsWith != "@test.com" {
		t.Errorf("unexpected Email filter: %+v", got.Email)
	}
	if got.Phone == nil || got.Phone.DoesNotContain != "555" {
		t.Errorf("unexpected Phone filter: %+v", got.Phone)
	}
	if got.Number == nil || got.Number.GreaterThan != 10 {
		t.Errorf("unexpected Number filter: %+v", got.Number)
	}
	if got.Checkbox == nil || got.Checkbox.Equals == nil || *got.Checkbox.Equals != true {
		t.Errorf("unexpected Checkbox filter: %+v", got.Checkbox)
	}
	if got.Select == nil || got.Select.DoesNotEqual != "bad" {
		t.Errorf("unexpected Select filter: %+v", got.Select)
	}
	if got.MultiSelect == nil || got.MultiSelect.Contains != "tag" {
		t.Errorf("unexpected MultiSelect filter: %+v", got.MultiSelect)
	}
	if got.Date == nil || got.Date.Before != "2020-01-01" {
		t.Errorf("unexpected Date filter: %+v", got.Date)
	}
	if got.CreatedTime == nil || got.CreatedTime.After != "2019-01-01" {
		t.Errorf("unexpected CreatedTime filter: %+v", got.CreatedTime)
	}
	if got.LastEditedTime == nil || got.LastEditedTime.OnOrBefore != "2021-01-01" {
		t.Errorf("unexpected LastEditedTime filter: %+v", got.LastEditedTime)
	}
	if got.Status == nil || got.Status.Equals != "Done" {
		t.Errorf("unexpected Status filter: %+v", got.Status)
	}
}

func TestCheckboxFilter_JSON(t *testing.T) {
	// Test with BoolTrue
	cf := CheckboxFilter{Equals: &BoolTrue}
	got := jsonRoundTrip(t, cf)
	if got.Equals == nil || *got.Equals != true {
		t.Errorf("expected Equals true, got %v", got.Equals)
	}

	// Test with BoolFalse
	cf2 := CheckboxFilter{DoesNotEqual: &BoolFalse}
	got2 := jsonRoundTrip(t, cf2)
	if got2.DoesNotEqual == nil || *got2.DoesNotEqual != false {
		t.Errorf("expected DoesNotEqual false, got %v", got2.DoesNotEqual)
	}
}

func TestTextFilter_JSON(t *testing.T) {
	tf := TextFilter{
		Equals:         "eq",
		DoesNotEqual:   "neq",
		Contains:       "cont",
		DoesNotContain: "ncont",
		StartsWith:     "sw",
		EndsWith:       "ew",
		IsEmpty:        true,
		IsNotEmpty:     true,
	}
	got := jsonRoundTrip(t, tf)
	if got.Equals != "eq" {
		t.Errorf("expected Equals %q, got %q", "eq", got.Equals)
	}
	if got.DoesNotEqual != "neq" {
		t.Errorf("expected DoesNotEqual %q, got %q", "neq", got.DoesNotEqual)
	}
	if got.Contains != "cont" {
		t.Errorf("expected Contains %q, got %q", "cont", got.Contains)
	}
	if got.DoesNotContain != "ncont" {
		t.Errorf("expected DoesNotContain %q, got %q", "ncont", got.DoesNotContain)
	}
	if got.StartsWith != "sw" {
		t.Errorf("expected StartsWith %q, got %q", "sw", got.StartsWith)
	}
	if got.EndsWith != "ew" {
		t.Errorf("expected EndsWith %q, got %q", "ew", got.EndsWith)
	}
	if !got.IsEmpty {
		t.Error("expected IsEmpty true")
	}
	if !got.IsNotEmpty {
		t.Error("expected IsNotEmpty true")
	}
}

func TestNumberFilter_JSON(t *testing.T) {
	nf := NumberFilter{
		Equals:               5,
		DoesNotEqual:         10,
		GreaterThan:          3,
		LessThan:             20,
		GreaterThanOrEqualTo: 5,
		LessThanOrEqualTo:    20,
		IsEmpty:              true,
		IsNotEmpty:           true,
	}
	got := jsonRoundTrip(t, nf)
	if got.Equals != 5 {
		t.Errorf("expected Equals 5, got %d", got.Equals)
	}
	if got.DoesNotEqual != 10 {
		t.Errorf("expected DoesNotEqual 10, got %d", got.DoesNotEqual)
	}
	if got.GreaterThan != 3 {
		t.Errorf("expected GreaterThan 3, got %d", got.GreaterThan)
	}
	if got.LessThan != 20 {
		t.Errorf("expected LessThan 20, got %d", got.LessThan)
	}
	if got.GreaterThanOrEqualTo != 5 {
		t.Errorf("expected GreaterThanOrEqualTo 5, got %d", got.GreaterThanOrEqualTo)
	}
	if got.LessThanOrEqualTo != 20 {
		t.Errorf("expected LessThanOrEqualTo 20, got %d", got.LessThanOrEqualTo)
	}
	if !got.IsEmpty {
		t.Error("expected IsEmpty true")
	}
	if !got.IsNotEmpty {
		t.Error("expected IsNotEmpty true")
	}
}

func TestSelectFilter_JSON(t *testing.T) {
	sf := SelectFilter{
		Equals:       "option1",
		DoesNotEqual: "option2",
		IsEmpty:      true,
		IsNotEmpty:   true,
	}
	got := jsonRoundTrip(t, sf)
	if got.Equals != "option1" {
		t.Errorf("expected Equals %q, got %q", "option1", got.Equals)
	}
	if got.DoesNotEqual != "option2" {
		t.Errorf("expected DoesNotEqual %q, got %q", "option2", got.DoesNotEqual)
	}
	if !got.IsEmpty {
		t.Error("expected IsEmpty true")
	}
	if !got.IsNotEmpty {
		t.Error("expected IsNotEmpty true")
	}
}

func TestMultiSelectFilter_JSON(t *testing.T) {
	msf := MultiSelectFilter{
		Contains:       "tag1",
		DoesNotContain: "tag2",
		IsEmpty:        true,
		IsNotEmpty:     true,
	}
	got := jsonRoundTrip(t, msf)
	if got.Contains != "tag1" {
		t.Errorf("expected Contains %q, got %q", "tag1", got.Contains)
	}
	if got.DoesNotContain != "tag2" {
		t.Errorf("expected DoesNotContain %q, got %q", "tag2", got.DoesNotContain)
	}
}

func TestDateFilter_JSON(t *testing.T) {
	onOrAfter := "2021-01-01"
	pastWeek := ""
	pastMonth := ""
	pastYear := ""
	nextWeek := ""
	nextMonth := ""
	nextYear := ""
	df := DateFilter{
		Equals:     "2020-01-01",
		Before:     "2020-06-01",
		After:      "2019-06-01",
		OnOrBefore: "2020-12-31",
		IsEmpty:    true,
		IsNotEmpty: true,
		OnOrAfter:  &onOrAfter,
		PastWeek:   &pastWeek,
		PastMonth:  &pastMonth,
		PastYear:   &pastYear,
		NextWeek:   &nextWeek,
		NextMonth:  &nextMonth,
		NextYear:   &nextYear,
	}
	got := jsonRoundTrip(t, df)
	if got.Equals != "2020-01-01" {
		t.Errorf("expected Equals %q, got %q", "2020-01-01", got.Equals)
	}
	if got.Before != "2020-06-01" {
		t.Errorf("expected Before %q, got %q", "2020-06-01", got.Before)
	}
	if got.After != "2019-06-01" {
		t.Errorf("expected After %q, got %q", "2019-06-01", got.After)
	}
	if got.OnOrBefore != "2020-12-31" {
		t.Errorf("expected OnOrBefore %q, got %q", "2020-12-31", got.OnOrBefore)
	}
	if !got.IsEmpty {
		t.Error("expected IsEmpty true")
	}
	if !got.IsNotEmpty {
		t.Error("expected IsNotEmpty true")
	}
	if got.OnOrAfter == nil || *got.OnOrAfter != "2021-01-01" {
		t.Errorf("expected OnOrAfter %q, got %v", "2021-01-01", got.OnOrAfter)
	}
	if got.PastWeek == nil {
		t.Error("expected PastWeek to be non-nil")
	}
	if got.PastMonth == nil {
		t.Error("expected PastMonth to be non-nil")
	}
	if got.PastYear == nil {
		t.Error("expected PastYear to be non-nil")
	}
	if got.NextWeek == nil {
		t.Error("expected NextWeek to be non-nil")
	}
	if got.NextMonth == nil {
		t.Error("expected NextMonth to be non-nil")
	}
	if got.NextYear == nil {
		t.Error("expected NextYear to be non-nil")
	}
}

func TestStatusFilter_JSON(t *testing.T) {
	sf := StatusFilter{
		Equals:       "Done",
		DoesNotEqual: "Not started",
		IsEmpty:      true,
		IsNotEmpty:   true,
	}
	got := jsonRoundTrip(t, sf)
	if got.Equals != "Done" {
		t.Errorf("expected Equals %q, got %q", "Done", got.Equals)
	}
	if got.DoesNotEqual != "Not started" {
		t.Errorf("expected DoesNotEqual %q, got %q", "Not started", got.DoesNotEqual)
	}
	if !got.IsEmpty {
		t.Error("expected IsEmpty true")
	}
	if !got.IsNotEmpty {
		t.Error("expected IsNotEmpty true")
	}
}

func TestDataSource_JSON(t *testing.T) {
	ds := DataSource{
		CommonObject: CommonObject{ID: "ds-123", Object: "data_source"},
		Title:        []Title{{PlainText: "Data Source 1"}},
		Properties:   map[string]Property{"Col": {ID: "c1", Type: "title"}},
		Parent:       Parent{Type: ParentTypeDatabase, DatabaseID: "db-123"},
	}
	got := jsonRoundTrip(t, ds)
	if got.ID != "ds-123" {
		t.Errorf("expected ID %q, got %q", "ds-123", got.ID)
	}
	if len(got.Title) != 1 || got.Title[0].PlainText != "Data Source 1" {
		t.Errorf("unexpected Title: %+v", got.Title)
	}
	if _, ok := got.Properties["Col"]; !ok {
		t.Error("expected 'Col' property")
	}
	if got.Parent.DatabaseID != "db-123" {
		t.Errorf("expected parent database ID %q, got %q", "db-123", got.Parent.DatabaseID)
	}
}

func TestDataSourceResponseList_JSON(t *testing.T) {
	list := DataSourceResponseList{
		Response: Response{Object: "list"},
		Results: []DataSource{
			{CommonObject: CommonObject{ID: "ds-1"}},
		},
	}
	got := jsonRoundTrip(t, list)
	if len(got.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got.Results))
	}
	if got.Results[0].ID != "ds-1" {
		t.Errorf("expected ID %q, got %q", "ds-1", got.Results[0].ID)
	}
}

func TestMovePageRequest_JSON(t *testing.T) {
	req := MovePageRequest{
		Parent:   Parent{Type: ParentTypePage, PageID: "page-456"},
		Position: &PagePosition{Type: "after", PageID: "sibling-789"},
	}
	got := jsonRoundTrip(t, req)
	if got.Parent.Type != ParentTypePage {
		t.Errorf("expected parent type %q, got %q", ParentTypePage, got.Parent.Type)
	}
	if got.Parent.PageID != "page-456" {
		t.Errorf("expected parent page ID %q, got %q", "page-456", got.Parent.PageID)
	}
	if got.Position == nil {
		t.Fatal("expected Position to be non-nil")
	}
	if got.Position.Type != "after" {
		t.Errorf("expected position type %q, got %q", "after", got.Position.Type)
	}
	if got.Position.PageID != "sibling-789" {
		t.Errorf("expected position page ID %q, got %q", "sibling-789", got.Position.PageID)
	}
}

// Variable verification tests

func TestColorVariables(t *testing.T) {
	colors := map[string]Color{
		"default":           ColorDefault,
		"gray":              ColorGray,
		"brown":             ColorBrown,
		"orange":            ColorOrange,
		"yellow":            ColorYellow,
		"green":             ColorGreen,
		"blue":              ColorBlue,
		"purple":            ColorPurple,
		"pink":              ColorPink,
		"red":               ColorRed,
		"gray_background":   ColorGrayBackground,
		"brown_background":  ColorBrownBackground,
		"orange_background": ColorOrangeBackground,
		"yellow_background": ColorYellowBackground,
		"green_background":  ColorGreenBackground,
		"blue_background":   ColorBlueBackground,
		"purple_background": ColorPurpleBackground,
		"pink_background":   ColorPinkBackground,
		"red_background":    ColorRedBackground,
	}
	for expected, got := range colors {
		if string(got) != expected {
			t.Errorf("expected color %q, got %q", expected, got)
		}
	}
}

func TestParentTypeVariables(t *testing.T) {
	tests := map[string]ParentType{
		"database_id": ParentTypeDatabase,
		"page_id":     ParentTypePage,
		"workspace":   ParentTypeWorkspace,
	}
	for expected, got := range tests {
		if string(got) != expected {
			t.Errorf("expected parent type %q, got %q", expected, got)
		}
	}
}

func TestValuePropertyTypeVariables(t *testing.T) {
	tests := map[string]ValuePropertyType{
		"rich_text":        ValuePropertyTypeRichText,
		"number":           ValuePropertyTypeNumber,
		"select":           ValuePropertyTypeSelect,
		"multi_select":     ValuePropertyTypeMultiSelect,
		"date":             ValuePropertyTypeDate,
		"formula":          ValuePropertyTypeFormula,
		"relation":         ValuePropertyTypeRelation,
		"rollup":           ValuePropertyTypeRollup,
		"title":            ValuePropertyTypeTitle,
		"people":           ValuePropertyTypePeople,
		"files":            ValuePropertyTypeFiles,
		"checkbox":         ValuePropertyTypeCheckbox,
		"url":              ValuePropertyTypeURL,
		"email":            ValuePropertyTypeEmail,
		"phone_number":     ValuePropertyTypePhoneNumber,
		"created_time":     ValuePropertyTypeCreatedTime,
		"created_by":       ValuePropertyTypeCreatedBy,
		"last_edited_time": ValuePropertyTypeLastEditedTime,
		"last_edited_by":   ValuePropertyTypeLastEditedBy,
		"status":           ValuePropertyTypeStatus,
	}
	for expected, got := range tests {
		if string(got) != expected {
			t.Errorf("expected value property type %q, got %q", expected, got)
		}
	}
}

func TestSortDirectionVariables(t *testing.T) {
	if string(SortDirectionAscending) != "ascending" {
		t.Errorf("expected %q, got %q", "ascending", SortDirectionAscending)
	}
	if string(SortDirectionDescending) != "descending" {
		t.Errorf("expected %q, got %q", "descending", SortDirectionDescending)
	}
}

func TestSortTimestampVariables(t *testing.T) {
	if string(SortTimestampCreatedTime) != "created_time" {
		t.Errorf("expected %q, got %q", "created_time", SortTimestampCreatedTime)
	}
	if string(SortTimestampLastEditedTime) != "last_edited_time" {
		t.Errorf("expected %q, got %q", "last_edited_time", SortTimestampLastEditedTime)
	}
}

func TestBoolVars(t *testing.T) {
	if BoolFalse != false {
		t.Errorf("expected BoolFalse to be false")
	}
	if BoolTrue != true {
		t.Errorf("expected BoolTrue to be true")
	}
}
