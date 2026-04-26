package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

func ReadUserData(w http.ResponseWriter, r *http.Request, requestStruct any) bool {
	if err := json.NewDecoder(r.Body).Decode(requestStruct); err != nil {
		log.WithError(err).Error("Failed to unmarshal JSON data")
		http.Error(w, "Incorrect input data format", http.StatusBadRequest)
		return false
	}
	return true
}

func GetIDRequest(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		log.WithError(err).Error("Failed to parse ID from URL")
	}
	return id, err
}
