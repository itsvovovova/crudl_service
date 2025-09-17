package db

import "crudl_service/src/types"

type SubscriptionRepository interface {
	Create(data *types.UserSubscription) (int64, error)
	Get(id int64) (*types.UserSubscription, error)
	Update(data *types.UserSubscription) error
	Delete(id int64) error
	List(userID string, afterID *int64, limit int) ([]types.UserSubscription, error)
	Sum(data *types.UserSumSubscriptionRequest) (int64, error)
}

type postgresRepository struct{}

func NewPostgresRepository() SubscriptionRepository {
	return &postgresRepository{}
}

func (r *postgresRepository) Create(data *types.UserSubscription) (int64, error) {
	return CreateUserSubscription(data)
}

func (r *postgresRepository) Get(id int64) (*types.UserSubscription, error) {
	return GetUserSubscription(id)
}

func (r *postgresRepository) Update(data *types.UserSubscription) error {
	return UpdateUserSubscription(data)
}

func (r *postgresRepository) Delete(id int64) error {
	return DeleteUserSubscription(id)
}

func (r *postgresRepository) List(userID string, afterID *int64, limit int) ([]types.UserSubscription, error) {
	return ListUserSubscriptions(userID, afterID, limit)
}

func (r *postgresRepository) Sum(data *types.UserSumSubscriptionRequest) (int64, error) {
	return GetSumUserSubscription(data)
}
