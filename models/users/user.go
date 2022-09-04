package users

import (
	"errors"
	"github.com/deltamc/otus-social-networks-backend/db"
	"github.com/deltamc/otus-social-networks-backend/tar"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
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

type UserTarantool struct {
	Id        int64  `db:"id" json:"id"`
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

func GetUsersTarantool(filter Filter) (users []UserTarantool, err error) {
	err = tar.Client("0").CallTyped("get_users", []interface{}{filter.FirstName, filter.LastName}, &users)

	if len(users) == 1 && users[0].Id == 0 {
		users = []UserTarantool{}
		return
	}
	return
}

func GetUserRnd() (user User) {
	ids = make(map[int]int)
	id := getId()
	logPas := "user-" + strconv.Itoa(id)
	if pass == "" {
		pass, _ = hashedPass("123456")
	}

	user = User{
		LastName:  getLastName(),
		FirstName: getFirstName(),
		Sex:       getSex(),
		Age:       int64(rnd(18, 99)),
		Login:     logPas,
		Password:  pass,
		City:      "London",
		Interests: "",
	}
	//user.hashedPass()

	return
}

func getId() int {
	id := rnd(100, 99999999999)
	if _, ok := ids[id]; ok {
		id = getId()
	}
	ids[id] = 1
	return id
}

func getLastName() (name string) {
	if len(lastNames) == 0 {
		file := os.Getenv("LAST_NAME_FILE")
		var err error
		lastNames, err = readFile(file)
		if err != nil {
			panic(err)
		}
	}
	r := rnd(0, len(lastNames)-1)
	name = lastNames[r]
	name = strings.Title(strings.ToLower(name))
	return
}

func getFirstName() (name string) {
	if len(firstNames) == 0 {
		file := os.Getenv("FIRST_NAME_FILE")
		var err error
		firstNames, err = readFile(file)
		if err != nil {
			panic(err)
		}
	}
	r := rnd(0, len(firstNames)-1)
	name = firstNames[r]
	name = strings.Title(strings.ToLower(name))
	return
}

func getSex() (sex int64) {
	sex = int64(rnd(1, 2))
	return
}

func rnd(min int, max int) (r int) {
	rand.Seed(time.Now().UnixNano())
	r = min + rand.Intn(max-min+1)
	return
}

func readFile(fname string) (list []string, err error) {
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(b), "\n")

	list = make([]string, 0, len(lines))

	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		list = append(list, strings.Title(l))
	}

	return list, nil
}