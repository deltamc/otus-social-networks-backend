package test

import (
	"github.com/deltamc/otus-social-networks-backend/db"
)

type Test struct {
	Id int64 `db:"id" json:"id"`
	Rnd1 int64 `db:"rnd1" json:"rnd1"`
	Rnd2 int64 `db:"rnd2" json:"rnd2"`
	Rnd3 int64 `db:"rnd3" json:"rnd3"`
}

func (t *Test) New() (lastID int64, err error)  {
	dbPool := db.OpenDB(db.MasterName)
	
	stmt, err := dbPool.Prepare(
		"INSERT INTO " +
			"`test` (`rnd1`, `rnd2`, `rnd3`) " +
			"VALUES (?, ?, ?)")
	if err != nil {
		return
	}
	defer stmt.Close()


	res, err := stmt.Exec(t.Rnd1, t.Rnd2, t.Rnd3)
	if err != nil {
		return
	}

	lastID, err = res.LastInsertId()
	if err != nil {
		return
	}

	t.Id = lastID

	return
}

func (t *Test) Save() (err error)  {
	dbPool := db.OpenDB(db.MasterName)

	stmt, err := dbPool.Prepare(`UPDATE test SET rnd1 = ?, rnd2 = ?, rnd3 =? WHERE id=?`)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(t.Rnd1, t.Rnd2, t.Rnd3, t.Id)
	if err != nil {
		return
	}

	return
}


