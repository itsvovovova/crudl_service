package db

import (
	"crudl_service/src/types"
	"database/sql"
	"fmt"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SubscriptionRepository interface {
	Create(data *types.UserSubscription) (int64, error)
	Get(id int64) (*types.UserSubscription, error)
	Update(data *types.UserSubscription) error
	Delete(id int64) error
	List(userID string, afterID *int64, limit int) ([]types.UserSubscription, error)
	Sum(data *types.UserSumSubscriptionRequest) (int64, error)
}

type UserRepository interface {
	GetUserByUsername(username string) (*User, error)
	CreateUser(username, hashedPassword string) (string, error)
}

// Repository combines subscription and user operations.
type Repository interface {
	SubscriptionRepository
	UserRepository
}

type postgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository returns a Repository backed by PostgreSQL.
func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) checkDB() error {
	if r.db == nil {
		return fmt.Errorf("database connection not initialized")
	}
	return nil
}

func (r *postgresRepository) GetUserByUsername(username string) (*User, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	user := &User{}
	if err := r.db.QueryRow(
		`SELECT id, username, password FROM users WHERE username = $1`, username,
	).Scan(&user.ID, &user.Username, &user.Password); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *postgresRepository) CreateUser(username, hashedPassword string) (string, error) {
	if err := r.checkDB(); err != nil {
		return "", err
	}
	var userID string
	err := r.db.QueryRow(
		`INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`,
		username, hashedPassword,
	).Scan(&userID)
	return userID, err
}
