package api

type Response struct {
	Object     string `json:"object,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more,omitempty"`
	Status     int    `json:"status,omitempty"`
	Code       string `json:"code,omitempty"`
	Message    string `json:"message,omitempty"`
}

type FailureResponse struct {
	Object  string `json:"object,omitempty"`
	Status  int    `json:"status,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type DatabaseResponseList struct {
	Response
	Results []Database `json:"results,omitempty"`
}

type Database struct {
	// TODO: Convert to time.Time type.
	CreatedTime    string `json:"created_time,omitempty"`
	LastEditedTime string `json:"last_edited_time,omitempty"`

	ID         string              `json:"id,omitempty"`
	Object     string              `json:"object,omitempty"`
	Title      []Title             `json:"title,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
}

type PageResponseList struct {
	Response
	Results []Page `json:"results,omitempty"`
}

type Page struct {
	// TODO: Convert to time.Time type.
	CreatedTime    string `json:"created_time,omitempty"`
	LastEditedTime string `json:"last_edited_time,omitempty"`

	ID         string                   `json:"id,omitempty"`
	Object     string                   `json:"object,omitempty"`
	Parent     Parent                   `json:"parent,omitempty"`
	Archived   bool                     `json:"archived,omitempty"`
	Properties map[string]ValueProperty `json:"properties,omitempty"`
}

type Title struct {
	Type        string     `json:"type,omitempty"`
	Text        Text       `json:"text,omitempty"`
	Annotations Annotation `json:"annotations,omitempty"`
	PlainText   string     `json:"plain_text,omitempty"`
	Href        string     `json:"href,omitempty"`
}

type Text struct {
	Content string `json:"content,omitempty"`
	Link    string `json:"link,omitempty"`
}

type Annotation struct {
	Bold          bool   `json:"bold,omitempty"`
	Italic        bool   `json:"italic,omitempty"`
	Strikethrough bool   `json:"strikethrough,omitempty"`
	Underline     bool   `json:"underline,omitempty"`
	Code          bool   `json:"code,omitempty"`
	Color         string `json:"color,omitempty"`
}

type Property struct {
	ID          string   `json:"id,omitempty"`
	Type        string   `json:"type,omitempty"`
	MultiSelect Select   `json:"multi_select,omitempty"`
	Number      Number   `json:"number,omitempty"`
	Select      Select   `json:"select,omitempty"`
	Date        Date     `json:"date,omitempty"`
	Checkbox    Checkbox `json:"checkbox,omitempty"`
}

type Select struct {
	Options []Option `json:"options,omitempty"`
}

type Option struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}

type Number struct {
	Format string `json:"format,omitempty"`
}

type Date struct{}

type Checkbox struct{}

type Parent struct {
	Type       ParentType `json:"type,omitempty"`
	DatabaseID string     `json:"database_id,omitempty"`
}

type ParentType string

var (
	ParentTypeDatabase = ParentType("database_id")
)

type ValueProperty struct {
	ID          string            `json:"id,omitempty"`
	Type        ValuePropertyType `json:"type,omitempty"`
	MultiSelect []Option          `json:"multi_select,omitempty"`
	Number      float64           `json:"number,omitempty"`
	Date        *DateRange        `json:"date,omitempty"`
	Checkbox    bool              `json:"checkbox,omitempty"`
	Select      *Option           `json:"select,omitempty"`
	Title       []Title           `json:"title,omitempty"`
	Relation    []Relation        `json:"relation,omitempty"`
}

type Relation struct {
	ID string `json:"id,omitempty"`
}

type ValuePropertyType string

var (
	ValuePropertyTypeNumber         = ValuePropertyType("number")
	ValuePropertyTypeSelect         = ValuePropertyType("select")
	ValuePropertyTypeMultiSelect    = ValuePropertyType("multi_select")
	ValuePropertyTypeDate           = ValuePropertyType("date")
	ValuePropertyTypeFormula        = ValuePropertyType("formula")
	ValuePropertyTypeRelation       = ValuePropertyType("relation")
	ValuePropertyTypeRollup         = ValuePropertyType("rollup")
	ValuePropertyTypeTitle          = ValuePropertyType("title")
	ValuePropertyTypePeople         = ValuePropertyType("people")
	ValuePropertyTypeFiles          = ValuePropertyType("files")
	ValuePropertyTypeCheckbox       = ValuePropertyType("checkbox")
	ValuePropertyTypeURL            = ValuePropertyType("url")
	ValuePropertyTypeEmail          = ValuePropertyType("email")
	ValuePropertyTypePhoneNumber    = ValuePropertyType("phone_number")
	ValuePropertyTypeCreatedTime    = ValuePropertyType("created_time")
	ValuePropertyTypeCreatedBy      = ValuePropertyType("created_by")
	ValuePropertyTypeLastEditedTime = ValuePropertyType("last_edited_time")
	ValuePropertyTypeLastEditedBy   = ValuePropertyType("last_edited_by")
	ValuePropertyTypeRichText       = ValuePropertyType("rich_text")
)

