package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"github.com/azehor/go-reddit-clone/internal/posts"
	pstore "github.com/azehor/go-reddit-clone/internal/posts/store"
	"github.com/azehor/go-reddit-clone/internal/subreddit"
	sstore "github.com/azehor/go-reddit-clone/internal/subreddit/store"
	transport "github.com/azehor/go-reddit-clone/internal/transport/http"
	"github.com/go-chi/chi/v5"
)

func Start() {
	db := initDatabase()
	r := chi.NewRouter()

	ps := pstore.New(db)
	p := posts.New(ps)
	ss := sstore.New(db)
	s := subreddit.New(ss)

	httpServer := transport.New(p, s)
	httpServer.AddRoutes(r)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		shutdown()
		os.Exit(1)
	}()

	http.ListenAndServe(":8080", r)
	shutdown()
}

func initDatabase() *sqlx.DB {
	db, err := sqlx.Connect("sqlite", "file:test.db?mode=rwc")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func shutdown() {
	fmt.Printf("Shutting down\n")
}
