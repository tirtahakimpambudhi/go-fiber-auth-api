package request

import "fmt"

// Struct for pagination parameters in requests.
type Page struct {
	Size   int    `query:"size" validate:"required,gt=0"`      // Page size (required, greater than 0)
	Before string `query:"before" validate:"omitempty,len=27"` // Optional pagination cursor for before
	After  string `query:"after" validate:"omitempty,len=27"`  // Optional pagination cursor for after
}

// Generates query parameters string for pagination.
func (page Page) GetQueryParams() string {
	queryParams := "?size=" + fmt.Sprint(page.Size) // Start query with page size
	if page.Before != "" {
		queryParams += "&before=" + page.Before // Add before cursor if present
	}
	if page.After != "" {
		queryParams += "&after=" + page.After // Add after cursor if present
	}
	return queryParams
}
