package users

import (
	userModel "date-app/models/users"
	"date-app/repository/users"
	"date-app/utils/jwt"
	"date-app/utils/password"
	"encoding/json"
	"errors"
	"time"
)

type Users interface {
	Create(data *userModel.CreateUserRequest) error
	Login(data *userModel.LoginRequest) (*userModel.LoginResponse, error)
	Logout(token string)
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
	if data.AdditionalData != nil {
		addData, err := json.Marshal(data.AdditionalData)
		if err != nil {
			return errors.New("error marshal additional data: " + err.Error())
		}
		data.AdditionalData = string(addData)
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

func (u *usersService) Logout(token string) {
	jwt.TakeOutToken(token)
}