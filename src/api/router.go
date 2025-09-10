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

// CreateSubscription creates a new subscription
//
//	@Summary		Create subscription
//	@Description	Create a new subscription for a user
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body	types.UserSubscription	true	"Subscription data"
//	@Success		201	"Subscription created successfully"
//	@Failure		400	{object}	string	"Bad request"
//	@Failure		500	{object}	string	"Internal server error"
//	@Router			/subscription [post]
func CreateSubscription(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating new subscription request received")
	var request types.UserSubscription
	service.ReadUserData(w, r, &request)

	log.Println("Validating start date format")
	if request.StartDate == nil {
		log.Println("Start date is required")
		http.Error(w, "Start date is required", http.StatusBadRequest)
		return
	}
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

// ReadSubscription gets a subscription
//
//	@Summary		Get subscription
//	@Description	Get subscription by user ID and service name
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body		types.UserSubscriptionData	true	"User and service data"
//	@Success		200				{object}	types.UserSubscription		"Subscription data"
//	@Failure		400				{object}	string						"Bad request"
//	@Failure		500				{object}	string						"Internal server error"
//	@Router			/subscription [get]
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

// UpdateSubscription updates an existing subscription
//
//	@Summary		Update subscription
//	@Description	Update subscription details
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body	types.UserSubscription	true	"Updated subscription data"
//	@Success		200	"Subscription updated successfully"
//	@Failure		400	{object}	string	"Bad request"
//	@Failure		500	{object}	string	"Internal server error"
//	@Router			/subscription [put]
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

// DeleteSubscription deletes a subscription
//
//	@Summary		Delete subscription
//	@Description	Delete subscription by user ID and service name
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body	types.UserSubscriptionData	true	"User and service data"
//	@Success		200	"Subscription deleted successfully"
//	@Failure		400	{object}	string	"Bad request"
//	@Failure		500	{object}	string	"Internal server error"
//	@Router			/subscription [delete]
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

// ListSubscription lists all user subscriptions
//
//	@Summary		List subscriptions
//	@Description	Get all subscriptions for a user
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			user	body		types.UserRequest			true	"User data"
//	@Success		200		{array}		types.UserSubscription		"List of subscriptions"
//	@Failure		400		{object}	string						"Bad request"
//	@Failure		500		{object}	string						"Internal server error"
//	@Router			/subscriptionList [get]
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

// SumUserSubscriptions calculates total subscription cost
//
//	@Summary		Calculate subscription sum
//	@Description	Calculate total cost of user subscriptions in date range
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		types.UserSumSubscriptionRequest	true	"Date range and user data"
//	@Success		200		{object}	types.UserSubscriptionSumResponse	"Total sum"
//	@Failure		400		{object}	string								"Bad request"
//	@Failure		500		{object}	string								"Internal server error"
//	@Router			/sum_subscriptions [get]
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
