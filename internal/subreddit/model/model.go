package model

import (
	"fmt"
	"time"
)

type Subreddit struct {
	ID            *string    `db:"id"`
	Description   *string    `db:"description"`
	Rules         *string    `db:"rules"`
	Over18        *bool      `db:"over18"`
	Subscribers   *int       `db:"subscribers"`
	SubredditType *string    `db:"subreddit_type"`
	Title         *string    `db:"title"`
	URL           *string    `db:"url"`
	CreatedAt     *time.Time `db:"created_at"`
}

func (s *Subreddit) String() string {
	return fmt.Sprintf("ID: %v\nDescription: %v\nRules: %v\nOver 18: %v\nSubscribers: %v\nType: %v\nTitle: %v\nURL: %v\nCreated At: %v\n",
		*s.ID, *s.Description, *s.Rules, *s.Over18, *s.Subscribers, *s.SubredditType, *s.Title, *s.URL, *s.CreatedAt)
}
