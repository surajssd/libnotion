# LibNotion

This repository contains the required API types to parse Notion databases, pages, responses, etc.

## Usage

Create a page in the database using the following `Page` struct:

```go
api.Page{
	Parent: api.Parent{
		Type:       api.ParentTypeDatabase,
		DatabaseID: "48f8fee9cd794180bc2fec0398253067",
	},
	Properties: map[string]api.ValueProperty{
		"Name": {
			Type: api.ValuePropertyTypeTitle,
			Title: []api.Title{
				{Type: "text", Text: api.Text{Content: "Tuscan Kale"}},
			},
		},
		"Date": {
			Type: api.ValuePropertyTypeDate,
			Date: &api.DateRange{Start: "2021-07-11"},
		},
		"DB": {
			Type:     api.ValuePropertyTypeRelation,
			Relation: []api.Relation{{ID: "897e5a76ae524b489fdfe71f5945d1af"}},
		},
	},
},
```
