package service

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
)

func ReadUserData(w http.ResponseWriter, r *http.Request, RequestStruct interface{}) {
	rv := reflect.ValueOf(RequestStruct)
	if rv.Kind() != reflect.Ptr {
		http.Error(w, "RequestStruct must be a pointer", http.StatusInternalServerError)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request data", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(data, RequestStruct); err != nil {
		http.Error(w, "Unmarshall data failed", http.StatusBadRequest)
		return
	}
}
