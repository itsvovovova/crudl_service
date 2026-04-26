package db

import (
	"crudl_service/src/types"
	"database/sql"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

var ErrNotFound = errors.New("not found")

func (r *postgresRepository) Create(data *types.UserSubscription) (int64, error) {
	if err := r.checkDB(); err != nil {
		return 0, err
	}
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
			  VALUES ($1, $2, $3, to_date($4, 'MM-YYYY'), CASE WHEN $5 IS NULL THEN NULL ELSE to_date($5, 'MM-YYYY') END) RETURNING id`
	var id int64
	if err := r.db.QueryRow(query, data.ServiceName, data.Price, data.UserId, data.StartDate, data.EndDate).Scan(&id); err != nil {
		log.WithError(err).Error("Failed to create subscription")
		return 0, err
	}
	return id, nil
}

func (r *postgresRepository) Get(id int64) (*types.UserSubscription, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	query := `SELECT service_name, price, user_id, to_char(start_date, 'MM-YYYY'), to_char(end_date, 'MM-YYYY')
			  FROM subscriptions WHERE id = $1`
	sub := &types.UserSubscription{Id: id}
	err := r.db.QueryRow(query, id).Scan(&sub.ServiceName, &sub.Price, &sub.UserId, &sub.StartDate, &sub.EndDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		log.WithError(err).Error("Failed to get subscription")
		return nil, err
	}
	return sub, nil
}

func (r *postgresRepository) Update(data *types.UserSubscription) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	query := `UPDATE subscriptions
			  SET price = $1, start_date = to_date($2, 'MM-YYYY'), end_date = CASE WHEN $3 IS NULL THEN NULL ELSE to_date($3, 'MM-YYYY') END
			  WHERE user_id = $4 AND service_name = $5`
	result, err := r.db.Exec(query, data.Price, data.StartDate, data.EndDate, data.UserId, data.ServiceName)
	if err != nil {
		log.WithError(err).Error("Failed to update subscription")
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *postgresRepository) Delete(id int64) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	result, err := r.db.Exec(`DELETE FROM subscriptions WHERE id = $1`, id)
	if err != nil {
		log.WithError(err).Error("Failed to delete subscription")
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *postgresRepository) List(userID string, afterID *int64, limit int) ([]types.UserSubscription, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	baseQuery := `SELECT id, service_name, price, user_id, to_char(start_date, 'MM-YYYY'), to_char(end_date, 'MM-YYYY')
				  FROM subscriptions WHERE user_id = $1`
	args := []interface{}{userID}

	if afterID != nil {
		baseQuery += " AND id > $2"
		args = append(args, *afterID)
	}
	baseQuery += " ORDER BY id ASC"
	if limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, limit)
	}

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		log.WithError(err).Error("Failed to list subscriptions")
		return nil, err
	}
	defer rows.Close()

	var subs []types.UserSubscription
	for rows.Next() {
		var s types.UserSubscription
		if err := rows.Scan(&s.Id, &s.ServiceName, &s.Price, &s.UserId, &s.StartDate, &s.EndDate); err != nil {
			log.WithError(err).Error("Failed to scan subscription row")
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, rows.Err()
}

func (r *postgresRepository) Sum(data *types.UserSumSubscriptionRequest) (int64, error) {
	if err := r.checkDB(); err != nil {
		return 0, err
	}
	if data == nil {
		return 0, fmt.Errorf("data cannot be nil")
	}
	query := `WITH params AS (
				  SELECT to_date($2, 'MM-YYYY') AS req_start, to_date($3, 'MM-YYYY') AS req_end
              ), selected AS (
				  SELECT s.price,
					     GREATEST(s.start_date, p.req_start) AS os,
					     LEAST(COALESCE(s.end_date, p.req_end), p.req_end) AS oe
				  FROM subscriptions s, params p
				  WHERE s.user_id = $1
					AND s.start_date <= p.req_end
					AND (s.end_date IS NULL OR s.end_date >= p.req_start)
              ), normalized AS (
				  SELECT price, os, oe FROM selected WHERE os <= oe
              )
              SELECT COALESCE(SUM(price * ((DATE_PART('year', age(oe, os))::int * 12) + DATE_PART('month', age(oe, os))::int + 1)), 0)
              FROM normalized`
	var total int64
	if err := r.db.QueryRow(query, data.UserId, data.StartDate, data.EndDate).Scan(&total); err != nil {
		log.WithError(err).Error("Failed to calculate subscription sum")
		return 0, err
	}
	return total, nil
}
