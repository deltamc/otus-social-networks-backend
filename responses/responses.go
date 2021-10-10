package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

type out422 struct {
	Errors interface{}  `json:"errors"`
}

func ResponseJson(w http.ResponseWriter, data interface{})  {
	w.Header().Add("Content-Type", "application/json")

	json.NewEncoder(w).Encode(data)
}

func Response422(w http.ResponseWriter, data interface{})  {
	w.Header().Add("Content-Type", "application/json")
	//http. Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
	w.WriteHeader(http.StatusUnprocessableEntity)

	json.NewEncoder(w).Encode( out422{data})
}

func Response500(w http.ResponseWriter, err error)  {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
}

func Response200(w http.ResponseWriter)  {
	w.WriteHeader(http.StatusOK)
}