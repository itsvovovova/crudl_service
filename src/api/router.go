package api

import (
	"crudl_service/src/db"
	"crudl_service/src/service"
	"crudl_service/src/types"
	"encoding/json"
	"net/http"
)

func CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var request types.UserSubscription
	service.ReadUserData(w, r, &request)
	if err := db.CreateUserSubscription(&request); err != nil {
		http.Error(w, "Couldn't link subscription and database", http.StatusInternalServerError)
	}
}

func ReadSubscription(w http.ResponseWriter, r *http.Request) {
	var request types.UserSubscriptionData
	service.ReadUserData(w, r, &request)
	responseBody, err := db.GetUserSubscription(&request)
	if err != nil {
		http.Error(w, "Couldn't find subscription", http.StatusBadRequest)
	}
	bodyMarshal, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, "Failed to marshal user's subscription data", http.StatusInternalServerError)
	}
	w.Write(bodyMarshal)
}

func UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	var request types.UserSubscription
	service.ReadUserData(w, r, &request)
	if err := db.UpdateUserSubscription(&request); err != nil {
		http.Error(w, "Couldn't link subscription and database", http.StatusInternalServerError)
	}
}

func DeleteSubsription(w http.ResponseWriter, r *http.Request) {
	var request types.UserSubscriptionData
	service.ReadUserData(w, r, &request)
	if err := db.DeleteUserSubscription(&request); err != nil {
		http.Error(w, "Couldn't link subscription and database", http.StatusInternalServerError)
	}
}

func ListSubscription(w http.ResponseWriter, r *http.Request) {
	var request types.ListUserRequest
	service.ReadUserData(w, r, &request)
	responseSlice, err := db.ListUserSubscriptions(&request)
	if err != nil {
		http.Error(w, "Couldn't link user and database", http.StatusInternalServerError)
	}
	bodyMarshal, err := json.Marshal(responseSlice)
	if err != nil {
		http.Error(w, "Failed to marshal user's subscription data", http.StatusInternalServerError)
	}
	w.Write(bodyMarshal)
}
