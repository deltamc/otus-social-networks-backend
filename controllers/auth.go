package controllers

import (
	"database/sql"
	"github.com/deltamc/otus-social-networks-backend/auth"
	"github.com/deltamc/otus-social-networks-backend/db"
	"github.com/deltamc/otus-social-networks-backend/models/users"
	"github.com/deltamc/otus-social-networks-backend/requests"
	"github.com/deltamc/otus-social-networks-backend/responses"
	"github.com/go-sql-driver/mysql"
	"github.com/thedevsaddam/govalidator"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
)

func HandleSignUp(w http.ResponseWriter, r *http.Request) {

	v := govalidator.New(requests.SignUp(r))
	e := v.Validate()
	if len(e) > 0 {
		responses.Response422(w, e)
		return
	}

	age, _ := strconv.ParseInt(r.FormValue("age"), 10, 64)
	sex, _ := strconv.ParseInt(r.FormValue("sex"), 10, 64)
	user := users.User{
		Password:  r.FormValue("password"),
		Login:     r.FormValue("login"),
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
		Age:       age,
		Sex:       sex,
		Interests: r.FormValue("interests"),
		City:      r.FormValue("city"),
	}
	user.HashedPass()

	_, err := user.Save()

	if nerr, ok := err.(*mysql.MySQLError); ok && nerr.Number == db.ERROR_DUPLICATE_ENTRY {
		res := map[string][]string{
			"login": []string{"The «Login» has already been taken."},
		}
		responses.Response422(w, res)
		return
	}

	if err != nil {
		responses.Response500(w, err)
		return
	}

	err = auth.GetBearerResponse(w, r, user)

	if err != nil {
		responses.Response500(w, err)
		return
	}
}


func HandleSignIn(w http.ResponseWriter, r *http.Request) {

	v := govalidator.New(requests.SignIn(r))
	e := v.Validate()
	if len(e) > 0 {
		responses.Response422(w, e)
		return
	}

	user, err := users.GetUserByLogin(r.FormValue("login"))

	if err !=nil {
		if err == sql.ErrNoRows {
			res := map[string][]string{
				"login": []string{"These credentials do not match our records."},
			}
			responses.Response422(w, res)
			return
		} else {
			responses.Response500(w, err)
			return
		}



	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.FormValue("password")))
	if err !=nil {
		log.Println(err)

		res := map[string][]string{
			"login": []string{"These credentials do not match our records."},
		}
		responses.Response422(w, res)
		return
	}

	err = auth.GetBearerResponse(w, r, user)

	if err != nil {
		responses.Response500(w, err)
		return
	}
}

func HandleRefresh(w http.ResponseWriter, r *http.Request, user users.User) {

	err := auth.GetBearerResponse(w, r, user)

	if err != nil {
		responses.Response500(w, err)
		return
	}
}