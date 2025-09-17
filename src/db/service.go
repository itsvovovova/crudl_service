package db

import (
	"crudl_service/src/types"
	"database/sql"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

var ErrNotFound = errors.New("not found")

func CreateUserSubscription(data *types.UserSubscription) (int64, error) {
	log.Info("Creating new subscription for user:", data.UserId, "service:", data.ServiceName)
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) 
			  VALUES ($1, $2, $3, to_date($4, 'MM-YYYY'), CASE WHEN $5 IS NULL THEN NULL ELSE to_date($5, 'MM-YYYY') END) RETURNING id`

	var id int64
	err := db.QueryRow(query, data.ServiceName, data.Price, data.UserId, data.StartDate, data.EndDate).Scan(&id)
	if err != nil {
		log.Error("Failed to create subscription")
		return 0, err
	}
	log.Info("Subscription created successfully for user:", data.UserId)
	return id, nil
}

func GetUserSubscription(idSubscription int64) (*types.UserSubscription, error) {
	log.Info("Getting subscription id:", idSubscription)
	query := `SELECT service_name, price, user_id, to_char(start_date, 'MM-YYYY') AS start_date, to_char(end_date, 'MM-YYYY') AS end_date 
			  FROM subscriptions 
			  WHERE id = $1`

	row := db.QueryRow(query, idSubscription)

	subscription := &types.UserSubscription{}
	err := row.Scan(&subscription.ServiceName, &subscription.Price, &subscription.UserId,
		&subscription.StartDate, &subscription.EndDate)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Info("Subscription not found", "id", idSubscription)
			return nil, err
		}
		log.Error("Failed to get subscription")
		return nil, err
	}

	log.Info("Subscription retrieved successfully", "id", idSubscription)
	return subscription, nil
}

func UpdateUserSubscription(data *types.UserSubscription) error {
	log.Info("Updating subscription", "user_id", data.UserId, "service_name", data.ServiceName)
	query := `UPDATE subscriptions 
			  SET price = $1, start_date = to_date($2, 'MM-YYYY'), end_date = CASE WHEN $3 IS NULL THEN NULL ELSE to_date($3, 'MM-YYYY') END
			  WHERE user_id = $4 AND service_name = $5`

	result, err := db.Exec(query, data.Price, data.StartDate, data.EndDate, data.UserId, data.ServiceName)
	if err != nil {
		log.Error("Failed to update subscription")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("Failed to get rows affected count")
		return err
	}

	if rowsAffected == 0 {
		log.Info("No subscription found for update", "user_id", data.UserId, "service_name", data.ServiceName)
		return ErrNotFound
	}

	log.Info("Subscription updated successfully", "user_id", data.UserId, "service_name", data.ServiceName)
	return nil
}

func DeleteUserSubscription(idSubscription int64) error {
	log.Info("Deleting subscription", "id", idSubscription)
	query := `DELETE FROM subscriptions WHERE id = $1`

	result, err := db.Exec(query, idSubscription)
	if err != nil {
		log.Error("Failed to delete subscription")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("Failed to get deleted rows count")
		return err
	}

	if rowsAffected == 0 {
		log.Info("No subscription found for deletion", "id", idSubscription)
		return ErrNotFound
	}

	log.Info("Subscription deleted successfully", "id", idSubscription)
	return nil
}

func ListUserSubscriptions(userID string, afterID *int64, limit int) ([]types.UserSubscription, error) {
	log.Info("Listing user subscriptions", "user_id", userID)

	baseQuery := `SELECT id, service_name, price, user_id, to_char(start_date, 'MM-YYYY') AS start_date, to_char(end_date, 'MM-YYYY') AS end_date 
                   FROM subscriptions 
                   WHERE user_id = $1`

	args := []interface{}{userID}

	if afterID != nil {
		baseQuery += " AND id > $2"
		args = append(args, *afterID)
	}

	baseQuery += " ORDER BY id ASC"

	if limit > 0 {
		if afterID != nil {
			baseQuery += " LIMIT $3"
		} else {
			baseQuery += " LIMIT $2"
		}
		args = append(args, limit)
	}

	rows, err := db.Query(baseQuery, args...)
	if err != nil {
		log.Error("Failed to get subscriptions list")
		return nil, err
	}
	defer rows.Close()

	var subscriptions []types.UserSubscription

	for rows.Next() {
		var subscription types.UserSubscription
		err := rows.Scan(&subscription.Id, &subscription.ServiceName, &subscription.Price, &subscription.UserId,
			&subscription.StartDate, &subscription.EndDate)
		if err != nil {
			log.Error("Failed to scan row")
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	if err = rows.Err(); err != nil {
		log.Error("Failed to iterate rows")
		return nil, err
	}

	log.Info("User subscriptions listed successfully", "user_id", userID, "count", len(subscriptions))
	return subscriptions, nil
}

func GetSumUserSubscription(data *types.UserSumSubscriptionRequest) (int64, error) {
	if data == nil {
		log.Error("Cannot calculate sum: nil data provided")
		return 0, fmt.Errorf("data cannot be nil")
	}
	log.Infof("Calculating sum of user subscriptions for user_id: %s, start_date: %s, end_date: %s", data.UserId, data.StartDate, data.EndDate)
	query := `WITH params AS (
                  SELECT to_date($2, 'MM-YYYY') AS req_start,
                         to_date($3, 'MM-YYYY') AS req_end
              ), selected AS (
                  SELECT s.price,
                         GREATEST(s.start_date, p.req_start) AS os,
                         LEAST(COALESCE(s.end_date, p.req_end), p.req_end) AS oe
                  FROM subscriptions s, params p
                  WHERE s.user_id = $1
                    AND s.start_date <= p.req_end
                    AND (s.end_date IS NULL OR s.end_date >= p.req_start)
              ), normalized AS (
                  SELECT price, os, oe
                  FROM selected
                  WHERE os <= oe
              )
              SELECT COALESCE(SUM(
                         price * (
                           (DATE_PART('year', age(oe, os))::int * 12)
                           + DATE_PART('month', age(oe, os))::int
                           + 1
                         )
                       ), 0)
              FROM normalized`

	var totalSum int64
	err := db.QueryRow(query, data.UserId, data.StartDate, data.EndDate).Scan(&totalSum)
	if err != nil {
		log.Error("Failed to calculate sum of user subscriptions")
		return 0, err
	}

	log.Info("Sum of user subscriptions calculated successfully", "user_id", data.UserId, "total_sum", totalSum)
	return totalSum, nil
}
