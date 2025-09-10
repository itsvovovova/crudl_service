package service

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"log"
)

func ReadUserData(w http.ResponseWriter, r *http.Request, RequestStruct interface{}) {
	log.Println("Reading user data from request body")
	rv := reflect.ValueOf(RequestStruct)
	if rv.Kind() != reflect.Ptr {
		log.Println("RequestStruct is not a pointer")
		http.Error(w, "RequestStruct must be a pointer", http.StatusInternalServerError)
		return
	}

	log.Println("Reading request body data")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read request body")
		http.Error(w, "Failed to read request data", http.StatusBadRequest)
		return
	}

	log.Println("Unmarshaling JSON data")
	if err := json.Unmarshal(data, RequestStruct); err != nil {
		log.Println("Failed to unmarshal JSON data")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	log.Println("User data read successfully")
}
