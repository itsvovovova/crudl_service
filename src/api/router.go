package api

import (
	"crudl_service/src/db"
	"crudl_service/src/service"
	"crudl_service/src/types"
	"encoding/json"
	"net/http"
	"time"
)

func CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var request types.UserSubscription
	service.ReadUserData(w, r, &request)
	if _, err := time.Parse("01-2006", *request.StartDate); err != nil {
		http.Error(w, "Incorrect time format", http.StatusBadRequest)
	}
	if request.EndDate != nil {
		_, err := time.Parse("01-2006", *request.EndDate)
		if err != nil {
			http.Error(w, "Incorrect time format", http.StatusBadRequest)
			return
		}
	}
	if err := db.CreateUserSubscription(&request); err != nil {
		http.Error(w, "Couldn't link subscription and database", http.StatusInternalServerError)
		return
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
	var request types.UserRequest
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

func SumUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	var request types.UserSumSubscriptionRequest
	service.ReadUserData(w, r, &request)
	responseSum, err := db.GetSumUserSubscription(&request)
	if err != nil {
		http.Error(w, "Couldn't link user and database", http.StatusInternalServerError)
	}
	var response = types.UserSubscriptionSumResponse{
		UserId:     request.UserId,
		CurrentSum: responseSum}
	bodyMarshal, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal user's subscription data", http.StatusInternalServerError)
	}
	w.Write(bodyMarshal)
}
