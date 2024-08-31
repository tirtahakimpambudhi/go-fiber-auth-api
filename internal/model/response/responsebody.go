package response

import (
	"fmt"
	"strings"
)

type Standard struct {
	Meta   map[string]any `json:"meta,omitempty"`
	Status int            `json:"status"`
	Code   string         `json:"code"`
	Data   any            `json:"data"`
}

type Error struct {
	Meta   map[string]any `json:"meta,omitempty"`
	Title  string         `json:"title"`
	Status int            `json:"status"`
	Code   string         `json:"code"`
	Detail string         `json:"detail"`
}

func (err Error) Error() string {
	return fmt.Sprintf("%s : %s", err.Title, err.Detail)
}

type StandardErrors struct {
	Errors []*Error `json:"errors"`
}

// Error implements the error interface for the StandardErrors type.
// It joins the error messages from each Error instance with a comma separator.
func (se StandardErrors) Error() string {
	errorMessages := make([]string, len(se.Errors))
	for i, err := range se.Errors {
		errorMessages[i] = err.Error()
	}
	return fmt.Sprintf("[%s]", strings.Join(errorMessages, ", "))
}

type LinksAble struct {
	Meta   map[string]any `json:"meta,omitempty"`
	Links  map[string]any `json:"links,omitempty"`
	Status int            `json:"status"`
	Code   string         `json:"code"`
	Data   any            `json:"data"`
}

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at,omitempty"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiredAt   int64  `json:"expired_at"`
}
