package requests

import (
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

func Post(r *http.Request) govalidator.Options {
	rules := govalidator.MapData{
		"body": []string{"required"},

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