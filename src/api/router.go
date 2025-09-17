package api

import (
	"crudl_service/src/db"
	"crudl_service/src/service"
	"crudl_service/src/types"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

var CurrentRepository db.SubscriptionRepository

// InitAPI initializes the API layer with repository
func InitAPI(repository db.SubscriptionRepository) {
	CurrentRepository = repository
	log.Info("API layer initialized with repository")
}

// ShutdownAPI gracefully shuts down the API layer
func ShutdownAPI() {
	log.Info("Shutting down API layer")
	CurrentRepository = nil
	log.Info("API layer shutdown complete")
}

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
	log.Info("Creating new subscription request received")
	var request types.UserSubscription
	service.ReadUserData(w, r, &request)

	log.Info("Validating start date format")
	if request.StartDate == nil {
		log.Error("Start date is required")
		http.Error(w, "Start date is required", http.StatusBadRequest)
		return
	}
	if _, err := time.Parse("01-2006", *request.StartDate); err != nil {
		log.Error("Invalid start date format provided")
		http.Error(w, "Incorrect time format", http.StatusBadRequest)
		return
	}

	if request.EndDate != nil {
		log.Info("Validating end date format")
		_, err := time.Parse("01-2006", *request.EndDate)
		if err != nil {
			log.Error("Invalid end date format provided")
			http.Error(w, "Incorrect time format", http.StatusBadRequest)
			return
		}
	}

	id, err := CurrentRepository.Create(&request)
	if err != nil {
		log.Error("Failed to create subscription in database")
		http.Error(w, "Couldn't link subscription and database", http.StatusInternalServerError)
		return
	}

	log.Info("Subscription created successfully via API")
	var response = types.CreateSubscriptionResponse{
		Result:         "ok",
		SubscriptionId: id,
	}

	bodyMarshal, err := json.Marshal(response)
	if err != nil {
		log.Error("Failed to marshal subscription data")
		http.Error(w, "Failed to marshal user's subscription data", http.StatusInternalServerError)
		return
	}
	w.Write(bodyMarshal)
}

// ReadSubscription gets a subscription
//
//	@Summary		Get subscription
//	@Description	Get subscription by ID
//	@Tags			subscriptions
//	@Produce		json
//	@Param			id	path	int	true	"Subscription ID"
//	@Success		200	{object}	types.UserSubscription	"Subscription data"
//	@Failure		400	{object}	string	"Bad request"
//	@Failure		404	{object}	string	"Not found"
//	@Failure		500	{object}	string	"Internal server error"
//	@Router			/subscription/{id} [get]
func ReadSubscription(w http.ResponseWriter, r *http.Request) {
	log.Info("Read subscription request received")
	intID, err := service.GetIDRequest(r)
	if err != nil {
		log.Error("Failed to parse ID")
		http.Error(w, "Failed to parse ID", http.StatusBadRequest)
		return
	}
	responseBody, err := CurrentRepository.Get(intID)
	if err != nil {
		log.Error("Failed to retrieve subscription from database")
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	log.Info("Marshaling subscription data to JSON")
	bodyMarshal, err := json.Marshal(responseBody)
	if err != nil {
		log.Error("Failed to marshal subscription data")
		http.Error(w, "Failed to marshal user's subscription data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bodyMarshal)
	log.Info("Subscription data sent successfully")
}

// UpdateSubscription updates an existing subscription
//
//	@Summary		Update subscription
//	@Description	Update subscription details
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Subscription ID"
//	@Param			subscription	body	types.UserSubscription	true	"Updated subscription data"
//	@Success		200	"Subscription updated successfully"
//	@Failure		400	{object}	string	"Bad request"
//	@Failure		404	{object}	string	"Not found"
//	@Failure		500	{object}	string	"Internal server error"
//	@Router			/subscription/{id} [put]
func UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	log.Info("Update subscription request received")
	var request types.UserSubscription
	service.ReadUserData(w, r, &request)

	if err := CurrentRepository.Update(&request); err != nil {
		log.Error("Failed to update subscription in database")
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	log.Info("Subscription updated successfully via API")
	w.WriteHeader(http.StatusOK)
}

// DeleteSubscription deletes a subscription
//
//	@Summary		Delete subscription
//	@Description	Delete subscription by ID
//	@Tags			subscriptions
//	@Produce		json
//	@Param			id	path	int	true	"Subscription ID"
//	@Success		200	"Subscription deleted successfully"
//	@Failure		400	{object}	string	"Bad request"
//	@Failure		404	{object}	string	"Not found"
//	@Failure		500	{object}	string	"Internal server error"
//	@Router			/subscription/{id} [delete]
func DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	log.Info("Delete subscription request received")
	intID, err := service.GetIDRequest(r)
	if err != nil {
		log.Error("Failed to parse ID")
		http.Error(w, "Failed to parse ID", http.StatusBadRequest)
		return
	}
	if err := CurrentRepository.Delete(intID); err != nil {
		log.Error("Failed to delete subscription from database")
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	log.Info("Subscription deleted successfully via API")
	w.WriteHeader(http.StatusOK)
}

// ListSubscription lists all user subscriptions
//
//	@Summary		List subscriptions
//	@Description	Get all subscriptions for a user with pagination
//	@Tags			subscriptions
//	@Produce		json
//	@Param			user_id	query	string	true	"User ID"
//	@Param			after_id	query	int	false	"Return items after this ID"
//	@Param			limit	query	int	false	"Maximum items to return"
//	@Success		200	{object}	map[string]interface{}	"{ data: [...], next_after_id: number|null }"
//	@Failure		400	{object}	string	"Bad request"
//	@Failure		500	{object}	string	"Internal server error"
//	@Router			/subscriptionList [get]
func ListSubscription(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	var limitInt int
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limitInt = n
		}
	}
	var afterID *int64
	if v := r.URL.Query().Get("after_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			afterID = &id
		}
	}

	items, err := CurrentRepository.List(userID, afterID, limitInt)
	if err != nil {
		http.Error(w, "failed to retrieve subscriptions", http.StatusInternalServerError)
		return
	}

	var nextAfterID *int64
	if len(items) > 0 {
		id := items[len(items)-1].Id
		nextAfterID = &id
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(struct {
		Data        []types.UserSubscription `json:"data"`
		NextAfterID *int64                   `json:"next_after_id"`
	}{
		Data:        items,
		NextAfterID: nextAfterID,
	})
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
//	@Failure		400		{object}	string					"Bad request"
//	@Failure		500		{object}	string					"Internal server error"
//	@Router			/sum_subscriptions [get]
func SumUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	log.Info("Sum user subscriptions request received")
	var request types.UserSumSubscriptionRequest
	service.ReadUserData(w, r, &request)

	responseSum, err := CurrentRepository.Sum(&request)
	if err != nil {
		log.Error("Failed to calculate subscriptions sum from database")
		http.Error(w, "Couldn't link user and database", http.StatusInternalServerError)
		return
	}

	log.Info("Creating response structure for sum calculation")
	var response = types.UserSubscriptionSumResponse{
		UserId:     request.UserId,
		CurrentSum: responseSum}

	log.Info("Marshaling sum response to JSON")
	bodyMarshal, err := json.Marshal(response)
	if err != nil {
		log.Error("Failed to marshal sum response")
		http.Error(w, "Failed to marshal user's subscription data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bodyMarshal)
	log.Info("Subscriptions sum sent successfully")
}
