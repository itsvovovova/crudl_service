package types

type UserSubscription struct {
	Id          int64   `json:"id"`
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
	UserId    string `json:"user_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type UserSubscriptionSumResponse struct {
	UserId     string `json:"user_id"`
	CurrentSum int64  `json:"current_sum"`
}

type CreateSubscriptionResponse struct {
	Result         string `json:"result"`
	SubscriptionId int64  `json:"subscription_id"`
}

type UserRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}
