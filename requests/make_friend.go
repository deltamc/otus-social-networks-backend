package requests

import (
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

func MakeFriend(r *http.Request) govalidator.Options {
	rules := govalidator.MapData{
		"user_id": []string{"required"},

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