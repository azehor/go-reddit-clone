package store

import (
	"database/sql"
	"log"
	"time"

	"github.com/azehor/go-reddit-clone/internal/subreddit/model"
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

// TODO: Extract initTable into SQL file
func initTable(db DB) {
	var postSchema = `
    CREATE TABLE IF NOT EXISTS subreddit (
        id integer NOT NULL PRIMARY KEY,
        description text,
        rules text,
        over18 boolean,
        subscribers integer DEFAULT 0,
        subreddit_type text,
        title text,
        url text,
        created_at datetime
    );`

	res := db.MustExec(postSchema)
	log.Print(res)
}

var timeNow = func() *time.Time {
	now := time.Now().UTC()
	return &now
}

func (s *Store) InsertSubreddit(sb *model.Subreddit) (*model.Subreddit, error) {
	sb.CreatedAt = timeNow()
	res, err := s.db.NamedQuery(`INSERT INTO
        subreddit(description, rules, over18, subreddit_type, title, url, created_at)
        VALUES(:description, :rules, :over18, :subreddit_type, :title, :url, :created_at)
        RETURNING *`, sb)
	if err != nil {
		log.Panic(err)
	}
	defer res.Close()

	if !res.Next() {
		log.Print(err)
	}
	subreddit := &model.Subreddit{}

	if err = res.StructScan(subreddit); err != nil {
		return nil, err
	}

	return subreddit, nil
}

func (s *Store) GetSubreddit(id string) (*model.Subreddit, error) {
	var subreddit model.Subreddit
	if err := s.db.Get(&subreddit, "SELECT * FROM subreddit WHERE id = $1", id); err != nil {
		return nil, err
	}
	return &subreddit, nil
}

func (s *Store) GetSubredditList(ordering string) ([]*model.Subreddit, error) {
	var list []*model.Subreddit
	stmt := "SELECT * FROM subreddit LIMIT 25"
	if err := s.db.Select(&list, stmt); err != nil {
		return nil, err
	}
	return list, nil
}
