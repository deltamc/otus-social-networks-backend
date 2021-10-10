package requests

import (
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

func SignIn(r *http.Request) govalidator.Options {
	rules := govalidator.MapData{
		"login": []string{"required"},
		"password": []string{"required"},
	}

	messages := govalidator.MapData{

	}

	return govalidator.Options{
		Request:         r,
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: true,
	}
}