package components

import "time"

type Post struct {
	Title      string
	Summary    string
	Tags       []string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Authors    []string
	LayoutName string
}
