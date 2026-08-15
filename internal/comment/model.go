package comment

import "time"

type Comment struct {
	ID        int64
	PostSlug  string
	Name      string
	Body      string
	CreatedAt time.Time
}
