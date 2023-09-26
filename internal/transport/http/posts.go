package http

import (
	"log"
	"net/http"

	"github.com/azehor/go-reddit-clone/internal/posts/model"
	"github.com/go-chi/chi/v5"
)

func (s *Server) getPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	data, err := s.posts.GetPost(id)
	if err != nil {
		log.Println(err)
	}
	err = s.templates.ExecuteTemplate(w, "post.html", data)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) getPosts(w http.ResponseWriter, r *http.Request) {
	var subreddit = chi.URLParam(r, "subreddit")
	log.Print(subreddit)
	data, err := s.posts.GetPosts(subreddit)
	if err != nil {
		log.Print(err)
	}
	for _, p := range data {
		log.Print(p)
	}
	err = s.templates.ExecuteTemplate(w, "post_list.html", data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) createPost(w http.ResponseWriter, r *http.Request) {
	var title = r.FormValue("newPostTitle")
	log.Printf("Posted Title: %v", title)
	var body = r.FormValue("newPostBody")
	var subreddit = r.FormValue("newPostSubreddit")

	//TODO: Validate inputs

	p := model.Post{
		Title:     &title,
		Body:      &body,
		Subreddit: &subreddit,
	}

	data, err := s.posts.InsertPost(&p)
	if err != nil {
		log.Print(err)
		w.WriteHeader(504)
	} else {
		log.Print(data)
		w.Header().Add("HX-Redirect", "/")
		w.WriteHeader(200)
	}
}
