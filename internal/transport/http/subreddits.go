package http

import (
	"log"
	"net/http"

	"github.com/azehor/go-reddit-clone/internal/subreddit/model"
	"github.com/go-chi/chi/v5"
)

func (s *Server) getSidebar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "subreddit")
	data, err := s.subreddits.GetSubreddit(id)
	if err != nil {
		log.Println(err)
	}
	err = s.templates.ExecuteTemplate(w, "sidebar.html", data)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) getSubredditList(w http.ResponseWriter, r *http.Request) {
	var ordering string
	if ordering := chi.URLParam(r, "ordering"); ordering == "" {
		ordering = "hot" //TODO: if signed in, default value should be pulled from user settings
	}
	data, err := s.subreddits.GetSubredditList(ordering)
	if err != nil {
		log.Fatal(err)
	}
	err = s.templates.ExecuteTemplate(w, "subreddit_list.html", data)
	if err != nil {
		log.Println(err)
	}

}

func (s *Server) createSubreddit(w http.ResponseWriter, r *http.Request) {
	var name = r.FormValue("subreddit_name")
	var title = r.FormValue("subreddit_title")
	var description = r.FormValue("subreddit_description")
	var rules = r.FormValue("subreddit_rules")
	var sub_type = r.FormValue("subreddit_type")
	var over18_string = r.Form["subreddit_over18"]

	//TODO: Validate inputs

	sub_url := "r/" + name
	over18 := len(over18_string) > 0
	log.Print(over18)
	sub := model.Subreddit{
		Description:   &description,
		Rules:         &rules,
		SubredditType: &sub_type,
		Title:         &title,
		URL:           &sub_url,
		Over18:        &over18,
	}

	data, err := s.subreddits.InsertSubreddit(&sub)
	if err != nil {
		w.WriteHeader(500)
	} else {
		log.Print(data)
		w.Header().Add("HX-Redirect", "//"+r.Host+"/"+*data.URL)
		w.WriteHeader(200)
	}
}
