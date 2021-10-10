package controllers

import (
	"github.com/deltamc/otus-social-networks-backend/db"
	"github.com/deltamc/otus-social-networks-backend/models/users"
	"github.com/deltamc/otus-social-networks-backend/requests"
	"github.com/deltamc/otus-social-networks-backend/responses"
	"github.com/go-sql-driver/mysql"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"strconv"
)

func HandleMy(w http.ResponseWriter, r *http.Request, user users.User) {
	responses.ResponseJson(w, user)
}

func HandleUsers(w http.ResponseWriter, r *http.Request) {
	userList, err := users.GetUsers()
	if err != nil {
		responses.Response500(w, err)
	}
	responses.ResponseJson(w, userList)
}


func HandleFriends(w http.ResponseWriter, r *http.Request, user users.User) {
	userList, err := user.GetFriends()
	if err != nil {
		responses.Response500(w, err)
		return
	}
	responses.ResponseJson(w, userList)
}

func HandleMakeFriend(w http.ResponseWriter, r *http.Request, user users.User) {

	v := govalidator.New(requests.MakeFriend(r))
	e := v.Validate()
	if len(e) > 0 {
		responses.Response422(w, e)
		return
	}

	userId, err := strconv.ParseInt(r.FormValue("user_id"), 10, 32)

	err = user.MakeFriend(userId)

	if err == nil {
		responses.Response200(w)
		return
	}

	if nerr, ok := err.(*mysql.MySQLError); ok && nerr.Number == db.ERROR_DUPLICATE_ENTRY {
		res := map[string][]string{
			"user_id": []string{"You are already friends"},
		}
		responses.Response422(w, res)
		return
	} else if err.Error() == users.ERROR_FRIENDS_WITH_YOURSELF {
		res := map[string][]string{
			"user_id": []string{users.ERROR_FRIENDS_WITH_YOURSELF},
		}
		responses.Response422(w, res)
		return
	}
	responses.Response500(w, err)
	return

}


