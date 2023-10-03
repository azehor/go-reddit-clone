package http

import (
	"html/template"
	"log"
	"net/http"

	p "github.com/azehor/go-reddit-clone/internal/posts/model"
	"github.com/go-chi/chi/v5"
)

func (s *Server) getFrontPage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		CurrentSubreddit, Ordering, Username string
		NavigationSessionID                  string
		IsFeed                               bool
	}{
		CurrentSubreddit:    chi.URLParam(r, "subreddit"),
		Ordering:            chi.URLParam(r, "ordering"),
		NavigationSessionID: chi.URLParam(r, "navigationSessionID"),
		Username:            "",
		IsFeed:              false,
	}

	if data.CurrentSubreddit == "" {
		data.IsFeed = true
		if data.Username == "" {
			data.CurrentSubreddit = "popular"
		} else {
			data.CurrentSubreddit = "home"
		}
	}

	//TODO: get username via sessionID

	//TODO: if user is signed in, default ordering should be determined by user settings
	if data.Ordering == "" {
		data.Ordering = "hot"
	}

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

func (s *Server) getCreateSubredditPage(w http.ResponseWriter, r *http.Request) {
	//TODO: should only be accesible if authenticated, else redirect to login page

	s.templates = template.Must(template.ParseGlob("./web/templates/*.html"))
	err := s.templates.ExecuteTemplate(w, "new_subreddit.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) getPostPage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		CurrentSubreddit, Username string
		Post                       *p.Post
	}{
		CurrentSubreddit: chi.URLParam(r, "subreddit"),
		Username:         "",
	}

	post, err := s.posts.GetPost(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
	}

	data.Post = post
	s.templates = template.Must(template.ParseGlob("./web/templates/*.html"))
	err = s.templates.ExecuteTemplate(w, "post_page.html", data)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) getEmpty(w http.ResponseWriter, r *http.Request) {
	log.Println("empty request made")
}
