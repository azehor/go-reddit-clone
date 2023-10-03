package http

import (
	"html/template"

	"github.com/go-chi/chi/v5"

	p "github.com/azehor/go-reddit-clone/internal/posts/model"
	s "github.com/azehor/go-reddit-clone/internal/subreddit/model"
)

type Posts interface {
	InsertPost(post *p.Post) (*p.Post, error)
	GetPost(id string) (*p.Post, error)
	GetPosts(subreddit string) ([]*p.Post, error)
}

type Subreddits interface {
	InsertSubreddit(subreddit *s.Subreddit) (*s.Subreddit, error)
	GetSubreddit(id string) (*s.Subreddit, error)
	GetSubredditList(ordering string) ([]*s.Subreddit, error)
}

type Server struct {
	templates  *template.Template
	posts      Posts
	subreddits Subreddits
}

func New(p Posts, s Subreddits) *Server {
	tmpl := template.Must(template.ParseGlob("./web/templates/*.html"))
	return &Server{
		posts:      p,
		subreddits: s,
		templates:  tmpl,
	}
}

func (s *Server) AddRoutes(r *chi.Mux) {
	r.Get("/", s.getFrontPage)

	r.Post("/search/", s.getEmpty)
	r.Get("/submit", s.getSubmitPage)
	r.Post("/submit", s.createPost)

	r.Route("/r", func(r chi.Router) {
		r.Route("/{subreddit}", func(r chi.Router) {
			r.Get("/", s.getFrontPage)
			r.Get("/{ordering}", s.getFrontPage)
			r.Get("/comments/{id}", s.getPostPage)
			r.Get("/submit", s.getSubmitPage)
			r.Get("/sidebar", s.getSidebar)
		})
	})

	r.Route("/feeds", func(r chi.Router) {
		r.Get("/community-feed", s.getSubredditFeed)
		r.Get("/all-feed", s.getSubredditFeed)
		r.Get("/home-feed", s.getHomeFeed)
		r.Get("/popular-feed", s.getSubredditFeed)
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/{username}/joined", s.getEmpty)
		r.Get("/{username}/submit", s.getSubmitPage)
	})

	r.Route("/subreddits", func(r chi.Router) {
		r.Get("/", s.getSubredditList)
		r.Get("/{ordering}", s.getSubredditList)
		r.Get("/create", s.getCreateSubredditPage)
		r.Post("/create", s.createSubreddit)
	})
}
