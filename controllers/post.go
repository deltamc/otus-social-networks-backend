package controllers

import (
	"github.com/deltamc/otus-social-networks-backend/models/post"
	"github.com/deltamc/otus-social-networks-backend/models/users"
	"github.com/deltamc/otus-social-networks-backend/queue"
	"github.com/deltamc/otus-social-networks-backend/requests"
	"github.com/deltamc/otus-social-networks-backend/responses"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"os"
)



func HandlePost(w http.ResponseWriter, r *http.Request, user users.User) {
	v := govalidator.New(requests.Post(r))
	e := v.Validate()
	if len(e) > 0 {
		responses.Response422(w, e)
		return
	}

	body := r.FormValue("body")

	p := post.Post{
		Body:body,
		UserId:user.Id,
	}

	_,err := p.New()
	if err != nil{
		responses.Response500(w, err)
		return
	}
	rmq := queue.NewRabbitMQ()
	err = rmq.Connect()
	if err != nil{
		responses.Response500(w, err)
	}
	defer rmq.Close()


	err = rmq.Push(os.Getenv("RABBITMQ_QUEQE_FEED"), p)

	if err != nil{
		responses.Response500(w, err)
		return
	}

}





