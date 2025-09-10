package api

import (
	"crudl_service/src/db"
	"crudl_service/src/service"
	"crudl_service/src/types"
	"encoding/json"
	"net/http"
	"time"

	"log"
)

func CreateSubscription(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating new subscription request received")
	var request types.UserSubscription
	service.ReadUserData(w, r, &request)

	log.Println("Validating start date format")
	if _, err := time.Parse("01-2006", *request.StartDate); err != nil {
		log.Println("Invalid start date format provided")
		http.Error(w, "Incorrect time format", http.StatusBadRequest)
		return
	}

	if request.EndDate != nil {
		log.Println("Validating end date format")
		_, err := time.Parse("01-2006", *request.EndDate)
		if err != nil {
			log.Println("Invalid end date format provided")
			http.Error(w, "Incorrect time format", http.StatusBadRequest)
			return
		}
	}

	if err := db.CreateUserSubscription(&request); err != nil {
		log.Println("Failed to create subscription in database")
		http.Error(w, "Couldn't link subscription and database", http.StatusInternalServerError)
		return
	}

	log.Println("Subscription created successfully via API")
	w.WriteHeader(http.StatusCreated)
}

func ReadSubscription(w http.ResponseWriter, r *http.Request) {
	log.Println("Read subscription request received")
	var request types.UserSubscriptionData
	service.ReadUserData(w, r, &request)

	responseBody, err := db.GetUserSubscription(&request)
	if err != nil {
		log.Println("Failed to retrieve subscription from database")
		http.Error(w, "Couldn't find subscription", http.StatusBadRequest)
		return
	}

	log.Println("Marshaling subscription data to JSON")
	bodyMarshal, err := json.Marshal(responseBody)
	if err != nil {
		log.Println("Failed to marshal subscription data")
		http.Error(w, "Failed to marshal user's subscription data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bodyMarshal)
	log.Println("Subscription data sent successfully")
}

func UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	log.Println("Update subscription request received")
	var request types.UserSubscription
	service.ReadUserData(w, r, &request)

	if err := db.UpdateUserSubscription(&request); err != nil {
		log.Println("Failed to update subscription in database")
		http.Error(w, "Couldn't link subscription and database", http.StatusInternalServerError)
		return
	}

	log.Println("Subscription updated successfully via API")
	w.WriteHeader(http.StatusOK)
}

func DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete subscription request received")
	var request types.UserSubscriptionData
	service.ReadUserData(w, r, &request)

	if err := db.DeleteUserSubscription(&request); err != nil {
		log.Println("Failed to delete subscription from database")
		http.Error(w, "Couldn't link subscription and database", http.StatusInternalServerError)
		return
	}

	log.Println("Subscription deleted successfully via API")
	w.WriteHeader(http.StatusOK)
}

func ListSubscription(w http.ResponseWriter, r *http.Request) {
	log.Println("List subscriptions request received")
	var request types.UserRequest
	service.ReadUserData(w, r, &request)

	responseSlice, err := db.ListUserSubscriptions(&request)
	if err != nil {
		log.Println("Failed to retrieve subscriptions list from database")
		http.Error(w, "Couldn't link user and database", http.StatusInternalServerError)
		return
	}

	log.Println("Marshaling subscriptions list to JSON")
	bodyMarshal, err := json.Marshal(responseSlice)
	if err != nil {
		log.Println("Failed to marshal subscriptions list")
		http.Error(w, "Failed to marshal user's subscription data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bodyMarshal)
	log.Println("Subscriptions list sent successfully")
}

func SumUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	log.Println("Sum user subscriptions request received")
	var request types.UserSumSubscriptionRequest
	service.ReadUserData(w, r, &request)

	responseSum, err := db.GetSumUserSubscription(&request)
	if err != nil {
		log.Println("Failed to calculate subscriptions sum from database")
		http.Error(w, "Couldn't link user and database", http.StatusInternalServerError)
		return
	}

	log.Println("Creating response structure for sum calculation")
	var response = types.UserSubscriptionSumResponse{
		UserId:     request.UserId,
		CurrentSum: responseSum}

	log.Println("Marshaling sum response to JSON")
	bodyMarshal, err := json.Marshal(response)
	if err != nil {
		log.Println("Failed to marshal sum response")
		http.Error(w, "Failed to marshal user's subscription data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bodyMarshal)
	log.Println("Subscriptions sum sent successfully")
}
