package users

import (
	userModel "date-app/models/users"
	users "date-app/repository/users/mysql"
	"date-app/utils/jwt"
	"date-app/utils/password"
	"errors"
	"time"
)

type Users interface {
	Create(data *userModel.CreateUserRequest) error
	Login(data *userModel.LoginRequest) (*userModel.LoginResponse, error)
}

type usersService struct {
	usersRepo users.Users
}

func NewUsersService(usersRepo users.Users) Users {
	return &usersService{
		usersRepo: usersRepo,
	}
}

func (u *usersService) Create(data *userModel.CreateUserRequest) error {
	var err error
	data.Password, err = password.HashPassword(data.Password)
	if err != nil {
		return errors.New("error hash password: " + err.Error())
	}
	err = u.usersRepo.Create(data)
	if err != nil {
		return errors.New("error create user: " + err.Error())
	}
	return nil
}

func (u *usersService) Login(data *userModel.LoginRequest) (*userModel.LoginResponse, error) {
	userData, err := u.usersRepo.GetByUsername(data.Username)
	if err != nil {
		return nil, errors.New("error get user: " + err.Error())
	}
	err = password.VerifyPassword(data.Password, userData.Password)
	if err != nil {
		return nil, errors.New("password is not matched")
	}
	expirationTime := time.Now().Add(2 * time.Hour)
	token, err := jwt.CreateToken(userData.ID, 0, expirationTime)
	if err != nil {
		return nil, errors.New("create token failed: " + err.Error())
	}
	tokenData := userModel.LoginResponse{
		Token:     token,
		ExpiredAt: expirationTime.Format("2006-01-02 15:04:05"),
	}
	return &tokenData, nil
}
