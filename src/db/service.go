package db

import (
	"crudl_service/src/types"
	"database/sql"
	"log"
)

func CreateUserSubscription(data *types.UserSubscription) error {
	log.Println("Creating new subscription for user:", data.UserId, "service:", data.ServiceName)
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) 
			  VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(query, data.ServiceName, data.Price, data.UserId, data.StartDate, data.EndDate)
	if err != nil {
		log.Println("Failed to create subscription")
		return err
	}
	log.Println("Subscription created successfully for user:", data.UserId)
	return nil
}

func GetUserSubscription(data *types.UserSubscriptionData) (*types.UserSubscription, error) {
	log.Println("Getting subscription for user:", data.UserId, "service:", data.ServiceName)
	query := `SELECT service_name, price, user_id, start_date, end_date 
			  FROM subscriptions 
			  WHERE user_id = $1 AND service_name = $2`

	row := db.QueryRow(query, data.UserId, data.ServiceName)

	subscription := &types.UserSubscription{}
	err := row.Scan(&subscription.ServiceName, &subscription.Price, &subscription.UserId,
		&subscription.StartDate, &subscription.EndDate)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Subscription not found", "user_id", data.UserId, "service_name", data.ServiceName)
			return nil, err
		}
		log.Println("Failed to get subscription")
		return nil, err
	}

	log.Println("Subscription retrieved successfully", "user_id", data.UserId, "service_name", data.ServiceName)
	return subscription, nil
}

func UpdateUserSubscription(data *types.UserSubscription) error {
	log.Println("Updating subscription", "user_id", data.UserId, "service_name", data.ServiceName)
	query := `UPDATE subscriptions 
			  SET price = $1, start_date = $2, end_date = $3
			  WHERE user_id = $4 AND service_name = $5`

	result, err := db.Exec(query, data.Price, data.StartDate, data.EndDate, data.UserId, data.ServiceName)
	if err != nil {
		log.Println("Failed to update subscription")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Failed to get rows affected count")
		return err
	}

	if rowsAffected == 0 {
		log.Println("No subscription found for update", "user_id", data.UserId, "service_name", data.ServiceName)
	} else {
		log.Println("Subscription updated successfully", "user_id", data.UserId, "service_name", data.ServiceName)
	}

	return nil
}

func DeleteUserSubscription(data *types.UserSubscriptionData) error {
	log.Println("Deleting subscription", "user_id", data.UserId, "service_name", data.ServiceName)
	query := `DELETE FROM subscriptions WHERE user_id = $1 AND service_name = $2`

	result, err := db.Exec(query, data.UserId, data.ServiceName)
	if err != nil {
		log.Println("Failed to delete subscription")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Failed to get deleted rows count")
		return err
	}

	if rowsAffected == 0 {
		log.Println("No subscription found for deletion", "user_id", data.UserId, "service_name", data.ServiceName)
	} else {
		log.Println("Subscription deleted successfully", "user_id", data.UserId, "service_name", data.ServiceName)
	}

	return nil
}

func ListUserSubscriptions(data *types.UserRequest) ([]types.UserSubscription, error) {
	log.Println("Listing user subscriptions", "user_id", data.UserId)
	query := `SELECT service_name, price, user_id, start_date, end_date 
			  FROM subscriptions 
			  WHERE user_id = $1`

	rows, err := db.Query(query, data.UserId)
	if err != nil {
		log.Println("Failed to get subscriptions list")
		return nil, err
	}
	defer rows.Close()

	var subscriptions []types.UserSubscription

	for rows.Next() {
		var subscription types.UserSubscription
		err := rows.Scan(&subscription.ServiceName, &subscription.Price, &subscription.UserId,
			&subscription.StartDate, &subscription.EndDate)
		if err != nil {
			log.Println("Failed to scan row")
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	if err = rows.Err(); err != nil {
		log.Println("Failed to iterate rows")
		return nil, err
	}

	log.Println("User subscriptions listed successfully", "user_id", data.UserId, "count", len(subscriptions))
	return subscriptions, nil
}

func GetSumUserSubscription(data *types.UserSumSubscriptionRequest) (int64, error) {
	log.Println("Calculating sum of user subscriptions", "user_id", data.UserId, "start_date", data.StartDate, "end_date", data.EndDate)
	query := `SELECT COALESCE(SUM(price), 0) 
			  FROM subscriptions 
			  WHERE user_id = $1 
			  AND start_date <= $2 
			  AND (end_date IS NULL OR end_date >= $3)`

	var totalSum int64
	err := db.QueryRow(query, data.UserId, data.EndDate, data.StartDate).Scan(&totalSum)
	if err != nil {
		log.Println("Failed to calculate sum of user subscriptions")
		return 0, err
	}

	log.Println("Sum of user subscriptions calculated successfully", "user_id", data.UserId, "total_sum", totalSum)
	return totalSum, nil
}
