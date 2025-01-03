package users

import (
	"context"
	userModel "date-app/models/users"
	users "date-app/repository/users/mysql"
	usersRedis "date-app/repository/users/redis"
	"date-app/utils/jwt"
	"date-app/utils/password"
	"errors"
	"strconv"
	"time"
)

type Users interface {
	Create(data *userModel.CreateUserRequest) error
	Login(data *userModel.LoginRequest) (*userModel.LoginResponse, error)
	GetRandomUser(userId int) (*userModel.UserResponse, error)
	UpdateIsPremium(data *userModel.UpdateIsPremiumRequest) error
	UserLike(userId int, data *userModel.UserLikeRequest) (*userModel.UserLikeResponse, error)
}

type usersService struct {
	usersRepo      users.Users
	usersRedisRepo usersRedis.Users
}

func NewUsersService(usersRepo users.Users, usersRedisRepo usersRedis.Users) Users {
	return &usersService{
		usersRepo:      usersRepo,
		usersRedisRepo: usersRedisRepo,
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
	token, err := jwt.CreateToken(userData.ID, expirationTime)
	if err != nil {
		return nil, errors.New("create token failed: " + err.Error())
	}
	tokenData := userModel.LoginResponse{
		Token:     token,
		ExpiredAt: expirationTime.Format("2006-01-02 15:04:05"),
	}
	return &tokenData, nil
}

func (u *usersService) GetRandomUser(userId int) (*userModel.UserResponse, error) {
	ctx := context.Background()
	viewedUsers, err := u.usersRedisRepo.GetViewedUser(ctx, "user:view:"+strconv.Itoa(userId))
	if err != nil {
		return nil, errors.New("error get user: " + err.Error())
	}

	user, err := u.usersRepo.GetByID(userId)
	if err != nil {
		return nil, errors.New("error get user: " + err.Error())
	}
	if user.IsPremium == 0 && len(viewedUsers) >= 10 {
		return nil, errors.New("daily limit reached")
	}
	excludedUsers := viewedUsers
	excludedUsers = append(excludedUsers, strconv.Itoa(userId))

	randomUser, err := u.usersRepo.GetRandomUser(userId, excludedUsers)
	if err != nil {
		return nil, errors.New("error get user: " + err.Error())
	}
	viewedUsers = append(viewedUsers, strconv.Itoa(randomUser.ID))
	err = u.usersRedisRepo.SetViewedUser(ctx, "user:view:"+strconv.Itoa(userId), viewedUsers)
	if err != nil {
		return nil, errors.New("error get user: " + err.Error())
	}
	return randomUser, nil
}

func (u *usersService) UpdateIsPremium(data *userModel.UpdateIsPremiumRequest) error {
	err := u.usersRepo.UpdateIsPremium(data.UserID, *data.IsPremium)
	if err != nil {
		return errors.New("error update is premium: " + err.Error())
	}
	return nil
}

func (u *usersService) UserLike(userId int, data *userModel.UserLikeRequest) (*userModel.UserLikeResponse, error) {
	var err error
	var response userModel.UserLikeResponse

	tx, err := u.usersRepo.BeginTrx()
	if err != nil {
		return nil, errors.New("error starting transaction: " + err.Error())
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	userLikeData := userModel.UserLikes{
		UserID:      userId,
		LikedUserID: data.UserID,
	}
	err = u.usersRepo.InsertUserLikes(&userLikeData)
	if err != nil {
		return nil, errors.New("error like user: " + err.Error())
	}
	userLikeId, err := u.usersRepo.GetUserLike(data.UserID, userId)
	if err != nil {
		return nil, errors.New("error like user: " + err.Error())
	}
	if userLikeId != 0 {
		userMatchData := userModel.UserMatches{}
		if userId < data.UserID {
			userMatchData.UserID1 = userId
			userMatchData.UserID2 = data.UserID
		} else {
			userMatchData.UserID1 = data.UserID
			userMatchData.UserID2 = userId
		}

		err = u.usersRepo.InsertUserMatches(&userMatchData)
		if err != nil {
			return nil, errors.New("error match user: " + err.Error())
		}
		response.Message = "you are matched !!!"
		tx.Commit()
		return &response, nil
	}
	response.Message = "like success"
	tx.Commit()
	return &response, nil
}
