package posts

import (
	"github.com/azehor/go-reddit-clone/internal/posts/model"
)

type Store interface {
	InsertPost(post *model.Post) (*model.Post, error)
	GetPost(id string) (*model.Post, error)
	GetPosts(subreddit string) ([]*model.Post, error)
	//TODO: Add the rest of CRUD operations
}

type Posts struct {
	store Store
}

func New(s Store) *Posts {
	return &Posts{
		store: s,
	}
}

func (p *Posts) InsertPost(post *model.Post) (*model.Post, error) {
	createdPost, err := p.store.InsertPost(post)
	if err != nil {
		return nil, err
	}

	return createdPost, nil
}

func (p *Posts) GetPost(id string) (*model.Post, error) {
	post, err := p.store.GetPost(id)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (p *Posts) GetPosts(subreddit string) ([]*model.Post, error) {
	posts, err := p.store.GetPosts(subreddit)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
