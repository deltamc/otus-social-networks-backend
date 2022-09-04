package users

import (
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var lastNames []string
var firstNames []string
var ids map[int]int
var pass string

func hashedPass(pass string) (hashedPass string, err error) {
	hashedPassByte, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	hashedPass = string(hashedPassByte)
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
