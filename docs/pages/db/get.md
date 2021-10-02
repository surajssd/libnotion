# Get database entries (pages)

This page shows how to retrieve one or multiple entries (pages) from the database.

## Database Schema

You can see that the database has six columns and five entries of the books.

![](./img/add-page-after.png)

## Code

This code retrieves names of all the books, in an ascending order, that have either `Category=Non-fiction` or `Sub Category=History`. See the code here:

<details open>

```go
package main

import (
	"fmt"
	"github.com/surajssd/libnotion/api"
	"github.com/surajssd/libnotion/pkg/rest"
)

func main() {
	// Get this ID before hand either using an API, from a config file or by
	// hardcoding it.
	booksDBID := ""

	// Create a new Notion client.
	nc := rest.NewNotionClient(rest.WithSecretToken(token))

	books, err := nc.QueryDatabase(booksDBID, &api.QueryDB{
		Filter: &api.Filter{
			// Use OR among all the conditions provided.
			Or: []api.PropertyFilter{
				// Find all the books of "Category" equal to "Non-fiction" OR of
				// "Sub Category" equal to "History".
				{Property: "Category", Select: &api.SelectFilter{Equals: "Non-fiction"}},
				{Property: "Sub Category", MultiSelect: &api.MultiSelectFilter{Contains: "History"}},
			},
		},
		Sorts: []api.Sort{
			// Sort the output in an ascending order of their "Name" column values.
			{Property: "Name", Direction: &api.SortDirectionAscending},
		},
	})
	if err != nil {
		panic(err)
	}

	for _, book := range books {
		fmt.Println(book.Properties["Name"].Title[0].Text.Content)
	}

}
```

</details>

## In Action

```bash
$ go run main.go
Chava
Getting Things Done: The Art of Stress-Free Productivity
Sapiens
```

The first book, viz. "Chava" has `Category="Literature - Marathi"` and `Sub Category="History"` and the later being our selection criteria this book is part of the results. The second book and the third book are part of the returned as a result of our condition of `Category="Non-fiction"`.