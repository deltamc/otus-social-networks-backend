package controllers

import (
	"github.com/deltamc/otus-social-networks-backend/db"
	"github.com/deltamc/otus-social-networks-backend/logger"
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

	filter := users.Filter{
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
	}

	userList, err := users.GetUsers(filter)
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
var successCount int
var totalCount int

func HandleGenUser(w http.ResponseWriter, r *http.Request) {
	totalCount++
	user := users.GetUserRnd()
	_, err := user.New()
	if err == nil {
		successCount++
	}

	logger.Log.Printf(`Записано успешно: %d  Всего: %d`, successCount, totalCount)

	if err != nil {
		responses.Response500(w, err)
		return
	}
	responses.ResponseJson(w, user)
}

func HandleProfile(w http.ResponseWriter, r *http.Request, user users.User) {
	v := govalidator.New(requests.Profile(r))
	e := v.Validate()
	if len(e) > 0 {
		responses.Response422(w, e)
		return
	}

	age, _ := strconv.ParseInt(r.FormValue("age"), 10, 64)
	sex, _ := strconv.ParseInt(r.FormValue("sex"), 10, 64)

	user.FirstName =  r.FormValue("first_name")
	user.LastName =  r.FormValue("last_name")
	user.Interests =  r.FormValue("interests")
	user.City =  r.FormValue("city")
	user.Age =  age
	user.Sex =  sex

	err := user.Save()
	if err != nil {
		responses.Response500(w, err)
		return
	}

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

	if nerr, ok := err.(*mysql.MySQLError); ok && nerr.Number == db.ErrorDuplicateEntry {
		res := map[string][]string{
			"user_id": []string{"You are already friends"},
		}
		responses.Response422(w, res)
		return
	} else if err.Error() == users.ErrorFriendsWithYourself {
		res := map[string][]string{
			"user_id": []string{users.ErrorFriendsWithYourself},
		}
		responses.Response422(w, res)
		return
	}
	responses.Response500(w, err)
	return

}




