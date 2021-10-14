package requests

import (
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

func Profile(r *http.Request) govalidator.Options {
	rules := govalidator.MapData{
		"first_name": []string{"required"},
		"last_name": []string{"required"},
		"age": []string{"required","numeric","numeric_between:1,99"},
		"sex": []string{"required","numeric","numeric_between:1,2"},
		"interests": []string{"required"},
		"city": []string{"required"},
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