package users

import (
	"errors"
	"github.com/deltamc/otus-social-networks-backend/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int64  `db:"id" json:"id"`
	Login     string `db:"login" json:"login"`
	Password  string `db:"password" json:"-"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Age       int64  `db:"age" json:"age"`
	Sex       int64  `db:"sex" json:"sex"`
	Interests string `db:"interests" json:"interests"`
	City      string `db:"city" json:"city"`
}

const ERROR_FRIENDS_WITH_YOURSELF string = "You can't make friends with yourself"

func (u *User) HashedPass() error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPass)
	return nil
}

func (u *User) New() (lastID int64, err error) {
	dbPool := db.OpenDB()

	stmt, err := dbPool.Prepare(
		"INSERT INTO " +
			"`users` (`login`, `password`, `first_name`, `last_name`, `age`, `sex`, `interests`, `city`) " +
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(u.Login, u.Password, u.FirstName, u.LastName, u.Age, u.Sex, u.Interests, u.City)
	if err != nil {
		return
	}

	lastID, err = res.LastInsertId()
	if err != nil {
		return
	}

	u.Id = lastID

	return
}

func (u *User) Save() (err error) {

	dbPool := db.OpenDB()

	stmt, err := dbPool.Prepare(`UPDATE users SET first_name = ?, last_name = ?, age =?, sex =?, interests=?, city=? WHERE id=?`)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.FirstName, u.LastName, u.Age, u.Sex, u.Interests, u.City, u.Id)
	if err != nil {
		return
	}

	return
}

func (u *User) GetFriends() (friends []User, err error) {

	dbPool := db.OpenDB()

	sqlStmt := `SELECT 
					users.id, login, first_name, last_name, age, sex, interests, city
				FROM 
					friends 
					LEFT JOIN users ON users.id = friends.friend_id
				WHERE 
					friends.user_id = ? 
				ORDER BY 
					users.id DESC`

	rows, err := dbPool.Query(sqlStmt, u.Id)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var user User

		err = rows.Scan(
			&user.Id,
			&user.Login,
			&user.FirstName,
			&user.LastName,
			&user.Age,
			&user.Sex,
			&user.Interests,
			&user.City)

		friends = append(friends, user)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (u *User) GetFlowers() (flowers []User, err error) {

	dbPool := db.OpenDB()

	sqlStmt := `SELECT 
					users.id, login, first_name, last_name, age, sex, interests, city
				FROM 
					friends 
					LEFT JOIN users ON users.id = friends.user_id
				WHERE 
					friends.friend_id = ? 
				ORDER BY 
					users.id DESC`

	rows, err := dbPool.Query(sqlStmt, u.Id)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var user User

		err = rows.Scan(
			&user.Id,
			&user.Login,
			&user.FirstName,
			&user.LastName,
			&user.Age,
			&user.Sex,
			&user.Interests,
			&user.City)

		flowers = append(flowers, user)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (u *User) MakeFriend(userId int64) (err error) {
	dbPool := db.OpenDB()

	if userId == u.Id {
		err = errors.New(ERROR_FRIENDS_WITH_YOURSELF)
		return
	}

	stmt, err := dbPool.Prepare(
		`INSERT INTO friends (user_id, friend_id) VALUES (?, ?)`)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Id, userId)
	if err != nil {
		return
	}

	return
}

func GetUserByLogin(login string) (user User, err error) {

	dbPool := db.OpenDB()

	sqlStmt := `SELECT 
					*
				FROM 
					users 
				WHERE 
					login = ?`

	// Prepare statement
	stmt, err := dbPool.Prepare(sqlStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(login).Scan(
		&user.Id,
		&user.Login,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Age,
		&user.Sex,
		&user.Interests,
		&user.City)

	return
}

func GetUserById(id int64) (user User, err error) {
	dbPool := db.OpenDB()

	sqlStmt := `SELECT 
					*
				FROM 
					users 
				WHERE 
					id = ?`

	// Prepare statement
	stmt, err := dbPool.Prepare(sqlStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(
		&user.Id,
		&user.Login,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Age,
		&user.Sex,
		&user.Interests,
		&user.City)

	return
}

func GetUsers(filter Filter) (users []User, err error) {

	dbPool := db.OpenDB()

	where, args := filter.getWhere()

	sqlStmt := `SELECT 
					id, login, first_name, last_name, age, sex, interests, city
				FROM 
					users 
					` + where + `
				ORDER BY 
					id DESC`

	//если нет фильтров, ограничиваем вывод
	if where == "" {
		sqlStmt += " LIMIT 100"
	}

	rows, err := dbPool.Query(sqlStmt, args...)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var user User

		if err = rows.Err(); err != nil {
			return
		}

		err = rows.Scan(
			&user.Id,
			&user.Login,
			&user.FirstName,
			&user.LastName,
			&user.Age,
			&user.Sex,
			&user.Interests,
			&user.City)

		users = append(users, user)
	}
	return
}