type DateRange struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

type SortDirection string

var (
	SortDirectionAscending  = SortDirection("ascending")
	SortDirectionDescending = SortDirection("descending")
)

type SortTimestamp string

var (
	SortTimestampCreatedTime    = SortTimestamp("created_time")
	SortTimestampLastEditedTime = SortTimestamp("last_edited_time")
)

type Sort struct {
	Property  string         `json:"property,omitempty"`
	Direction *SortDirection `json:"direction,omitempty"`
	Timestamp *SortTimestamp `json:"timestamp,omitempty"`
}

type QueryDB struct {
	Sorts       []Sort  `json:"sorts,omitempty"`
	Filter      *Filter `json:"filter,omitempty"`
	PageSize    int     `json:"page_size,omitempty"`
	StartCursor string  `json:"start_cursor,omitempty"`
}

type Filter struct {
	Or  []PropertyFilter `json:"or,omitempty"`
	And []PropertyFilter `json:"and,omitempty"`
}

type PropertyFilter struct {
	Property       string             `json:"property,omitempty"`
	Title          *TextFilter        `json:"title,omitempty"`
	RichText       *TextFilter        `json:"rich_text,omitempty"`
	URL            *TextFilter        `json:"url,omitempty"`
	Email          *TextFilter        `json:"email,omitempty"`
	Phone          *TextFilter        `json:"phone,omitempty"`
	Number         *NumberFilter      `json:"number,omitempty"`
	Checkbox       *CheckboxFilter    `json:"checkbox,omitempty"`
	Select         *SelectFilter      `json:"select,omitempty"`
	MultiSelect    *MultiSelectFilter `json:"multi_select,omitempty"`
	Date           *DateFilter        `json:"date,omitempty"`
	CreatedTime    *DateFilter        `json:"created_time,omitempty"`
	LastEditedTime *DateFilter        `json:"last_edited_time,omitempty"`
}

var (
	BoolFalse = false
	BoolTrue  = true
)

type TextFilter struct {
	Equals         string `json:"equals,omitempty"`
	DoesNotEqual   string `json:"does_not_equal,omitempty"`
	Contains       string `json:"contains,omitempty"`
	DoesNotContain string `json:"does_not_contain,omitempty"`
	StartsWith     string `json:"starts_with,omitempty"`
	EndsWith       string `json:"ends_with,omitempty"`
	IsEmpty        bool   `json:"is_empty,omitempty"`
	IsNotEmpty     bool   `json:"is_not_empty,omitempty"`
}

type NumberFilter struct {
	Equals               int  `json:"equals,omitempty"`
	DoesNotEqual         int  `json:"does_not_equal,omitempty"`
	GreaterThan          int  `json:"greater_than,omitempty"`
	LessThan             int  `json:"less_than,omitempty"`
	GreaterThanOrEqualTo int  `json:"greater_than_or_equal_to,omitempty"`
	LessThanOrEqualTo    int  `json:"less_than_or_equal_to,omitempty"`
	IsEmpty              bool `json:"is_empty,omitempty"`
	IsNotEmpty           bool `json:"is_not_empty,omitempty"`
}

type CheckboxFilter struct {
	Equals       *bool `json:"equals,omitempty"`
	DoesNotEqual *bool `json:"does_not_equal,omitempty"`
}

type SelectFilter struct {
	Equals       string `json:"equals,omitempty"`
	DoesNotEqual string `json:"does_not_equal,omitempty"`
	IsEmpty      bool   `json:"is_empty,omitempty"`
	IsNotEmpty   bool   `json:"is_not_empty,omitempty"`
}

type MultiSelectFilter struct {
	Contains       string `json:"contains,omitempty"`
	DoesNotContain string `json:"does_not_contain,omitempty"`
	IsEmpty        bool   `json:"is_empty,omitempty"`
	IsNotEmpty     bool   `json:"is_not_empty,omitempty"`
}

type DateFilter struct {
	Equals     string  `json:"equals,omitempty"`
	Before     string  `json:"before,omitempty"`
	After      string  `json:"after,omitempty"`
	OnOrBefore string  `json:"on_or_before,omitempty"`
	IsEmpty    bool    `json:"is_empty,omitempty"`
	IsNotEmpty bool    `json:"is_not_empty,omitempty"`
	OnOrAfter  *string `json:"on_or_after,omitempty"`
	PastWeek   *string `json:"past_week,omitempty"`
	PastMonth  *string `json:"past_month,omitempty"`
	PastYear   *string `json:"past_year,omitempty"`
	NextWeek   *string `json:"next_week,omitempty"`
	NextMonth  *string `json:"next_month,omitempty"`
	NextYear   *string `json:"next_year,omitempty"`
}
