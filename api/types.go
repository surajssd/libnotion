package api

// Response returned by the Notion API when the status code is 200.
type Response struct {
	Object     string `json:"object,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more,omitempty"`
	Status     int    `json:"status,omitempty"`
	Code       string `json:"code,omitempty"`
	Message    string `json:"message,omitempty"`
}

// FailureResponse is used to parse the response if the returned status code from the API is
// non-200.
type FailureResponse struct {
	Object  string `json:"object,omitempty"`
	Status  int    `json:"status,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// DatabaseResponseList is used to parse the response when querying database endpoint.
type DatabaseResponseList struct {
	Response
	Results []Database `json:"results,omitempty"`
}

// Database is used for storing database metadata like the columns in a database.
type Database struct {
	// TODO: Convert to time.Time type.
	// Date and time when this database was created. Formatted as an ISO 8601 date time string. e.g. "2020-03-17T19:10:04.968Z"
	CreatedTime string `json:"created_time,omitempty"`
	// Date and time when this database was updated. Formatted as an ISO 8601 date time string. e.g. "2020-03-17T21:49:37.913Z"
	LastEditedTime string `json:"last_edited_time,omitempty"`

	// Unique identifier for the database.
	ID string `json:"id,omitempty"`

	// Always "database".
	Object string `json:"object,omitempty"`

	// Title of database as it appears in Notion. An array of rich text objects.
	Title []Title `json:"title,omitempty"`

	// Property schema of database. This corresponds with the columns in the database. The keys are
	// the names of properties as they appear in Notion and the values are property schema objects.
	Properties map[string]Property `json:"properties,omitempty"`
}

// PageResponseList is used to parse the response when querying pages endpoint.
type PageResponseList struct {
	Response
	Results []Page `json:"results,omitempty"`
}

// Page is used for storing all the information related to a page. If the page is a part of a
// database then it includes properties as well.
type Page struct {
	// TODO: Convert to time.Time type.
	// Date and time when this page was created. Formatted as an ISO 8601 date time string.
	CreatedTime string `json:"created_time,omitempty"`
	// Date and time when this page was updated. Formatted as an ISO 8601 date time string.
	LastEditedTime string `json:"last_edited_time,omitempty"`

	// Unique identifier of the page.
	ID string `json:"id,omitempty"`

	// Always "page".
	Object string `json:"object,omitempty"`

	// The parent of this page. Can be a database, page, or workspace.
	Parent Parent `json:"parent,omitempty"`

	// The archived status of the page.
	Archived bool `json:"archived,omitempty"`

	// Property values of this page. If parent.type is "page_id" or "workspace", then the only valid
	// key is title. If parent.type is "database_id", then the keys and values of this field are
	// determined by the properties of the database this page belongs to.
	//
	// key string: Name of a property as it appears in Notion.
	// value object: A Property value object.
	Properties map[string]ValueProperty `json:"properties,omitempty"`
}

// Title of database as it appears in Notion.
type Title struct {
	Type        string     `json:"type,omitempty"`
	Text        Text       `json:"text,omitempty"`
	Annotations Annotation `json:"annotations,omitempty"`
	PlainText   string     `json:"plain_text,omitempty"`
	Href        string     `json:"href,omitempty"`
}

type Text struct {
	// Text content. This field contains the actual content of your text and is probably the field
	// you'll use most often.
	Content string `json:"content,omitempty"`

	// Any inline link in this text.
	Link string `json:"link,omitempty"`
}

type Annotation struct {
	// Whether the text is bolded.
	Bold bool `json:"bold,omitempty"`
	// Whether the text is italicized.
	Italic bool `json:"italic,omitempty"`
	// Whether the text is struck through.
	Strikethrough bool `json:"strikethrough,omitempty"`
	// Whether the text is underlined.
	Underline bool `json:"underline,omitempty"`
	// Whether the text is code style.
	Code bool `json:"code,omitempty"`
	// Color of the text.
	Color Color `json:"color,omitempty"`
}

type Color string

var (
	ColorDefault          = Color("default")
	ColorGray             = Color("gray")
	ColorBrown            = Color("brown")
	ColorOrange           = Color("orange")
	ColorYellow           = Color("yellow")
	ColorGreen            = Color("green")
	ColorBlue             = Color("blue")
	ColorPurple           = Color("purple")
	ColorPink             = Color("pink")
	ColorRed              = Color("red")
	ColorGrayBackground   = Color("gray_background")
	ColorBrownBackground  = Color("brown_background")
	ColorOrangeBackground = Color("orange_background")
	ColorYellowBackground = Color("yellow_background")
	ColorGreenBackground  = Color("green_background")
	ColorBlueBackground   = Color("blue_background")
	ColorPurpleBackground = Color("purple_background")
	ColorPinkBackground   = Color("pink_background")
	ColorRedBackground    = Color("red_background")
)

type Property struct {
	ID          string   `json:"id,omitempty"`
	Type        string   `json:"type,omitempty"`
	MultiSelect Select   `json:"multi_select,omitempty"`
	Number      Number   `json:"number,omitempty"`
	Select      Select   `json:"select,omitempty"`
	Date        Date     `json:"date,omitempty"`
	Checkbox    Checkbox `json:"checkbox,omitempty"`
}

// Sorted list of options available for this property.
type Select struct {
	Options []Option `json:"options,omitempty"`
}

type Option struct {
	ID string `json:"id,omitempty"`

	// Name of the option as it appears in Notion.
	Name string `json:"name"`

	// Color of the option. Possible values include: default, gray, brown, orange, yellow, green,
	// blue, purple, pink, red.
	Color string `json:"color,omitempty"`
}

type Number struct {
	Format string `json:"format,omitempty"`
}

// Date database property schema objects have no additional configuration within the date property.
type Date struct{}

// Checkbox database property schema objects have no additional configuration within the checkbox
// property.
type Checkbox struct{}

type Parent struct {
	// The paret type could be page, database or workspace.
	Type ParentType `json:"type,omitempty"`

	// The ID of the database that this page belongs to.
	DatabaseID string `json:"database_id,omitempty"`

	// The ID of the page that this page belongs to.
	PageID string `json:"page_id,omitempty"`

	// True if the parent is a workspace.
	Workspace bool `json:"workspace,omitempty"`
}

type ParentType string

var (
	// If the parent is database.
	ParentTypeDatabase = ParentType("database_id")

	// If the parent is a page.
	ParentTypePage = ParentType("page_id")

	// A page with a workspace parent is a top-level page within a Notion workspace.
	ParentTypeWorkspace = ParentType("workspace")
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
	ValuePropertyTypeRichText       = ValuePropertyType("rich_text")
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
)

// DateRange is used to specify span of date between start and end.
type DateRange struct {
	// An ISO 8601 format date, with optional time. e.g. "2020-12-08T12:00:00Z".
	Start string `json:"start,omitempty"`

	// An ISO 8601 formatted date, with optional time. Represents the end of a date range.
	// If null, this property's date value is not a range. e.g. "2020-12-08T12:00:00Z".
	End string `json:"end,omitempty"`
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

// QueryDB lets you filter the resulting pages as per the given criteria.
type QueryDB struct {
	// When supplied, orders the results based on the provided sort criteria.
	Sorts []Sort `json:"sorts,omitempty"`

	// When supplied, limits which pages are returned based on the filter conditions.
	Filter *Filter `json:"filter,omitempty"`

	// The number of items from the full list desired in the response. Maximum: 100
	PageSize int `json:"page_size,omitempty"`

	// When supplied, returns a page of results starting after the cursor provided. If not supplied,
	// this endpoint will return the first page of results.
	StartCursor string `json:"start_cursor,omitempty"`
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
