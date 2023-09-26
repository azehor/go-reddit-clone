package http

import (
	"html/template"
	"log"
	"net/http"

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
		r.Get("/all", s.getPosts)
		r.Route("/{subreddit}", func(r chi.Router) {
			r.Get("/", s.getPosts)
			r.Get("/{ordering}", s.getPosts)
			r.Get("/comments/{id}", s.getPost)
			r.Get("/submit", s.getSubmitPage)
			r.Get("/sidebar", s.getSidebar)
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/{username}/joined", s.getEmpty)
		r.Get("/{username}/submit", s.getSubmitPage)
	})

	r.Get("/subreddits/{ordering}", s.getEmpty) //TODO: gets list of subreddits ordered
	r.Get("/subreddits/create", s.getNewSubredditPage)
	r.Post("/subreddits/create", s.createSubreddit)

}

func (s *Server) getFrontPage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		CurrentSubreddit, Ordering, Username string
		LoggedIn                             string
	}{
		CurrentSubreddit: "",
		Ordering:         "hot",
		Username:         "",
		LoggedIn:         "",
	}
	//TODO: should get a listing with posts from either r/all or logged user's 'joined' subs
	//TODO: remove reloading template when done with css
	s.templates = template.Must(template.ParseGlob("./web/templates/*.html"))
	err := s.templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) getSubmitPage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Subreddit, Username string
	}{
		Subreddit: chi.URLParam(r, "subreddit"),
		Username:  chi.URLParam(r, "username"),
	}

	s.templates = template.Must(template.ParseGlob("./web/templates/*.html"))
	err := s.templates.ExecuteTemplate(w, "submit.html", data)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) getNewSubredditPage(w http.ResponseWriter, r *http.Request) {
	//TODO: should only be accesible if authenticated, else redirect to login page

	s.templates = template.Must(template.ParseGlob("./web/templates/*.html"))
	err := s.templates.ExecuteTemplate(w, "newsubreddit.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) getEmpty(w http.ResponseWriter, r *http.Request) {
	log.Println("empty request made")
}
