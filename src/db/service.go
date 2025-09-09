package db

import (
	"crudl_service/src/types"
)

func CreateUserSubscription(data *types.UserSubscription) error {
	return nil
}

func GetUserSubscription(data *types.UserSubscriptionData) (*types.UserSubscription, error) {
	return &types.UserSubscription{}, nil
}

func UpdateUserSubscription(data *types.UserSubscription) error {
	return nil
}

func DeleteUserSubscription(data *types.UserSubscriptionData) error {
	return nil
}

func ListUserSubscriptions(data *types.ListUserRequest) ([]types.UserSubscription, error) {
	return []types.UserSubscription{}, nil
}
