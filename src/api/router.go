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

const dateFormat = "01-2006"

// App holds all application dependencies.
type App struct {
	repo      db.Repository
	jwtSecret string
}

func NewApp(repo db.Repository, jwtSecret string) *App {
	return &App{repo: repo, jwtSecret: jwtSecret}
}

// CreateSubscription creates a new subscription
//
//	@Summary		Create subscription
//	@Description	Create a new subscription for a user
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body		types.UserSubscription				true	"Subscription data"
//	@Success		201				{object}	types.CreateSubscriptionResponse	"Subscription created"
//	@Failure		400				{object}	string								"Bad request"
//	@Failure		500				{object}	string								"Internal server error"
//	@Router			/subscription [post]
func (a *App) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var request types.UserSubscription
	if !service.ReadUserData(w, r, &request) {
		return
	}
	request.UserId = r.Header.Get("User-ID")

	if request.StartDate == nil {
		http.Error(w, "Start date is required", http.StatusBadRequest)
		return
	}
	if _, err := time.Parse(dateFormat, *request.StartDate); err != nil {
		http.Error(w, "Incorrect time format, expected MM-YYYY", http.StatusBadRequest)
		return
	}
	if request.EndDate != nil {
		if _, err := time.Parse(dateFormat, *request.EndDate); err != nil {
			http.Error(w, "Incorrect time format, expected MM-YYYY", http.StatusBadRequest)
			return
		}
	}

	id, err := a.repo.Create(&request)
	if err != nil {
		log.WithError(err).Error("Failed to create subscription")
		http.Error(w, "Failed to create subscription", http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(types.CreateSubscriptionResponse{Result: "ok", SubscriptionId: id})
	if err != nil {
		log.WithError(err).Error("Failed to marshal response")
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

// ReadSubscription gets a subscription by ID
//
//	@Summary		Get subscription
//	@Description	Get subscription by ID
//	@Tags			subscriptions
//	@Produce		json
//	@Param			id	path		int						true	"Subscription ID"
//	@Success		200	{object}	types.UserSubscription	"Subscription data"
//	@Failure		400	{object}	string					"Bad request"
//	@Failure		403	{object}	string					"Forbidden"
//	@Failure		404	{object}	string					"Not found"
//	@Router			/subscription/{id} [get]
func (a *App) ReadSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := service.GetIDRequest(r)
	if err != nil {
		http.Error(w, "Failed to parse ID", http.StatusBadRequest)
		return
	}
	sub, err := a.repo.Get(id)
	if err != nil {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}
	if sub.UserId != r.Header.Get("User-ID") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	body, err := json.Marshal(sub)
	if err != nil {
		log.WithError(err).Error("Failed to marshal subscription")
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// UpdateSubscription updates an existing subscription
//
//	@Summary		Update subscription
//	@Description	Update subscription details
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			id				path	int						true	"Subscription ID"
//	@Param			subscription	body	types.UserSubscription	true	"Updated subscription data"
//	@Success		200				"Subscription updated"
//	@Failure		400				{object}	string	"Bad request"
//	@Failure		404				{object}	string	"Not found"
//	@Router			/subscription/{id} [put]
func (a *App) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	var request types.UserSubscription
	if !service.ReadUserData(w, r, &request) {
		return
	}
	request.UserId = r.Header.Get("User-ID")

	if err := a.repo.Update(&request); err != nil {
		log.WithError(err).Error("Failed to update subscription")
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteSubscription deletes a subscription
//
//	@Summary		Delete subscription
//	@Description	Delete subscription by ID
//	@Tags			subscriptions
//	@Param			id	path	int	true	"Subscription ID"
//	@Success		200	"Subscription deleted"
//	@Failure		400	{object}	string	"Bad request"
//	@Failure		403	{object}	string	"Forbidden"
//	@Failure		404	{object}	string	"Not found"
//	@Router			/subscription/{id} [delete]
func (a *App) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := service.GetIDRequest(r)
	if err != nil {
		http.Error(w, "Failed to parse ID", http.StatusBadRequest)
		return
	}
	existing, err := a.repo.Get(id)
	if err != nil {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}
	if existing.UserId != r.Header.Get("User-ID") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if err := a.repo.Delete(id); err != nil {
		log.WithError(err).Error("Failed to delete subscription")
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// ListSubscription lists subscriptions for the authenticated user
//
//	@Summary		List subscriptions
//	@Description	Paginated list of subscriptions for the current user
//	@Tags			subscriptions
//	@Produce		json
//	@Param			after_id	query		int		false	"Cursor: return items after this ID"
//	@Param			limit		query		int		false	"Max items to return"
//	@Success		200			{object}	map[string]interface{}	"{ data: [...], next_after_id: number|null }"
//	@Failure		500			{object}	string	"Internal server error"
//	@Router			/subscriptionList [get]
func (a *App) ListSubscription(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("User-ID")

	var limit int
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	var afterID *int64
	if v := r.URL.Query().Get("after_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			afterID = &id
		}
	}

	items, err := a.repo.List(userID, afterID, limit)
	if err != nil {
		log.WithError(err).Error("Failed to list subscriptions")
		http.Error(w, "Failed to retrieve subscriptions", http.StatusInternalServerError)
		return
	}

	var nextAfterID *int64
	if len(items) > 0 {
		id := items[len(items)-1].Id
		nextAfterID = &id
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(struct {
		Data        []types.UserSubscription `json:"data"`
		NextAfterID *int64                   `json:"next_after_id"`
	}{Data: items, NextAfterID: nextAfterID}); err != nil {
		log.WithError(err).Error("Failed to encode list response")
	}
}

// SumUserSubscriptions calculates total subscription cost
//
//	@Summary		Calculate subscription sum
//	@Description	Total cost of subscriptions in a date range
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		types.UserSumSubscriptionRequest	true	"Date range"
//	@Success		200		{object}	types.UserSubscriptionSumResponse	"Total sum"
//	@Failure		400		{object}	string								"Bad request"
//	@Failure		500		{object}	string								"Internal server error"
//	@Router			/sum_subscriptions [post]
func (a *App) SumUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	var request types.UserSumSubscriptionRequest
	if !service.ReadUserData(w, r, &request) {
		return
	}
	request.UserId = r.Header.Get("User-ID")

	total, err := a.repo.Sum(&request)
	if err != nil {
		log.WithError(err).Error("Failed to calculate subscription sum")
		http.Error(w, "Failed to calculate subscription sum", http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(types.UserSubscriptionSumResponse{UserId: request.UserId, CurrentSum: total})
	if err != nil {
		log.WithError(err).Error("Failed to marshal sum response")
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
