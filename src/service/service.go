package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

func ReadUserData(w http.ResponseWriter, r *http.Request, requestStruct interface{}) {
	log.Info("Reading user data from request body")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(requestStruct); err != nil {
		log.Error("Failed to unmarshal JSON data")
		http.Error(w, "Incorrect input data format", http.StatusBadRequest)
		return
	}
	log.Info("User data read successfully")
}

func GetIDRequest(r *http.Request) (int64, error) {
	log.Info("Reading id from URL path parameter")
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Failed to parse ID")
		return 0, err
	}
	return id, nil
}
