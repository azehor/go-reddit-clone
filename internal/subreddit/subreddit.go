package subreddit

import (
	"github.com/azehor/go-reddit-clone/internal/subreddit/model"
)

type Store interface {
	InsertSubreddit(subreddit *model.Subreddit) (*model.Subreddit, error)
	GetSubreddit(id string) (*model.Subreddit, error)
	GetSubredditList(ordering string) ([]*model.Subreddit, error)
	//TODO: finish CRUD
}

type Subreddits struct {
	store Store
}

func New(s Store) *Subreddits {
	return &Subreddits{
		store: s,
	}
}

func (s *Subreddits) InsertSubreddit(sb *model.Subreddit) (*model.Subreddit, error) {
	createdSubreddit, err := s.store.InsertSubreddit(sb)
	if err != nil {
		return nil, err
	}
	return createdSubreddit, nil
}

func (s *Subreddits) GetSubreddit(id string) (*model.Subreddit, error) {
	subreddit, err := s.store.GetSubreddit(id)
	if err != nil {
		return nil, err
	}

	return subreddit, nil
}

func (s *Subreddits) GetSubredditList(ordering string) ([]*model.Subreddit, error) {
	subreddit_list, err := s.store.GetSubredditList(ordering)
	if err != nil {
		return nil, err
	}
	return subreddit_list, nil
}
