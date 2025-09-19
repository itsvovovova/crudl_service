package db

import "crudl_service/src/types"

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

func GetUserByUsername(username string) (*User, error) {
	query := `SELECT id, username, password FROM users WHERE username = $1`
	row := DB.QueryRow(query, username)

	user := &User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func CreateUser(username, hashedPassword string) (string, error) {
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
	var userID string
	err := DB.QueryRow(query, username, hashedPassword).Scan(&userID)
	return userID, err
}

func CheckTaskOwnership(userID, taskID string) (bool, error) {
	query := `SELECT COUNT(*) FROM tasks WHERE id = $1 AND user_id = $2`
	var count int
	err := DB.QueryRow(query, taskID, userID).Scan(&count)
	return count > 0, err
}
