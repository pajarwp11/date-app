package users

import (
	"database/sql"
	"date-app/models/users"
	"fmt"
	"strings"
)

type Users interface {
	Create(data *users.CreateUserRequest) error
	GetByUsername(username string) (*users.UserData, error)
	GetRandomUser(userId int, excludedId []string) (*users.UserResponse, error)
}
type usersRepository struct {
	DB *sql.DB
}

func NewUsersRepository(db *sql.DB) Users {
	return &usersRepository{
		DB: db,
	}
}

func (u *usersRepository) Create(data *users.CreateUserRequest) error {
	query := "INSERT INTO users (username,password,fullname,gender,location,education,occupation,bio) VALUES (?,?,?,?,?,?,?,?)"
	_, err := u.DB.Exec(query, data.Username, data.Password, data.Fullname, data.Gender, data.Location, data.Education, data.Occupation, data.Bio)
	return err
}

func (u *usersRepository) GetByUsername(username string) (*users.UserData, error) {
	query := "SELECT id,password,is_premium FROM users WHERE username=?"
	row := u.DB.QueryRow(query, username)

	var user users.UserData
	err := row.Scan(&user.ID, &user.Password, &user.IsPremium)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with username %s", username)
		}
		return nil, err
	}

	return &user, nil
}

func (u *usersRepository) GetRandomUser(userId int, excludedId []string) (*users.UserResponse, error) {
	placeholders := make([]string, len(excludedId))
	for i := range excludedId {
		placeholders[i] = "?"
	}
	placeholdersStr := strings.Join(placeholders, ",")

	query := fmt.Sprintf(`
		SELECT id, username, fullname, gender, location, education, occupation, bio
		FROM users
		WHERE id NOT IN (%s)
		AND id NOT IN (SELECT liked_user_id FROM date_app.user_likes WHERE user_id = ?)
		ORDER BY RAND() LIMIT 1`, placeholdersStr)

	args := make([]interface{}, len(excludedId)+1)
	for i, id := range excludedId {
		args[i] = id
	}
	args[len(excludedId)] = userId

	row := u.DB.QueryRow(query, args...)
	var user users.UserResponse
	err := row.Scan(&user.ID, &user.Username, &user.Fullname, &user.Gender, &user.Location, &user.Education, &user.Occupation, &user.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unable to find user")
		}
		return nil, err
	}

	return &user, nil
}
