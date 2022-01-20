package post

import (
	"github.com/deltamc/otus-social-networks-backend/db"
	"time"
)

type Post struct {
	Id int64 `db:"id" json:"id"`
	UserId int64 `db:"user_id" json:"user_id"`
	Body string `db:"body" json:"body"`
	CreatedAt time.Time `json:"created_at"  db:"created_at"`
	UserName string `json:"user_name" db:"user_name"`
}

func (p *Post) New() (lastID int64, err error) {
	dbPool := db.OpenDB()


	stmt, err := dbPool.Prepare(
		"INSERT INTO " +
			"`posts` (`user_id`, `body`, `created_at`) " +
			"VALUES (?, ?, NOW())")
	if err != nil {
		return
	}
	defer stmt.Close()


	res, err := stmt.Exec(p.UserId, p.Body)
	if err != nil {
		return
	}

	lastID, err = res.LastInsertId()
	if err != nil {
		return
	}

	p.Id = lastID


	return
}