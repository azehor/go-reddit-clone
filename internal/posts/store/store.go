package store

import (
	"database/sql"
	"log"
	"time"

	"github.com/azehor/go-reddit-clone/internal/posts/model"
	"github.com/jmoiron/sqlx"
)

type DB interface {
	MustExec(query string, args ...interface{}) sql.Result
	NamedQuery(query string, args interface{}) (*sqlx.Rows, error)
	NamedExec(query string, args interface{}) (sql.Result, error)
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
}

type Store struct {
	db DB
}

func New(db DB) *Store {
	initTable(db)
	return &Store{
		db: db,
	}
}

var timeNow = func() *time.Time {
	now := time.Now().UTC()
	return &now
}

var timestamp = func(t *time.Time) *string {
	timestamp := time.Since(*t).Truncate(time.Hour).String()
	return &timestamp
}

// TODO: Extract initTable to a SQL file
func initTable(db DB) {
	var postSchema = `
    CREATE TABLE IF NOT EXISTS posts (
        id integer NOT NULL PRIMARY KEY,
        title text,
        body text,
        summary text,
        upvotes integer DEFAULT 0,
        comment_ammount integer DEFAULT 0,
        subreddit text,
        created_at datetime,
        timestamp datetime
    );`

	res := db.MustExec(postSchema)
	log.Print(res)
}

var summarize = func(s string) string {
	if len(s) > 50 {
		return string([]rune(s)[0:50])
	} else {
		return s
	}
}

func (s *Store) InsertPost(post *model.Post) (*model.Post, error) {
	var summary = summarize(*post.Body)

	post.CreatedAt = timeNow()
	post.Summary = &summary
	post.TimeStamp = timestamp(timeNow())

	res, err := s.db.NamedQuery(`INSERT INTO
        posts(title, body, summary, subreddit, created_at, timestamp)
        VALUES (:title, :body, :summary, :subreddit, :created_at, :timestamp)
        RETURNING *`, post)
	if err != nil {
		log.Print(err)
	}
	defer res.Close()

	if !res.Next() {
		return nil, nil //TODO: Handle this error
	}

	createdPost := &model.Post{}

	if err := res.StructScan(&createdPost); err != nil {
		return nil, err
	}

	return createdPost, nil
}

func (s *Store) GetPost(id string) (*model.Post, error) {
	var post model.Post
	if err := s.db.Get(&post, "SELECT * FROM posts WHERE id = $1", id); err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *Store) GetPosts(subreddit string) ([]*model.Post, error) {
	var posts []*model.Post
	var stmt string
	if subreddit != "" {
		stmt = "SELECT * FROM posts WHERE subreddit = '" + subreddit + "' LIMIT 25"
		log.Print(stmt)
	} else {
		stmt = "SELECT * FROM posts LIMIT 25"
	}
	if err := s.db.Select(&posts, stmt); err != nil {
		return nil, err
	}
	return posts, nil
}
