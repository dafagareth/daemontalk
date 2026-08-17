package templates

import (
	"time"

	"daemontalk/internal/comment"
)

// CommentGroup groups consecutive comments posted by the same user.
type CommentGroup struct {
	Name      string
	CreatedAt time.Time
	Comments  []comment.Comment
}

func groupComments(comments []comment.Comment) []CommentGroup {
	var groups []CommentGroup
	for _, c := range comments {
		if len(groups) > 0 && groups[len(groups)-1].Name == c.Name {
			last := &groups[len(groups)-1]
			last.Comments = append(last.Comments, c)
		} else {
			groups = append(groups, CommentGroup{
				Name:      c.Name,
				CreatedAt: c.CreatedAt,
				Comments:  []comment.Comment{c},
			})
		}
	}
	return groups
}
