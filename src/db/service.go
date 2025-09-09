package db

import (
	"crudl_service/src/types"
	"database/sql"
	"log"
)

func CreateUserSubscription(data *types.UserSubscription) error {
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) 
			  VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(query, data.ServiceName, data.Price, data.UserId, data.StartDate, data.EndDate)
	if err != nil {
		log.Printf("Failed to create subscription")
		return err
	}
	return nil
}

func GetUserSubscription(data *types.UserSubscriptionData) (*types.UserSubscription, error) {
	query := `SELECT service_name, price, user_id, start_date, end_date 
			  FROM subscriptions 
			  WHERE user_id = $1 AND service_name = $2`

	row := db.QueryRow(query, data.UserId, data.ServiceName)

	subscription := &types.UserSubscription{}
	err := row.Scan(&subscription.ServiceName, &subscription.Price, &subscription.UserId,
		&subscription.StartDate, &subscription.EndDate)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Subscription not found for user %s and service %s", data.UserId, data.ServiceName)
			return nil, err
		}
		log.Printf("Failed to get subscription")
		return nil, err
	}

	return subscription, nil
}

func UpdateUserSubscription(data *types.UserSubscription) error {
	query := `UPDATE subscriptions 
			  SET price = $1, start_date = $2, end_date = $3
			  WHERE user_id = $4 AND service_name = $5`

	result, err := db.Exec(query, data.Price, data.StartDate, data.EndDate, data.UserId, data.ServiceName)
	if err != nil {
		log.Printf("Failed to update subscription")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected count")
		return err
	}

	if rowsAffected == 0 {
		log.Printf("No subscription found for update")
	}

	return nil
}

func DeleteUserSubscription(data *types.UserSubscriptionData) error {
	query := `DELETE FROM subscriptions WHERE user_id = $1 AND service_name = $2`

	result, err := db.Exec(query, data.UserId, data.ServiceName)
	if err != nil {
		log.Printf("Failed to delete subscription")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get deleted rows count")
		return err
	}

	if rowsAffected == 0 {
		log.Printf("No subscription found for deletion")
	}

	return nil
}

func ListUserSubscriptions(data *types.UserRequest) ([]types.UserSubscription, error) {
	query := `SELECT service_name, price, user_id, start_date, end_date 
			  FROM subscriptions 
			  WHERE user_id = $1`

	rows, err := db.Query(query, data.UserId)
	if err != nil {
		log.Printf("Failed to get subscriptions list")
		return nil, err
	}
	defer rows.Close()

	var subscriptions []types.UserSubscription

	for rows.Next() {
		var subscription types.UserSubscription
		err := rows.Scan(&subscription.ServiceName, &subscription.Price, &subscription.UserId,
			&subscription.StartDate, &subscription.EndDate)
		if err != nil {
			log.Printf("Failed to scan row")
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Failed to iterate rows")
		return nil, err
	}

	return subscriptions, nil
}

func GetSumUserSubscription(data *types.UserSumSubscriptionRequest) (int64, error) {
	query := `SELECT COALESCE(SUM(price), 0) 
			  FROM subscriptions 
			  WHERE user_id = $1 
			  AND start_date <= $2 
			  AND (end_date IS NULL OR end_date >= $3)`

	var totalSum int64
	err := db.QueryRow(query, data.UserId, data.EndDate, data.StartDate).Scan(&totalSum)
	if err != nil {
		log.Printf("Failed to calculate sum of user subscriptions")
		return 0, err
	}

	return totalSum, nil
}
