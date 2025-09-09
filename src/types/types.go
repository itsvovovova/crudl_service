package types

import "time"

type UserSubscription struct {
	ServiceName string     `json:"service_name"`
	Price       int64      `json:"price"`
	UserId      string     `json:"user_id"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

type UserSubscriptionData struct {
	UserId      string `json:"user_id"`
	ServiceName string `json:"service_name"`
}

type ListUserRequest struct {
	UserId string `json:"user_id"`
}
