package controllers

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/deltamc/otus-social-networks-backend/models/users"
	"github.com/deltamc/otus-social-networks-backend/responses"
	"io"
	"net/http"
	"os"
)

var mc  *memcache.Client

func HandleFeed(w http.ResponseWriter, r *http.Request, user users.User) {
	w.Header().Add("Content-Type", "application/json")
	if mc == nil {
		mc = memcache.New(fmt.Sprintf("%s:%s", os.Getenv("MEMCACHE_HOST"), os.Getenv("MEMCACHE_PORT")))
	}

	feed, err := mc.Get(fmt.Sprintf("feed_%d", user.Id))
	list := "[]"
	if err != nil {
		if err.Error() != "memcache: cache miss" {
			responses.Response500(w, err)
			return
		}
	} else {
		list = string(feed.Value)
	}


	_, err = io.WriteString(w, list)

	if err != nil {
		responses.Response500(w, err)
		return
	}
}





