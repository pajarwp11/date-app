package users

import (
	"database/sql"
	"date-app/models/users"
	"fmt"
)

type Users interface {
	Create(data *users.CreateUserRequest) error
	GetByUsername(username string) (*users.UserData, error)
	GetRandomUser(excludedId string) (*users.UserData, error)
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
	query := "SELECT id,password FROM users WHERE username=?"
	row := u.DB.QueryRow(query, username)

	var user users.UserData
	err := row.Scan(&user.ID, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with username %s", username)
		}
		return nil, err
	}

	return &user, nil
}

func (u *usersRepository) GetRandomUser(excludedId string) (*users.UserData, error) {
	query := "SELECT username,fullname,gender,location,education,occupation,bio FROM users WHERE id NOT IN (?) ORDER BY RAND() LIMIT 1"
	row := u.DB.QueryRow(query, excludedId)

	var user users.UserData
	err := row.Scan(&user.Password, &user.Fullname, &user.Gender, &user.Location, &user.Education, &user.Occupation, &user.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unable to find user")
		}
		return nil, err
	}

	return &user, nil
}