package model

import (
	"fmt"
	"time"
)

type Post struct {
	ID             *string    `db:"id"`
	Title          *string    `db:"title"`
	Body           *string    `db:"body"`
	Summary        *string    `db:"summary"`
	Upvotes        *string    `db:"upvotes"`
	CommentAmmount *string    `db:"comment_ammount"`
	Subreddit      *string    `db:"subreddit"`
	TimeStamp      *string    `db:"timestamp"`
	CreatedAt      *time.Time `db:"created_at"`
}

func (p *Post) String() string {
	return fmt.Sprintf("ID: %v\nTitle: %v\nBody: %v\nSummary: %v\nUpvotes: %v\nComment Ammount: %v\nSubreddit: %v\nTimestamp: %v\nCreated At: %v\n",
		*p.ID, *p.Title, *p.Body, *p.Summary, *p.Upvotes, *p.CommentAmmount, *p.Subreddit, *p.TimeStamp, *p.CreatedAt)
}
