package types

import "time"

type UserSubscription struct {
	ServiceName string  `json:"service_name"`
	Price       int64   `json:"price"`
	UserId      string  `json:"user_id"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

type UserSubscriptionData struct {
	UserId      string `json:"user_id"`
	ServiceName string `json:"service_name"`
}

type UserRequest struct {
	UserId string `json:"user_id"`
}

type UserSumSubscriptionRequest struct {
	UserId    string    `json:"user_id"`
	StartDate time.Time `json:"start_time"`
	EndDate   time.Time `json:"end_time"`
}

type UserSubscriptionSumResponse struct {
	UserId     string `json:"user_id"`
	CurrentSum int64  `json:"current_sum"`
}
