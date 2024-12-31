package users

import (
	"database/sql"
	"date-app/models/users"
	"fmt"
	"strings"
)

type Users interface {
	BeginTrx() (*sql.Tx, error)
	Create(data *users.CreateUserRequest) error
	GetByUsername(username string) (*users.UserData, error)
	GetRandomUser(userId int, excludedId []string) (*users.UserResponse, error)
	GetByID(id int) (*users.UserData, error)
	UpdateIsPremium(id int, status int) error
	InsertUserLikes(data *users.UserLikes) error
	InsertUserMatches(data *users.UserMatches) error
	GetUserLike(userId int, likedUserId int) (int, error)
}
type usersRepository struct {
	DB *sql.DB
}

func NewUsersRepository(db *sql.DB) Users {
	return &usersRepository{
		DB: db,
	}
}

func (u *usersRepository) BeginTrx() (*sql.Tx, error) {
	return u.DB.Begin()
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

func (u *usersRepository) GetByID(id int) (*users.UserData, error) {
	query := "SELECT is_premium FROM users WHERE id=?"
	row := u.DB.QueryRow(query, id)

	var user users.UserData
	err := row.Scan(&user.IsPremium)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with id %d", id)
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

func (u *usersRepository) UpdateIsPremium(id int, status int) error {
	query := "UPDATE users SET is_premium=? WHERE id=?"
	_, err := u.DB.Exec(query, status, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *usersRepository) InsertUserLikes(data *users.UserLikes) error {
	query := "INSERT INTO user_likes (user_id,liked_user_id) VALUES (?,?)"
	_, err := u.DB.Exec(query, data.UserID, data.LikedUserID)
	return err
}

func (u *usersRepository) InsertUserMatches(data *users.UserMatches) error {
	query := "INSERT INTO user_matches (user_id_1,user_id_2) VALUES (?,?)"
	_, err := u.DB.Exec(query, data.UserID1, data.UserID2)
	return err
}

func (u *usersRepository) GetUserLike(userId int, likedUserId int) (int, error) {
	query := "SELECT id FROM user_likes WHERE user_id=? AND liked_user_id=?"
	row := u.DB.QueryRow(query, userId, likedUserId)

	var id int
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return id, nil
		}
		return id, err
	}

	return id, nil
}
